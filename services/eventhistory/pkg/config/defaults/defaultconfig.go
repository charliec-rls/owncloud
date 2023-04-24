package defaults

import (
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/structs"
	"github.com/owncloud/ocis/v2/services/eventhistory/pkg/config"
)

// FullDefaultConfig returns the full default config
func FullDefaultConfig() *config.Config {
	cfg := DefaultConfig()
	EnsureDefaults(cfg)
	Sanitize(cfg)
	return cfg
}

// DefaultConfig return the default configuration
func DefaultConfig() *config.Config {
	return &config.Config{
		Service: config.Service{
			Name: "eventhistory",
		},
		Events: config.Events{
			Endpoint:  "127.0.0.1:9233",
			Cluster:   "ocis-cluster",
			EnableTLS: false,
		},
		Store: config.Store{
			Store:    "memory",
			Database: "eventhistory",
			Table:    "events",
			TTL:      336 * time.Hour,
		},
		GRPC: config.GRPCConfig{
			Addr:      "127.0.0.1:0",
			Namespace: "com.owncloud.api",
		},
	}
}

// EnsureDefaults ensures the config contains default values
func EnsureDefaults(cfg *config.Config) {
	// provide with defaults for shared logging, since we need a valid destination address for "envdecode".
	if cfg.Log == nil && cfg.Commons != nil && cfg.Commons.Log != nil {
		cfg.Log = &config.Log{
			Level:  cfg.Commons.Log.Level,
			Pretty: cfg.Commons.Log.Pretty,
			Color:  cfg.Commons.Log.Color,
			File:   cfg.Commons.Log.File,
		}
	} else if cfg.Log == nil {
		cfg.Log = &config.Log{}
	}

	if cfg.GRPCClientTLS == nil && cfg.Commons != nil {
		cfg.GRPCClientTLS = structs.CopyOrZeroValue(cfg.Commons.GRPCClientTLS)
	}

	if cfg.GRPC.TLS == nil && cfg.Commons != nil {
		cfg.GRPC.TLS = structs.CopyOrZeroValue(cfg.Commons.GRPCServiceTLS)
	}
}

// Sanitize sanitizes the config
func Sanitize(cfg *config.Config) {
	// nothing to sanitize here atm
}
