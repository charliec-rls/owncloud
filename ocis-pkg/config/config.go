package config

import (
	"github.com/owncloud/ocis/ocis-pkg/shared"

	accounts "github.com/owncloud/ocis/extensions/accounts/pkg/config"
	appProvider "github.com/owncloud/ocis/extensions/app-provider/pkg/config"
	appRegistry "github.com/owncloud/ocis/extensions/app-registry/pkg/config"
	audit "github.com/owncloud/ocis/extensions/audit/pkg/config"
	authbasic "github.com/owncloud/ocis/extensions/auth-basic/pkg/config"
	authbearer "github.com/owncloud/ocis/extensions/auth-bearer/pkg/config"
	authmachine "github.com/owncloud/ocis/extensions/auth-machine/pkg/config"
	frontend "github.com/owncloud/ocis/extensions/frontend/pkg/config"
	gateway "github.com/owncloud/ocis/extensions/gateway/pkg/config"
	glauth "github.com/owncloud/ocis/extensions/glauth/pkg/config"
	graphExplorer "github.com/owncloud/ocis/extensions/graph-explorer/pkg/config"
	graph "github.com/owncloud/ocis/extensions/graph/pkg/config"
	group "github.com/owncloud/ocis/extensions/group/pkg/config"
	idm "github.com/owncloud/ocis/extensions/idm/pkg/config"
	idp "github.com/owncloud/ocis/extensions/idp/pkg/config"
	nats "github.com/owncloud/ocis/extensions/nats/pkg/config"
	notifications "github.com/owncloud/ocis/extensions/notifications/pkg/config"
	ocdav "github.com/owncloud/ocis/extensions/ocdav/pkg/config"
	ocs "github.com/owncloud/ocis/extensions/ocs/pkg/config"
	proxy "github.com/owncloud/ocis/extensions/proxy/pkg/config"
	search "github.com/owncloud/ocis/extensions/search/pkg/config"
	settings "github.com/owncloud/ocis/extensions/settings/pkg/config"
	sharing "github.com/owncloud/ocis/extensions/sharing/pkg/config"
	storagepublic "github.com/owncloud/ocis/extensions/storage-publiclink/pkg/config"
	storageshares "github.com/owncloud/ocis/extensions/storage-shares/pkg/config"
	storagesystem "github.com/owncloud/ocis/extensions/storage-system/pkg/config"
	storageusers "github.com/owncloud/ocis/extensions/storage-users/pkg/config"
	store "github.com/owncloud/ocis/extensions/store/pkg/config"
	thumbnails "github.com/owncloud/ocis/extensions/thumbnails/pkg/config"
	user "github.com/owncloud/ocis/extensions/user/pkg/config"
	web "github.com/owncloud/ocis/extensions/web/pkg/config"
	webdav "github.com/owncloud/ocis/extensions/webdav/pkg/config"
)

const (
	// SUPERVISED sets the runtime mode as supervised threads.
	SUPERVISED = iota

	// UNSUPERVISED sets the runtime mode as a single thread.
	UNSUPERVISED
)

type Mode int

// Runtime configures the oCIS runtime when running in supervised mode.
type Runtime struct {
	Port       string `yaml:"port" env:"OCIS_RUNTIME_PORT"`
	Host       string `yaml:"host" env:"OCIS_RUNTIME_HOST"`
	Extensions string `yaml:"extensions" env:"OCIS_RUN_EXTENSIONS"`
}

// Config combines all available configuration parts.
type Config struct {
	*shared.Commons `yaml:"shared"`

	Tracing *shared.Tracing `yaml:"tracing"`
	Log     *shared.Log     `yaml:"log"`

	Mode    Mode // DEPRECATED
	File    string
	OcisURL string `yaml:"ocis_url"`

	Registry          string               `yaml:"registry"`
	TokenManager      *shared.TokenManager `yaml:"token_manager"`
	MachineAuthAPIKey string               `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY"`
	TransferSecret    string               `yaml:"transfer_secret" env:"STORAGE_TRANSFER_SECRET"`
	MetadataUserID    string               `yaml:"metadata_user_id" env:"METADATA_USER_ID"`
	SystemUserAPIKey  string               `yaml:"system_user_api_key" env:"OCIS_SYSTEM_USER_API_KEY"`
	Runtime           Runtime              `yaml:"runtime"`

	Accounts          *accounts.Config      `yaml:"accounts"`
	AppProvider       *appProvider.Config   `yaml:"app-provider"`
	AppRegistry       *appRegistry.Config   `yaml:"app-registry"`
	Audit             *audit.Config         `yaml:"audit"`
	AuthBasic         *authbasic.Config     `yaml:"auth-basic"`
	AuthBearer        *authbearer.Config    `yaml:"auth-bearer"`
	AuthMachine       *authmachine.Config   `yaml:"auth-machine"`
	Frontend          *frontend.Config      `yaml:"frontend"`
	Gateway           *gateway.Config       `yaml:"gateway"`
	GLAuth            *glauth.Config        `yaml:"glauth"`
	Graph             *graph.Config         `yaml:"graph"`
	GraphExplorer     *graphExplorer.Config `yaml:"graph-explorer"`
	Group             *group.Config         `yaml:"group"`
	IDM               *idm.Config           `yaml:"idm"`
	IDP               *idp.Config           `yaml:"idp"`
	Nats              *nats.Config          `yaml:"nats"`
	Notifications     *notifications.Config `yaml:"notifications"`
	OCDav             *ocdav.Config         `yaml:"ocdav"`
	OCS               *ocs.Config           `yaml:"ocs"`
	Proxy             *proxy.Config         `yaml:"proxy"`
	Settings          *settings.Config      `yaml:"settings"`
	Sharing           *sharing.Config       `yaml:"sharing"`
	StorageSystem     *storagesystem.Config `yaml:"storage-system"`
	StoragePublicLink *storagepublic.Config `yaml:"storage-public"`
	StorageShares     *storageshares.Config `yaml:"storage-shares"`
	StorageUsers      *storageusers.Config  `yaml:"storage-users"`
	Store             *store.Config         `yaml:"store"`
	Thumbnails        *thumbnails.Config    `yaml:"thumbnails"`
	User              *user.Config          `yaml:"user"`
	Web               *web.Config           `yaml:"web"`
	WebDAV            *webdav.Config        `yaml:"webdav"`
	Search            *search.Config        `yaml:"search"`
}
