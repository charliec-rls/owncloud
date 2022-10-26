package http

import (
	"strings"
	"time"

	"github.com/owncloud/ocis/v2/ocis-pkg/broker"
	"github.com/owncloud/ocis/v2/ocis-pkg/registry"
	oregistry "github.com/owncloud/ocis/v2/ocis-pkg/registry"

	mhttps "github.com/go-micro/plugins/v4/server/http"
	"go-micro.dev/v4"
	"go-micro.dev/v4/server"
)

// Service simply wraps the go-micro web service.
type Service struct {
	micro.Service
}

// NewService initializes a new http service.
func NewService(registry registry.Registry, opts ...Option) (Service, error) {

	reg, err := oregistry.GetRegistry(registry)
	if err != nil {
		return Service{}, err
	}

	noopBroker := broker.NoOp{}
	sopts := newOptions(opts...)
	wopts := []micro.Option{
		micro.Server(mhttps.NewServer(server.TLSConfig(sopts.TLSConfig))),
		micro.Broker(noopBroker),
		micro.Address(sopts.Address),
		micro.Name(strings.Join([]string{sopts.Namespace, sopts.Name}, ".")),
		micro.Version(sopts.Version),
		micro.Context(sopts.Context),
		micro.Flags(sopts.Flags...),
		micro.Registry(reg),
		micro.RegisterTTL(time.Second * 30),
		micro.RegisterInterval(time.Second * 10),
	}

	return Service{micro.NewService(wopts...)}, nil
}
