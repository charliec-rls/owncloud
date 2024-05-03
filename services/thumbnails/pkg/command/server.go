package command

import (
	"context"
	"fmt"
	"os"

	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/logging"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/server/debug"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/server/grpc"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/server/http"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", cfg.Service.Name),
		Category: "server",
		Before: func(_ *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(_ *cli.Context) error {
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			cfg.GrpcClient, err = ogrpc.NewClient(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			if err != nil {
				return err
			}

			var (
				gr          = run.Group{}
				ctx, cancel = func() (context.Context, context.CancelFunc) {
					if cfg.Context == nil {
						return context.WithCancel(context.Background())
					}
					return context.WithCancel(cfg.Context)
				}()
				m = metrics.New()
			)

			defer cancel()

			m.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			service := grpc.NewService(
				grpc.Logger(logger),
				grpc.Context(ctx),
				grpc.Config(cfg),
				grpc.Name(cfg.Service.Name),
				grpc.Namespace(cfg.GRPC.Namespace),
				grpc.Address(cfg.GRPC.Addr),
				grpc.Metrics(m),
				grpc.TraceProvider(traceProvider),
			)

			gr.Add(service.Run, func(_ error) {
				logger.Error().
					Err(err).
					Str("server", "grpc").
					Msg("Shutting down server")

				cancel()
				os.Exit(1)
			})

			server, err := debug.Server(
				debug.Logger(logger),
				debug.Config(cfg),
			)
			if err != nil {
				logger.Info().Err(err).Str("transport", "debug").Msg("Failed to initialize server")
				return err
			}

			gr.Add(server.ListenAndServe, func(_ error) {
				_ = server.Shutdown(ctx)
				cancel()
			})

			httpServer, err := http.Server(
				http.Logger(logger),
				http.Context(ctx),
				http.Config(cfg),
				http.Metrics(m),
				http.Namespace(cfg.HTTP.Namespace),
				http.TraceProvider(traceProvider),
			)
			if err != nil {
				logger.Info().
					Err(err).
					Str("transport", "http").
					Msg("Failed to initialize server")

				return err
			}

			gr.Add(httpServer.Run, func(_ error) {
				logger.Error().
					Err(err).
					Str("server", "http").
					Msg("Shutting down server")
				cancel()
			})

			return gr.Run()
		},
	}
}
