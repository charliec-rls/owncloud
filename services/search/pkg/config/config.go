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

	GRPC GRPC `yaml:"grpc"`

	Datapath string `yaml:"data_path" env:"SEARCH_DATA_PATH" desc:"Path for the search persistence directory."`
	Reva     *Reva  `yaml:"reva"`
	Events   Events `yaml:"events"`

	MachineAuthAPIKey string `yaml:"machine_auth_api_key" env:"OCIS_MACHINE_AUTH_API_KEY;SEARCH_MACHINE_AUTH_API_KEY" desc:"Machine auth API key used to validate internal requests necessary for the access to resources from other services."`

	Context context.Context `yaml:"-"`
}

// Events combines the configuration options for the event bus.
type Events struct {
	Endpoint             string `yaml:"endpoint" env:"SEARCH_EVENTS_ENDPOINT" desc:"The address of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture."`
	Cluster              string `yaml:"cluster" env:"SEARCH_EVENTS_CLUSTER" desc:"The clusterID of the event system. The event system is the message queuing service. It is used as message broker for the microservice architecture. Mandatory when using NATS as event system."`
	ConsumerGroup        string `yaml:"group" env:"SEARCH_EVENTS_GROUP" desc:"The customer group of the service. One group will only get one copy of an event"`
	TLSInsecure          bool   `yaml:"tls_insecure" env:"OCIS_INSECURE;SEARCH_EVENTS_TLS_INSECURE" desc:"Whether to verify the server TLS certificates."`
	TLSRootCACertificate string `yaml:"tls_root_ca_certificate" env:"SEARCH_EVENTS_TLS_ROOT_CA_CERTIFICATE" desc:"The root CA certificate used to validate the server's TLS certificate. If provided SEARCH_EVENTS_TLS_INSECURE will be seen as false."`
}
