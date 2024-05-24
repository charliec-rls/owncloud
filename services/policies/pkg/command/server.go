package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/signal"

	"github.com/cs3org/reva/v2/pkg/events/stream"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/runner"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/debug"
	"github.com/owncloud/ocis/v2/ocis-pkg/service/grpc"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/ocis-pkg/version"
	svcProtogen "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/policies/v0"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config"
	"github.com/owncloud/ocis/v2/services/policies/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/policies/pkg/engine/opa"
	svcEvent "github.com/owncloud/ocis/v2/services/policies/pkg/service/event"
	svcGRPC "github.com/owncloud/ocis/v2/services/policies/pkg/service/grpc"
	"github.com/urfave/cli/v2"
)

// Server is the entrypoint for the server command.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "server",
		Usage:    fmt.Sprintf("start the %s service without runtime (unsupervised mode)", "authz"),
		Category: "server",
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			var cancel context.CancelFunc
			ctx := cfg.Context
			if ctx == nil {
				ctx, cancel = signal.NotifyContext(context.Background(), runner.StopSignals...)
				defer cancel()
			}

			logger := log.NewLogger(
				log.Name(cfg.Service.Name),
				log.Level(cfg.Log.Level),
				log.Pretty(cfg.Log.Pretty),
				log.Color(cfg.Log.Color),
				log.File(cfg.Log.File),
			).SubloggerWithRequestID(ctx)

			traceProvider, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}

			e, err := opa.NewOPA(cfg.Engine.Timeout, logger, cfg.Engine)
			if err != nil {
				return err
			}

			gr := runner.NewGroup()
			{
				grpcClient, err := grpc.NewClient(
					append(
						grpc.GetClientOptions(cfg.GRPCClientTLS),
						grpc.WithTraceProvider(traceProvider),
					)...,
				)
				if err != nil {
					return err
				}

				svc, err := grpc.NewServiceWithClient(
					grpcClient,
					grpc.Logger(logger),
					grpc.TLSEnabled(cfg.GRPC.TLS.Enabled),
					grpc.TLSCert(
						cfg.GRPC.TLS.Cert,
						cfg.GRPC.TLS.Key,
					),
					grpc.Name(cfg.Service.Name),
					grpc.Context(ctx),
					grpc.Address(cfg.GRPC.Addr),
					grpc.Namespace(cfg.GRPC.Namespace),
					grpc.Version(version.GetString()),
					grpc.TraceProvider(traceProvider),
				)
				if err != nil {
					return err
				}

				grpcSvc, err := svcGRPC.New(e)
				if err != nil {
					return err
				}

				if err := svcProtogen.RegisterPoliciesProviderHandler(
					svc.Server(),
					grpcSvc,
				); err != nil {
					return err
				}

				gr.Add(runner.NewGoMicroGrpcServerRunner("policies_grpc", svc))
			}

			{

				bus, err := stream.NatsFromConfig(cfg.Service.Name, false, stream.NatsConfig(cfg.Events))
				if err != nil {
					return err
				}

				eventSvc, err := svcEvent.New(ctx, bus, logger, traceProvider, e, cfg.Postprocessing.Query)
				if err != nil {
					return err
				}

				gr.Add(runner.New("policies_svc", func() error {
					return eventSvc.Run()
				}, func() {
					eventSvc.Close()
				}))
			}

			{
				server := debug.NewService(
					debug.Logger(logger),
					debug.Name(cfg.Service.Name),
					debug.Version(version.GetString()),
					debug.Address(cfg.Debug.Addr),
					debug.Token(cfg.Debug.Token),
					debug.Pprof(cfg.Debug.Pprof),
					debug.Zpages(cfg.Debug.Zpages),
					debug.Health(
						func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "text/plain")
							w.WriteHeader(http.StatusOK)

							// TODO: check if services are up and running

							_, err := io.WriteString(w, http.StatusText(http.StatusOK))
							// io.WriteString should not fail but if it does, we want to know.
							if err != nil {
								panic(err)
							}
						},
					),
					debug.Ready(
						func(w http.ResponseWriter, r *http.Request) {
							w.Header().Set("Content-Type", "text/plain")
							w.WriteHeader(http.StatusOK)

							// TODO: check if services are up and running

							_, err := io.WriteString(w, http.StatusText(http.StatusOK))
							// io.WriteString should not fail but if it does, we want to know.
							if err != nil {
								panic(err)
							}
						},
					),
				)

				gr.Add(runner.NewGolangHttpServerRunner("policies_debug", server))
			}

			grResults := gr.Run(ctx)

			// return the first non-nil error found in the results
			for _, grResult := range grResults {
				if grResult.RunnerError != nil {
					return grResult.RunnerError
				}
			}
			return nil
		},
	}
}
