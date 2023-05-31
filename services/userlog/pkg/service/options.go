package service

import (
	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/go-chi/chi/v5"
	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	ehsvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/eventhistory/v0"
	"github.com/owncloud/ocis/v2/services/userlog/pkg/config"
	"go-micro.dev/v4/store"
)

// Option for the userlog service
type Option func(*Options)

// Options for the userlog service
type Options struct {
	Logger           log.Logger
	Consumer         events.Consumer
	Mux              *chi.Mux
	Store            store.Store
	Config           *config.Config
	HistoryClient    ehsvc.EventHistoryService
	GatewaySelector  pool.Selectable[gateway.GatewayAPIClient]
	RegisteredEvents []events.Unmarshaller
}

// Logger configures a logger for the userlog service
func Logger(log log.Logger) Option {
	return func(o *Options) {
		o.Logger = log
	}
}

// Consumer configures an event consumer for the userlog service
func Consumer(c events.Consumer) Option {
	return func(o *Options) {
		o.Consumer = c
	}
}

// Mux defines the muxer for the userlog service
func Mux(m *chi.Mux) Option {
	return func(o *Options) {
		o.Mux = m
	}
}

// Store defines the store for the userlog service
func Store(s store.Store) Option {
	return func(o *Options) {
		o.Store = s
	}
}

// Config adds the config for the userlog service
func Config(c *config.Config) Option {
	return func(o *Options) {
		o.Config = c
	}
}

// HistoryClient adds a grpc client for the eventhistory service
func HistoryClient(hc ehsvc.EventHistoryService) Option {
	return func(o *Options) {
		o.HistoryClient = hc
	}
}

// GatewaySelector adds a grpc client selector for the gateway service
func GatewaySelector(gatewaySelector pool.Selectable[gateway.GatewayAPIClient]) Option {
	return func(o *Options) {
		o.GatewaySelector = gatewaySelector
	}
}

// RegisteredEvents registers the events the service should listen to
func RegisteredEvents(e []events.Unmarshaller) Option {
	return func(o *Options) {
		o.RegisteredEvents = e
	}
}
