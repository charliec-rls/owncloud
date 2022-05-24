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

	Reva              *Reva  `yaml:"reva"`
	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;IDP_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to impersonate users when looking up their userinfo via the 'cs3' backend."`

	Asset   Asset    `yaml:"asset"`
	IDP     Settings `yaml:"idp"`
	Clients []Client `yaml:"clients"`
	Ldap    Ldap     `yaml:"ldap"`

	Context context.Context `yaml:"-"`
}

// Ldap defines the available LDAP configuration.
type Ldap struct {
	URI       string `yaml:"uri" env:"LDAP_URI;IDP_LDAP_URI"`
	TLSCACert string `yaml:"cacert" env:"LDAP_CACERT;IDP_LDAP_TLS_CACERT"`

	BindDN       string `yaml:"bind_dn" env:"LDAP_BIND_DN;IDP_LDAP_BIND_DN"`
	BindPassword string `yaml:"bind_password" env:"LDAP_BIND_PASSWORD;IDP_LDAP_BIND_PASSWORD"`

	BaseDN string `yaml:"base_dn" env:"LDAP_USER_BASE_DN;IDP_LDAP_BASE_DN"`
	Scope  string `yaml:"scope" env:"LDAP_USER_SCOPE;IDP_LDAP_SCOPE"`

	LoginAttribute    string `yaml:"login_attribute" env:"IDP_LDAP_LOGIN_ATTRIBUTE"`
	EmailAttribute    string `yaml:"email_attribute" env:"LDAP_USER_SCHEMA_MAIL;IDP_LDAP_EMAIL_ATTRIBUTE"`
	NameAttribute     string `yaml:"name_attribute" env:"LDAP_USER_SCHEMA_USERNAME;IDP_LDAP_NAME_ATTRIBUTE"`
	UUIDAttribute     string `yaml:"uuid_attribute" env:"LDAP_USER_SCHEMA_ID;IDP_LDAP_UUID_ATTRIBUTE"`
	UUIDAttributeType string `yaml:"uuid_attribute_type" env:"IDP_LDAP_UUID_ATTRIBUTE_TYPE"`

	Filter      string `yaml:"filter" env:"LDAP_USER_FILTER;IDP_LDAP_FILTER"`
	ObjectClass string `yaml:"objectclass" env:"LDAP_USER_OBJECTCLASS;IDP_LDAP_OBJECTCLASS"`
}

// Asset defines the available asset configuration.
type Asset struct {
	Path string `yaml:"asset" env:"IDP_ASSET_PATH"`
}

type Client struct {
	ID              string   `yaml:"id"`
	Name            string   `yaml:"name"`
	Trusted         bool     `yaml:"trusted"`
	Secret          string   `yaml:"secret"`
	RedirectURIs    []string `yaml:"redirect_uris"`
	Origins         []string `yaml:"origins"`
	ApplicationType string   `yaml:"application_type"`
}

type Settings struct {
	// don't change the order of elements in this struct
	// it needs to match github.com/libregraph/lico/bootstrap.Settings

	Iss string `yaml:"iss" env:"OCIS_URL;OCIS_OIDC_ISSUER;IDP_ISS" desc:"The OIDC issuer URL to use."`

	IdentityManager string `yaml:"identity_manager" env:"IDP_IDENTITY_MANAGER" desc:"The identity manager implementation to use, defaults to 'ldap', can be changed to 'cs3', 'kc', 'libregraph', 'cookie' or 'guest'."`

	URIBasePath string `yaml:"uri_base_path" env:"IDP_URI_BASE_PATH"`

	SignInURI    string `yaml:"sign_in_uri" env:"IDP_SIGN_IN_URI"`
	SignedOutURI string `yaml:"signed_out_uri" env:"IDP_SIGN_OUT_URI"`

	AuthorizationEndpointURI string `yaml:"authorization_endpoint_uri" env:"IDP_ENDPOINT_URI"`
	EndsessionEndpointURI    string `yaml:"end_session_endpoint_uri" env:"IDP_ENDSESSION_ENDPOINT_URI"`

	Insecure bool `yaml:"insecure" env:"IDP_INSECURE" desc:"Allow insecure connections to the backend."`

	TrustedProxy []string `yaml:"trusted_proxy"` //TODO: how to configure this via env?

	AllowScope                     []string `yaml:"allow_scope"` // TODO: is this even needed?
	AllowClientGuests              bool     `yaml:"allow_client_guests" env:"IDP_ALLOW_CLIENT_GUESTS"`
	AllowDynamicClientRegistration bool     `yaml:"allow_dynamic_client_registration" env:"IDP_ALLOW_DYNAMIC_CLIENT_REGISTRATION"`

	EncryptionSecretFile string `yaml:"encrypt_secret_file" env:"IDP_ENCRYPTION_SECRET_FILE"`

	Listen string

	IdentifierClientDisabled          bool   `yaml:"identifier_client_disabled" env:"IDP_DISABLE_IDENTIFIER_WEBAPP"`
	IdentifierClientPath              string `yaml:"-"`
	IdentifierRegistrationConf        string `yaml:"-"`
	IdentifierScopesConf              string `yaml:"identifier_scopes_conf" env:"IDP_IDENTIFIER_SCOPES_CONF"`
	IdentifierDefaultBannerLogo       string
	IdentifierDefaultSignInPageText   string
	IdentifierDefaultUsernameHintText string
	IdentifierUILocales               []string

	SigningKid             string   `yaml:"signing_kid" env:"IDP_SIGNING_KID"`
	SigningMethod          string   `yaml:"signing_method" env:"IDP_SIGNING_METHOD"`
	SigningPrivateKeyFiles []string `yaml:"signing_private_key_files" env:"IDP_SIGNING_PRIVATE_KEY_FILES"`
	ValidationKeysPath     string   `yaml:"validation_keys_path" env:"IDP_VALIDATION_KEYS_PATH"`

	CookieBackendURI string
	CookieNames      []string

	AccessTokenDurationSeconds        uint64 `yaml:"access_token_duration_seconds" env:"IDP_ACCESS_TOKEN_EXPIRATION"`
	IDTokenDurationSeconds            uint64 `yaml:"id_token_duration_seconds" env:"IDP_ID_TOKEN_EXPIRATION"`
	RefreshTokenDurationSeconds       uint64 `yaml:"refresh_token_duration_seconds" env:"IDP_REFRESH_TOKEN_EXPIRATION"`
	DyamicClientSecretDurationSeconds uint64 `yaml:"dynamic_client_secret_duration_seconds" env:""`
}
