package registry

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	mRegistry "go-micro.dev/v4/registry"
	"go-micro.dev/v4/util/addr"
)

func BuildGRPCService(serviceID, uuid, address string, version string, nodeMetadata map[string]string) *mRegistry.Service {
	var host string
	var port int

	parts := strings.Split(address, ":")
	if len(parts) > 1 {
		host = strings.Join(parts[:len(parts)-1], ":")
		port, _ = strconv.Atoi(parts[len(parts)-1])
	} else {
		host = parts[0]
	}

	addr, err := addr.Extract(host)
	if err != nil {
		addr = host
	}

	node := &mRegistry.Node{
		Id:       serviceID + "-" + uuid,
		Address:  net.JoinHostPort(addr, fmt.Sprint(port)),
		Metadata: make(map[string]string),
	}

	node.Metadata["registry"] = GetRegistry().String()
	node.Metadata["server"] = "grpc"
	node.Metadata["transport"] = "grpc"
	node.Metadata["protocol"] = "grpc"
	for k, v := range nodeMetadata {
		node.Metadata[k] = v
	}

	return &mRegistry.Service{
		Name:      serviceID,
		Version:   version,
		Nodes:     []*mRegistry.Node{node},
		Endpoints: make([]*mRegistry.Endpoint, 0),
	}
}

func BuildHTTPService(serviceID, uuid, address string, version string) *mRegistry.Service {
	var host string
	var port int

	parts := strings.Split(address, ":")
	if len(parts) > 1 {
		host = strings.Join(parts[:len(parts)-1], ":")
		port, _ = strconv.Atoi(parts[len(parts)-1])
	} else {
		host = parts[0]
	}

	addr, err := addr.Extract(host)
	if err != nil {
		addr = host
	}

	node := &mRegistry.Node{
		Id:       serviceID + "-" + uuid,
		Address:  net.JoinHostPort(addr, fmt.Sprint(port)),
		Metadata: make(map[string]string),
	}

	node.Metadata["registry"] = GetRegistry().String()
	node.Metadata["server"] = "http"
	node.Metadata["transport"] = "http"
	node.Metadata["protocol"] = "http"

	return &mRegistry.Service{
		Name:      serviceID,
		Version:   version,
		Nodes:     []*mRegistry.Node{node},
		Endpoints: make([]*mRegistry.Endpoint, 0),
	}
}
