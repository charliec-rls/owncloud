package config

import "github.com/owncloud/ocis/v2/ocis-pkg/shared"

// HTTP defines the available http configuration.
type HTTP struct {
	Addr      string                `yaml:"addr" env:"WEB_HTTP_ADDR" desc:"The bind address of the HTTP service."`
	TLS       shared.HTTPServiceTLS `yaml:"tls"`
	Namespace string                `yaml:"-"`
	Root      string                `yaml:"root" env:"WEB_HTTP_ROOT" desc:"Subdirectory that serves as the root for this HTTP service."`
	CacheTTL  int                   `yaml:"cache_ttl" env:"WEB_CACHE_TTL" desc:"Cache policy in seconds for ownCloud Web assets."`
}
