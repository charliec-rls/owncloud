package command

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/events/server"
	"github.com/go-micro/plugins/v4/events/natsjs"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	ociscrypto "github.com/owncloud/ocis/v2/ocis-pkg/crypto"
	"github.com/owncloud/ocis/v2/services/audit/pkg/config"
	"github.com/owncloud/ocis/v2/services/audit/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/audit/pkg/logging"
	svc "github.com/owncloud/ocis/v2/services/audit/pkg/service"
	"github.com/owncloud/ocis/v2/services/audit/pkg/types"
	"github.com/urfave/cli/v2"
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
			logger := logging.Configure(cfg.Service.Name, cfg.Log)

			ctx := cfg.Context
			if ctx == nil {
				ctx = context.Background()
			}
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			evtsCfg := cfg.Events

			var tlsConf *tls.Config
			if evtsCfg.EnableTLS {
				var rootCAPool *x509.CertPool
				if evtsCfg.TLSRootCACertificate != "" {
					rootCrtFile, err := os.Open(evtsCfg.TLSRootCACertificate)
					if err != nil {
						return err
					}

					rootCAPool, err = ociscrypto.NewCertPoolFromPEM(rootCrtFile)
					if err != nil {
						return err
					}
					evtsCfg.TLSInsecure = false
				}

				tlsConf = &tls.Config{
					MinVersion:         tls.VersionTLS12,
					InsecureSkipVerify: evtsCfg.TLSInsecure, //nolint:gosec
					RootCAs:            rootCAPool,
				}
			}
			client, err := server.NewNatsStream(
				natsjs.TLSConfig(tlsConf),
				natsjs.Address(evtsCfg.Endpoint),
				natsjs.ClusterID(evtsCfg.Cluster),
			)
			if err != nil {
				return err
			}
			evts, err := events.Consume(client, evtsCfg.ConsumerGroup, types.RegisteredEvents()...)
			if err != nil {
				return err
			}

			svc.AuditLoggerFromConfig(ctx, cfg.Auditlog, evts, logger)
			return nil
		},
	}
}
