package flagset

import (
	"path"

	"github.com/owncloud/ocis/ocis-pkg/config/defaults"
	"github.com/owncloud/ocis/ocis-pkg/flags"
	"github.com/owncloud/ocis/settings/pkg/config"
	"github.com/urfave/cli/v2"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set logging level",
			EnvVars:     []string{"SETTINGS_LOG_LEVEL", "OCIS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"SETTINGS_LOG_PRETTY", "OCIS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"SETTINGS_LOG_COLOR", "OCIS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// HealthWithConfig applies cfg to the root flagset
func HealthWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9194"),
			Usage:       "Address to debug endpoint",
			EnvVars:     []string{"SETTINGS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Value:       "",
			Usage:       "Path to config file",
			EnvVars:     []string{"SETTINGS_CONFIG_FILE"},
			Destination: &cfg.File,
		},
		&cli.BoolFlag{
			Name:        "tracing-enabled",
			Usage:       "Enable sending traces",
			EnvVars:     []string{"SETTINGS_TRACING_ENABLED", "OCIS_TRACING_ENABLED"},
			Destination: &cfg.Tracing.Enabled,
		},
		&cli.StringFlag{
			Name:        "tracing-type",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Type, "jaeger"),
			Usage:       "Tracing backend type",
			EnvVars:     []string{"SETTINGS_TRACING_TYPE", "OCIS_TRACING_TYPE"},
			Destination: &cfg.Tracing.Type,
		},
		&cli.StringFlag{
			Name:        "tracing-endpoint",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Endpoint, ""),
			Usage:       "Endpoint for the agent",
			EnvVars:     []string{"SETTINGS_TRACING_ENDPOINT", "OCIS_TRACING_ENDPOINT"},
			Destination: &cfg.Tracing.Endpoint,
		},
		&cli.StringFlag{
			Name:        "tracing-collector",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Collector, ""),
			Usage:       "Endpoint for the collector",
			EnvVars:     []string{"SETTINGS_TRACING_COLLECTOR", "OCIS_TRACING_COLLECTOR"},
			Destination: &cfg.Tracing.Collector,
		},
		&cli.StringFlag{
			Name:        "tracing-service",
			Value:       flags.OverrideDefaultString(cfg.Tracing.Service, "settings"),
			Usage:       "Service name for tracing",
			EnvVars:     []string{"SETTINGS_TRACING_SERVICE"},
			Destination: &cfg.Tracing.Service,
		},
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       flags.OverrideDefaultString(cfg.Debug.Addr, "127.0.0.1:9194"),
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"SETTINGS_DEBUG_ADDR"},
			Destination: &cfg.Debug.Addr,
		},
		&cli.StringFlag{
			Name:        "debug-token",
			Value:       flags.OverrideDefaultString(cfg.Debug.Token, ""),
			Usage:       "Token to grant metrics access",
			EnvVars:     []string{"SETTINGS_DEBUG_TOKEN"},
			Destination: &cfg.Debug.Token,
		},
		&cli.BoolFlag{
			Name:        "debug-pprof",
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"SETTINGS_DEBUG_PPROF"},
			Destination: &cfg.Debug.Pprof,
		},
		&cli.BoolFlag{
			Name:        "debug-zpages",
			Usage:       "Enable zpages debugging",
			EnvVars:     []string{"SETTINGS_DEBUG_ZPAGES"},
			Destination: &cfg.Debug.Zpages,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Addr, "127.0.0.1:9190"),
			Usage:       "Address to bind http server",
			EnvVars:     []string{"SETTINGS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Namespace, "com.owncloud.web"),
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"SETTINGS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       flags.OverrideDefaultString(cfg.HTTP.Root, "/"),
			Usage:       "Root path of http server",
			EnvVars:     []string{"SETTINGS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.IntFlag{
			Name:        "http-cache-ttl",
			Value:       flags.OverrideDefaultInt(cfg.HTTP.CacheTTL, 604800), // 10 days
			Usage:       "Set the static assets caching duration in seconds",
			EnvVars:     []string{"SETTINGS_CACHE_TTL"},
			Destination: &cfg.HTTP.CacheTTL,
		},
		&cli.StringSliceFlag{
			Name:    "cors-allowed-origins",
			Value:   cli.NewStringSlice("*"),
			Usage:   "Set the allowed CORS origins",
			EnvVars: []string{"SETTINGS_CORS_ALLOW_ORIGINS", "OCIS_CORS_ALLOW_ORIGINS"},
		},
		&cli.StringSliceFlag{
			Name:    "cors-allowed-methods",
			Value:   cli.NewStringSlice("GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"),
			Usage:   "Set the allowed CORS origins",
			EnvVars: []string{"SETTINGS_CORS_ALLOW_METHODS", "OCIS_CORS_ALLOW_METHODS"},
		},
		&cli.StringSliceFlag{
			Name:    "cors-allowed-headers",
			Value:   cli.NewStringSlice("Authorization", "Origin", "Content-Type", "Accept", "X-Requested-With"),
			Usage:   "Set the allowed CORS origins",
			EnvVars: []string{"SETTINGS_CORS_ALLOW_HEADERS", "OCIS_CORS_ALLOW_HEADERS"},
		},
		&cli.BoolFlag{
			Name:    "cors-allow-credentials",
			Value:   flags.OverrideDefaultBool(cfg.HTTP.CORS.AllowCredentials, true),
			Usage:   "Allow credentials for CORS",
			EnvVars: []string{"SETTINGS_CORS_ALLOW_CREDENTIALS", "OCIS_CORS_ALLOW_CREDENTIALS"},
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       flags.OverrideDefaultString(cfg.GRPC.Addr, "127.0.0.1:9191"),
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"SETTINGS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       flags.OverrideDefaultString(cfg.Asset.Path, ""),
			Usage:       "Path to custom assets",
			EnvVars:     []string{"SETTINGS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"SETTINGS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "settings"),
			Usage:       "service name",
			EnvVars:     []string{"SETTINGS_NAME"},
			Destination: &cfg.Service.Name,
		},
		&cli.StringFlag{
			Name:        "data-path",
			Value:       flags.OverrideDefaultString(cfg.Service.DataPath, path.Join(defaults.BaseDataPath(), "settings")),
			Usage:       "Mount path for the storage",
			EnvVars:     []string{"SETTINGS_DATA_PATH"},
			Destination: &cfg.Service.DataPath,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       flags.OverrideDefaultString(cfg.TokenManager.JWTSecret, "Pive-Fumkiu4"),
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"SETTINGS_JWT_SECRET", "OCIS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
		&cli.StringFlag{
			Name:  "extensions",
			Usage: "Run specific extensions during supervised mode. This flag is set by the runtime",
		},
	}
}

// ListSettingsWithConfig applies list command flags to cfg
func ListSettingsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       flags.OverrideDefaultString(cfg.GRPC.Namespace, "com.owncloud.api"),
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"SETTINGS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       flags.OverrideDefaultString(cfg.Service.Name, "settings"),
			Usage:       "service name",
			EnvVars:     []string{"SETTINGS_NAME"},
			Destination: &cfg.Service.Name,
		},
	}
}
