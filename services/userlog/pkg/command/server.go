package command

import (
	"context"
	"fmt"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	ogrpc "github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/store"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/logging"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/metrics"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/server/http"
	"github.com/urfave/cli/v2"
	microstore "go-micro.dev/v4/store"
)

// all events we care about
var _registeredEvents = []events.Unmarshaller{
	// space related
	events.SpaceDisabled{},
	events.SpaceDeleted{},
	events.SpaceShared{},
	events.SpaceUnshared{},
	events.SpaceMembershipExpired{},

	// share related
	events.ShareCreated{},
	events.ShareRemoved{},
	events.ShareExpired{},
}

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
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			err := ogrpc.Configure(ogrpc.GetClientOptions(cfg.GRPCClientTLS)...)
			if err != nil {
				return err
			}

			gr := run.Group{}
			ctx, cancel := func() (context.Context, context.CancelFunc) {
				if cfg.Context == nil {
					return context.WithCancel(context.Background())
				}
				return context.WithCancel(cfg.Context)
			}()

			mtrcs := metrics.New()
			mtrcs.BuildInfo.WithLabelValues(version.GetString()).Set(1)

			defer cancel()

			consumer, err := stream.NatsFromConfig(stream.NatsConfig(cfg.Events))
			if err != nil {
				return err
			}

			st := store.Create(
				store.WithCacheOptions(store.CacheOptions{
					Type: cfg.Store.Type,
					TTL:  cfg.Store.TTL,
					Size: cfg.Store.Size,
				}),
				microstore.Nodes(cfg.Store.Addresses...),
				microstore.Database(cfg.Store.Database),
				microstore.Table(cfg.Store.Table),
			)

			tm, err := pool.StringToTLSMode(cfg.GRPCClientTLS.Mode)
			if err != nil {
				return err
			}
			gwclient, err := pool.GetGatewayServiceClient(
				cfg.RevaGateway,
				pool.WithTLSCACert(cfg.GRPCClientTLS.CACert),
				pool.WithTLSMode(tm),
			)
			if err != nil {
				return fmt.Errorf("could not get reva client: %s", err)
			}

			hClient := ehsvc.NewEventHistoryService("com.owncloud.api.eventhistory", grpc.DefaultClient())

			{
				server, err := http.Server(
					http.Logger(logger),
					http.Context(ctx),
					http.Config(cfg),
					http.Metrics(mtrcs),
					http.Store(st),
					http.Consumer(consumer),
					http.Gateway(gwclient),
					http.History(hClient),
					http.RegisteredEvents(_registeredEvents),
				)

				if err != nil {
					logger.Info().Err(err).Str("transport", "http").Msg("Failed to initialize server")
					return err
				}

				gr.Add(func() error {
					return server.Run()
				}, func(err error) {
					logger.Error().
						Str("transport", "http").
						Err(err).
						Msg("Shutting down server")

					cancel()
				})
			}

			return gr.Run()
		},
	}
}
