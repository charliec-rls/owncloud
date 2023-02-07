package config

import (
	"context"

	"github.com/owncloud/ocis/v2/ocis-pkg/shared"
)

// Config combines all available configuration parts.
type Config struct {
	Commons *shared.Commons `yaml:"-"` // don't use this directly as configuration for a service

	Service Service `yaml:"-"`

	Tracing *Tracing `yaml:"tracing"`
	Log     *Log     `yaml:"log"`
	Debug   Debug    `yaml:"debug"`

	HTTP HTTP `yaml:"http"`

	Asset Asset  `yaml:"asset"`
	File  string `yaml:"file" env:"WEB_UI_CONFIG" desc:"Read the ownCloud Web configuration from this file."` // TODO: rename this to a more self explaining string
	Web   Web    `yaml:"web"`

	Context context.Context `yaml:"-"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"path" env:"WEB_ASSET_PATH" desc:"Serve ownCloud Web assets from a path on the filesystem instead of the builtin assets."`
}

// WebConfig defines the available web configuration for a dynamically rendered config.json.
type WebConfig struct {
	Server        string                 `json:"server,omitempty" yaml:"server" env:"OCIS_URL;WEB_UI_CONFIG_SERVER" desc:"URL, where the oCIS APIs are reachable for ownCloud Web."`
	Theme         string                 `json:"theme,omitempty" yaml:"-"`
	OpenIDConnect OIDC                   `json:"openIdConnect,omitempty" yaml:"oidc"`
	Apps          []string               `json:"apps" yaml:"apps"`
	Applications  []Application          `json:"applications,omitempty" yaml:"applications"`
	ExternalApps  []ExternalApp          `json:"external_apps,omitempty" yaml:"external_apps"`
	Options       map[string]interface{} `json:"options,omitempty" yaml:"options"`
}

// OIDC defines the available oidc configuration
type OIDC struct {
	MetadataURL  string `json:"metadata_url,omitempty" yaml:"metadata_url" env:"WEB_OIDC_METADATA_URL" desc:"URL for the OIDC well-known configuration endpoint. Defaults to the oCIS API URL + \"/.well-known/openid-configuration\"."`
	Authority    string `json:"authority,omitempty" yaml:"authority" env:"OCIS_URL;OCIS_OIDC_ISSUER;WEB_OIDC_AUTHORITY" desc:"URL of the OIDC issuer. It defaults to URL of the builtin IDP."`
	ClientID     string `json:"client_id,omitempty" yaml:"client_id" env:"WEB_OIDC_CLIENT_ID" desc:"OIDC client ID, which ownCloud Web uses. This client needs to be set up in your IDP."`
	ResponseType string `json:"response_type,omitempty" yaml:"response_type" env:"WEB_OIDC_RESPONSE_TYPE" desc:"OIDC response type to use for authentication."`
	Scope        string `json:"scope,omitempty" yaml:"scope" env:"WEB_OIDC_SCOPE" desc:"OIDC scopes to request during authentication to authorize access to user details. Defaults to 'openid profile email'. Values are separated by blank. More example values but not limited to are 'address' or 'phone' etc."`
}

// Application defines an application for the Web app switcher.
type Application struct {
	Icon   string            `json:"icon,omitempty" yaml:"icon"`
	Target string            `json:"target,omitempty" yaml:"target"`
	Title  map[string]string `json:"title,omitempty" yaml:"title"`
	Menu   string            `json:"menu,omitempty" yaml:"menu"`
	URL    string            `json:"url,omitempty" yaml:"url"`
}

// ExternalApp defines an external web app.
// {
//	"name": "hello",
//	"path": "http://localhost:9105/hello.js",
//	  "config": {
//	    "url": "http://localhost:9105"
//	  }
//  }
type ExternalApp struct {
	ID   string `json:"id,omitempty" yaml:"id"`
	Path string `json:"path,omitempty" yaml:"path"`
	// Config is completely dynamic, because it depends on the extension
	Config map[string]interface{} `json:"config,omitempty" yaml:"config"`
}

// ExternalAppConfig defines an external web app configuration.
type ExternalAppConfig struct {
	URL string `json:"url,omitempty" yaml:"url" env:""`
}

// Web defines the available web configuration.
type Web struct {
	Path        string    `yaml:"path" env:"WEB_UI_PATH" desc:"Read the ownCloud Web configuration from this file path."`
	ThemeServer string    `yaml:"theme_server" env:"OCIS_URL;WEB_UI_THEME_SERVER" desc:"URL to load themes from. Will be prepended to the theme path."` // used to build Theme in WebConfig
	ThemePath   string    `yaml:"theme_path" env:"WEB_UI_THEME_PATH" desc:"URL path to load themes from. The theme server will be prepended."`          // used to build Theme in WebConfig
	Config      WebConfig `yaml:"config"`
}
