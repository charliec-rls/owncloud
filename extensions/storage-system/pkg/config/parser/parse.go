package parser

import (
	"errors"

	"github.com/owncloud/ocis/extensions/storage-system/pkg/config"
	"github.com/owncloud/ocis/extensions/storage-system/pkg/config/defaults"
	ociscfg "github.com/owncloud/ocis/ocis-pkg/config"
	"github.com/owncloud/ocis/ocis-pkg/shared"

	"github.com/owncloud/ocis/ocis-pkg/config/envdecode"
)

// ParseConfig loads configuration from known paths.
func ParseConfig(cfg *config.Config) error {
	_, err := ociscfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	defaults.EnsureDefaults(cfg)

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	defaults.Sanitize(cfg)

	return Validate(cfg)
}

func Validate(cfg *config.Config) error {
	if cfg.TokenManager.JWTSecret == "" {
		return shared.MissingJWTTokenError(cfg.Service.Name)
	}

	if cfg.SystemUserAPIKey == "" {
		return shared.MissingMachineAuthApiKeyError(cfg.Service.Name)
	}

	if cfg.SystemUserID == "" {
		return shared.MissingSystemUserID(cfg.Service.Name)
	}
	return nil
}
