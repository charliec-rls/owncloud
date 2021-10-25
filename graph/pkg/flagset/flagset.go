package flagset

import (
	"github.com/owncloud/ocis/graph/pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/urfave/cli/v2"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       flags.OverrideDefaultString(cfg.File, ""),
			Usage:       "Path to config file",
			EnvVars:     []string{"GRAPH_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"GRAPH_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"GRAPH_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"GRAPH_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9124"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"GRAPH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-file",
			Usage:       "Enable log to file",
			EnvVars:     []string{"GRAPH_LOG_FILE", "OCIS_LOG_FILE"},
			Destination: &cfg.Log.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"GRAPH_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"GRAPH_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"GRAPH_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"GRAPH_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "graph"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"GRAPH_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9124"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"GRAPH_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"GRAPH_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"GRAPH_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"GRAPH_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "127.0.0.1:9120"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"GRAPH_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/graph"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"GRAPH_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http service for service discovery",
			EnvVars:     []string{"GRAPH_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},

		&cli.StringFlag{
			Name:        "spaces-webdav-base",
			Value:       flags.OverrideDefaultString(cfg.Spaces.WebDavBase, "https://localhost:9200"),
			Usage:       "spaces webdav base URL to use when rendering drive WabDAV URLs",
			EnvVars:     []string{"GRAPH_SPACES_WEBDAV_BASE", "OCIS_URL"},
			Destination: &cfg.Spaces.WebDavBase,
		},
		&cli.StringFlag{
			Name:        "spaces-webdav-path",
			Value:       flags.OverrideDefaultString(cfg.Spaces.WebDavPath, "/dav/spaces/"),
			Usage:       "spaces webdav path to use when rendering drive WabDAV URLs",
			EnvVars:     []string{"GRAPH_SPACES_WEBDAV_PATH"},
			Destination: &cfg.Spaces.WebDavPath,
		},

		&cli.StringFlag{
			Name:        "default-space-quota",
			Value:       flags.OverrideDefaultString(cfg.Spaces.DefaultQuota, "1000000000"),
			Usage:       "default quota used for all spaces if no custom quota was given",
			EnvVars:     []string{"GRAPH_SPACES_DEFAULT_QUOTA"},
			Destination: &cfg.Spaces.DefaultQuota,
		},

		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       flags.OverrideDefaultString(cfg.TokenManager.JWTSecret, "Pive-Fumkiu4"),
			Usage:       "Used to validate the reva access JWT, should equal reva's jwt-secret",
			EnvVars:     []string{"GRAPH_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		&cli.StringFlag{
			Name:        "reva-gateway-addr",
			Value:       flags.OverrideDefaultString(cfg.Reva.Address, "127.0.0.1:9142"),
			Usage:       "Address of REVA gateway endpoint",
			EnvVars:     []string{"REVA_GATEWAY"},
			Destination: &cfg.Reva.Address,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
		},
	}
}
