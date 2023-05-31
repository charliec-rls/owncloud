package svc

import (
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"net/http"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/config"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/thumbnail/imgsource"
	"github.com/owncloud/ocis/v2/services/thumbnails/pkg/thumbnail/storage"
)

// Option defines a single option function.
type Option func(o *Options)

// Options defines the available options for this package.
type Options struct {
	Logger             log.Logger
	Config             *config.Config
	Middleware         []func(http.Handler) http.Handler
	ThumbnailStorage   storage.Storage
	ImageSource        imgsource.Source
	CS3Source          imgsource.Source
	CS3GatewaySelector pool.Selectable[gateway.GatewayAPIClient]
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

// ThumbnailStorage provides a function to set the thumbnail storage option.
func ThumbnailStorage(val storage.Storage) Option {
	return func(o *Options) {
		o.ThumbnailStorage = val
	}
}

// ThumbnailSource provides a function to set the image source option.
func ThumbnailSource(val imgsource.Source) Option {
	return func(o *Options) {
		o.ImageSource = val
	}
}

func CS3Source(val imgsource.Source) Option {
	return func(o *Options) {
		o.CS3Source = val
	}
}

func CS3ClientSelector(gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) Option {
	return func(o *Options) {
		o.CS3GatewaySelector = gatewaySelector
	}
}
