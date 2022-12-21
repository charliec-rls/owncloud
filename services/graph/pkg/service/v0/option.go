package svc

import (
	"net/http"

	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/ocis-pkg/roles"
	searchsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/search/v0"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/graph/pkg/config"
	"github.com/owncloud/ocis/v2/services/graph/pkg/identity"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger                   log.Logger
	Config                   *config.Config
	Middleware               []func(http.Handler) http.Handler
	RequireAdminMiddleware   func(http.Handler) http.Handler
	GatewayClient            GatewayClient
	IdentityBackend          identity.Backend
	IdentityEducationBackend identity.EducationBackend
	RoleService              RoleService
	PermissionService        Permissions
	RoleManager              *roles.Manager
	EventsPublisher          events.Publisher
	SearchService            searchsvc.SearchProviderService
}

// newOptions initializes the available default options.
func newOptions(opts ...Option) Options {
	opt := Options{}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Logger provides a function to set the logger option.
func Logger(val log.Logger) Option {
	return func(o *Options) {
		o.Logger = val
	}
}

// Config provides a function to set the config option.
func Config(val *config.Config) Option {
	return func(o *Options) {
		o.Config = val
	}
}

// Middleware provides a function to set the middleware option.
func Middleware(val ...func(http.Handler) http.Handler) Option {
	return func(o *Options) {
		o.Middleware = val
	}
}

// WithRequireAdminMiddleware provides a function to set the RequireAdminMiddleware option.
func WithRequireAdminMiddleware(val func(http.Handler) http.Handler) Option {
	return func(o *Options) {
		o.RequireAdminMiddleware = val
	}
}

// WithGatewayClient provides a function to set the gateway client option.
func WithGatewayClient(val GatewayClient) Option {
	return func(o *Options) {
		o.GatewayClient = val
	}
}

// WithIdentityBackend provides a function to set the IdentityBackend option.
func WithIdentityBackend(val identity.Backend) Option {
	return func(o *Options) {
		o.IdentityBackend = val
	}
}

// WithIdentityEducationBackend provides a function to set the IdentityEducationBackend option.
func WithIdentityEducationBackend(val identity.EducationBackend) Option {
	return func(o *Options) {
		o.IdentityEducationBackend = val
	}
}

// WithRoleService provides a function to set the RoleService option.
func WithRoleService(val RoleService) Option {
	return func(o *Options) {
		o.RoleService = val
	}
}

// WithSearchService provides a function to set the SearchService option.
func WithSearchService(val searchsvc.SearchProviderService) Option {
	return func(o *Options) {
		o.SearchService = val
	}
}

// PermissionService provides a function to set the PermissionService option.
func PermissionService(val settingssvc.PermissionService) Option {
	return func(o *Options) {
		o.PermissionService = val
	}
}

// RoleManager provides a function to set the RoleManager option.
func RoleManager(val *roles.Manager) Option {
	return func(o *Options) {
		o.RoleManager = val
	}
}

// EventsPublisher provides a function to set the EventsPublisher option.
func EventsPublisher(val events.Publisher) Option {
	return func(o *Options) {
		o.EventsPublisher = val
	}
}
