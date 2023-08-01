package registry

import (
	"os"
	"strings"
	"time"

	consulr "github.com/go-micro/plugins/v4/registry/consul"
	etcdr "github.com/go-micro/plugins/v4/registry/etcd"
	kubernetesr "github.com/go-micro/plugins/v4/registry/kubernetes"
	mdnsr "github.com/go-micro/plugins/v4/registry/mdns"
	memr "github.com/go-micro/plugins/v4/registry/memory"
	natsr "github.com/go-micro/plugins/v4/registry/nats"

	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/registry/cache"
)

var (
	registryEnv        = "MICRO_REGISTRY"
	registryAddressEnv = "MICRO_REGISTRY_ADDRESS"
)

// GetRegistry returns a configured micro registry based on Micro env vars.
// It defaults to mDNS, so mind that systems with mDNS disabled by default (i.e SUSE) will have a hard time
// and it needs to explicitly use etcd. Os awareness for providing a working registry out of the box should be done.
func GetRegistry() registry.Registry {
	addresses := strings.Split(os.Getenv(registryAddressEnv), ",")

	var r registry.Registry
	switch os.Getenv(registryEnv) {
	case "nats":
		r = natsr.NewRegistry(
			registry.Addrs(addresses...),
		)
	case "kubernetes":
		r = kubernetesr.NewRegistry(
			registry.Addrs(addresses...),
		)
	case "etcd":
		r = etcdr.NewRegistry(
			registry.Addrs(addresses...),
		)
	case "consul":
		r = consulr.NewRegistry(
			registry.Addrs(addresses...),
		)
	case "memory":
		r = memr.NewRegistry()
	default:
		r = mdnsr.NewRegistry()
	}

	// always use cached registry to prevent registry
	// lookup for every request
	return cache.New(r, cache.WithTTL(20*time.Second))
}
