package command

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/store"
	"github.com/cs3org/reva/v2/pkg/token/manager/jwt"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/justinas/alice"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	pkgmiddleware "github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	"github.com/owncloud/ocis/v2/ocis-pkg/oidc"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	storesvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/store/v0"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/autoprovision"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/logging"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/middleware"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/proxy"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/router"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/server/debug"
	proxyHTTP "github.com/owncloud/ocis/v2/services/proxy/pkg/server/http"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/tracing"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/user/backend"
	"github.com/owncloud/ocis/v2/services/proxy/pkg/userroles"
	"github.com/urfave/cli/v2"
	microstore "go-micro.dev/v4/store"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			userInfoCache := store.Create(
				store.Store(cfg.OIDC.UserinfoCache.Store),
				store.TTL(cfg.OIDC.UserinfoCache.TTL),
				store.Size(cfg.OIDC.UserinfoCache.Size),
				microstore.Nodes(cfg.OIDC.UserinfoCache.Nodes...),
				microstore.Database(cfg.OIDC.UserinfoCache.Database),
				microstore.Table(cfg.OIDC.UserinfoCache.Table),
			)

			logger := logging.Configure(cfg.Service.Name, cfg.Log)
			err := tracing.Configure(cfg)
			if err != nil {
				return err
			}
			err = grpc.Configure(grpc.GetClientOptions(cfg.GRPCClientTLS)...)
			if err != nil {
				return err
			}

			var oidcHTTPClient = &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						MinVersion:         tls.VersionTLS12,
						InsecureSkipVerify: cfg.OIDC.Insecure, //nolint:gosec
					},
					DisableKeepAlives: true,
				},
				Timeout: time.Second * 10,
			}

			oidcClient := oidc.NewOIDCClient(
				oidc.WithAccessTokenVerifyMethod(cfg.OIDC.AccessTokenVerifyMethod),
				oidc.WithLogger(logger),
				oidc.WithHTTPClient(oidcHTTPClient),
				oidc.WithOidcIssuer(cfg.OIDC.Issuer),
				oidc.WithJWKSOptions(cfg.OIDC.JWKS),
			)

			var (
				m = metrics.New()
			)

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			defer cancel()

			m.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			rp, err := proxy.NewMultiHostReverseProxy(
				proxy.Logger(logger),
				proxy.Config(cfg),
			)

			lh := StaticRouteHandler{
				prefix:        cfg.HTTP.Root,
				userInfoCache: userInfoCache,
				logger:        logger,
				config:        *cfg,
				oidcClient:    oidcClient,
				proxy:         rp,
			}
			if err != nil {
				return fmt.Errorf("failed to initialize reverse proxy: %w", err)
			}

			{
				middlewares := loadMiddlewares(ctx, logger, cfg, userInfoCache)
				server, err := proxyHTTP.Server(
					proxyHTTP.Handler(lh.handler()),
					proxyHTTP.Logger(logger),
					proxyHTTP.Context(ctx),
					proxyHTTP.Config(cfg),
					proxyHTTP.Metrics(metrics.New()),
					proxyHTTP.Middlewares(middlewares),
				)

				if err != nil {
					logger.Error().
						Err(err).
						Str("server", "http").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(err error) {
					logger.Error().
						Err(err).
						Str("server", "http").
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Error().Err(err).Str("server", "debug").Msg("Failed to initialize server")
					return err
				}

				gr.Add(server.ListenAndServe, func(_ error) {
					_ = server.Shutdown(ctx)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}

// StaticRouteHandler defines a Route Handler for static routes
type StaticRouteHandler struct {
	prefix        string
	proxy         http.Handler
	userInfoCache microstore.Store
	logger        log.Logger
	config        config.Config
	oidcClient    oidc.OIDCClient
}

func (h *StaticRouteHandler) handler() http.Handler {
	m := chi.NewMux()
	m.Route(h.prefix, func(r chi.Router) {
		// Wrapper for backchannel logout
		r.Post("/backchannel_logout", h.backchannelLogout)

		// TODO: migrate oidc well knowns here in a second wrapper
		r.HandleFunc("/*", h.proxy.ServeHTTP)

	})
	// This is commented out due to a race issue in chi
	//var methods = []string{"PROPFIND", "DELETE", "PROPPATCH", "MKCOL", "COPY", "MOVE", "LOCK", "UNLOCK", "REPORT"}
	//for _, k := range methods {
	//	chi.RegisterMethod(k)
	//}

	// To avoid using the chi.RegisterMethod() this is basically a catchAll for all HTTP Methods that are not
	// covered in chi by default
	m.MethodNotAllowed(h.proxy.ServeHTTP)

	return m
}

type jse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// handle backchannel logout requests as per https://openid.net/specs/openid-connect-backchannel-1_0.html#BCRequest
func (h *StaticRouteHandler) backchannelLogout(w http.ResponseWriter, r *http.Request) {
	// parse the application/x-www-form-urlencoded POST request
	logger := h.logger.SubloggerWithRequestID(r.Context())
	if err := r.ParseForm(); err != nil {
		logger.Warn().Err(err).Msg("ParseForm failed")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	logoutToken, err := h.oidcClient.VerifyLogoutToken(r.Context(), r.PostFormValue("logout_token"))
	if err != nil {
		logger.Warn().Err(err).Msg("VerifyLogoutToken failed")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	records, err := h.userInfoCache.Read(logoutToken.SessionId)
	if errors.Is(err, microstore.ErrNotFound) || len(records) == 0 {
		render.Status(r, http.StatusOK)
		render.JSON(w, r, nil)
		return
	}

	if err != nil {
		logger.Error().Err(err).Msg("Error reading userinfo cache")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
		return
	}

	for _, record := range records {
		err = h.userInfoCache.Delete(string(record.Value))
		if err != nil && !errors.Is(err, microstore.ErrNotFound) {
			// Spec requires us to return a 400 BadRequest when the session could not be destroyed
			logger.Err(err).Msg("could not delete user info from cache")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, jse{Error: "invalid_request", ErrorDescription: err.Error()})
			return
		}
		logger.Debug().Msg("Deleted userinfo from cache")
	}

	// we can ignore errors when cleaning up the lookup table
	err = h.userInfoCache.Delete(logoutToken.SessionId)
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to cleanup sessionid lookup entry")
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, nil)
}

func loadMiddlewares(ctx context.Context, logger log.Logger, cfg *config.Config, userInfoCache microstore.Store) alice.Chain {
	grpcClient, err := grpc.NewClient(append(grpc.GetClientOptions(cfg.GRPCClientTLS), grpc.WithTraceProvider(tracing.TraceProvider))...)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get gateway client")
	}
	rolesClient := settingssvc.NewRoleService("com.owncloud.api.settings", grpcClient)
	revaClient, err := pool.GetGatewayServiceClient(cfg.Reva.Address, cfg.Reva.GetRevaOptions()...)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get gateway client")
	}
	tokenManager, err := jwt.New(map[string]interface{}{
		"secret": cfg.TokenManager.JWTSecret,
	})
	if err != nil {
		logger.Fatal().Err(err).
			Msg("Failed to create token manager")
	}
	autoProvsionCreator := autoprovision.NewCreator(autoprovision.WithTokenManager(tokenManager))
	var userProvider backend.UserBackend
	switch cfg.AccountBackend {
	case "cs3":

		userProvider = backend.NewCS3UserBackend(
			backend.WithLogger(logger),
			backend.WithRevaAuthenticator(revaClient),
			backend.WithMachineAuthAPIKey(cfg.MachineAuthAPIKey),
			backend.WithOIDCissuer(cfg.OIDC.Issuer),
			backend.WithAutoProvisonCreator(autoProvsionCreator),
		)
	default:
		logger.Fatal().Msgf("Invalid accounts backend type '%s'", cfg.AccountBackend)
	}

	var roleAssigner userroles.UserRoleAssigner
	switch cfg.RoleAssignment.Driver {
	case "default":
		roleAssigner = userroles.NewDefaultRoleAssigner(
			userroles.WithRoleService(rolesClient),
			userroles.WithLogger(logger),
		)
	case "oidc":
		roleAssigner = userroles.NewOIDCRoleAssigner(
			userroles.WithRoleService(rolesClient),
			userroles.WithLogger(logger),
			userroles.WithRolesClaim(cfg.RoleAssignment.OIDCRoleMapper.RoleClaim),
			userroles.WithRoleMapping(cfg.RoleAssignment.OIDCRoleMapper.RolesMap),
			userroles.WithAutoProvisonCreator(autoProvsionCreator),
		)
	default:
		logger.Fatal().Msgf("Invalid role assignment driver '%s'", cfg.RoleAssignment.Driver)
	}

	storeClient := storesvc.NewStoreService("com.owncloud.api.store", grpcClient)
	if err != nil {
		logger.Error().Err(err).
			Str("gateway", cfg.Reva.Address).
			Msg("Failed to create reva gateway service client")
	}

	var oidcHTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS12,
				InsecureSkipVerify: cfg.OIDC.Insecure, //nolint:gosec
			},
			DisableKeepAlives: true,
		},
		Timeout: time.Second * 10,
	}

	var authenticators []middleware.Authenticator
	if cfg.EnableBasicAuth {
		logger.Warn().Msg("basic auth enabled, use only for testing or development")
		authenticators = append(authenticators, middleware.BasicAuthenticator{
			Logger:        logger,
			UserProvider:  userProvider,
			UserCS3Claim:  cfg.UserCS3Claim,
			UserOIDCClaim: cfg.UserOIDCClaim,
		})
	}

	authenticators = append(authenticators, middleware.NewOIDCAuthenticator(
		middleware.Logger(logger),
		middleware.UserInfoCache(userInfoCache),
		middleware.DefaultAccessTokenTTL(cfg.OIDC.UserinfoCache.TTL),
		middleware.HTTPClient(oidcHTTPClient),
		middleware.OIDCIss(cfg.OIDC.Issuer),
		middleware.OIDCClient(oidc.NewOIDCClient(
			oidc.WithAccessTokenVerifyMethod(cfg.OIDC.AccessTokenVerifyMethod),
			oidc.WithLogger(logger),
			oidc.WithHTTPClient(oidcHTTPClient),
			oidc.WithOidcIssuer(cfg.OIDC.Issuer),
			oidc.WithJWKSOptions(cfg.OIDC.JWKS),
		)),
	))
	authenticators = append(authenticators, middleware.PublicShareAuthenticator{
		Logger:            logger,
		RevaGatewayClient: revaClient,
	})
	authenticators = append(authenticators, middleware.SignedURLAuthenticator{
		Logger:             logger,
		PreSignedURLConfig: cfg.PreSignedURL,
		UserProvider:       userProvider,
		UserRoleAssigner:   roleAssigner,
		Store:              storeClient,
	})

	return alice.New(
		// first make sure we log all requests and redirect to https if necessary
		middleware.Tracer(),
		pkgmiddleware.TraceContext,
		chimiddleware.RealIP,
		chimiddleware.RequestID,
		middleware.AccessLog(logger),
		middleware.HTTPSRedirect,
		middleware.OIDCWellKnownRewrite(
			logger, cfg.OIDC.Issuer,
			cfg.OIDC.RewriteWellKnown,
			oidcHTTPClient,
		),
		router.Middleware(cfg.PolicySelector, cfg.Policies, logger),
		middleware.Authentication(
			authenticators,
			middleware.CredentialsByUserAgent(cfg.AuthMiddleware.CredentialsByUserAgent),
			middleware.Logger(logger),
			middleware.OIDCIss(cfg.OIDC.Issuer),
			middleware.EnableBasicAuth(cfg.EnableBasicAuth),
		),
		middleware.AccountResolver(
			middleware.Logger(logger),
			middleware.UserProvider(userProvider),
			middleware.UserRoleAssigner(roleAssigner),
			middleware.UserOIDCClaim(cfg.UserOIDCClaim),
			middleware.UserCS3Claim(cfg.UserCS3Claim),
			middleware.AutoprovisionAccounts(cfg.AutoprovisionAccounts),
		),
		middleware.SelectorCookie(
			middleware.Logger(logger),
			middleware.PolicySelectorConfig(*cfg.PolicySelector),
		),
		middleware.Policies(logger, cfg.PoliciesMiddleware.Query),
		// finally, trigger home creation when a user logs in
		middleware.CreateHome(
			middleware.Logger(logger),
			middleware.RevaGatewayClient(revaClient),
			middleware.RoleQuotas(cfg.RoleQuotas),
		),
	)
}
