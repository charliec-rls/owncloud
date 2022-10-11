package revaconfig

import (
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
)

// StorageUsersConfigFromStruct will adapt an oCIS config struct into a reva mapstructure to start a reva service.
func StorageUsersConfigFromStruct(cfg *config.Config) map[string]interface{} {
	rcfg := map[string]interface{}{
		"core": map[string]interface{}{
			"tracing_enabled":      cfg.Tracing.Enabled,
			"tracing_endpoint":     cfg.Tracing.Endpoint,
			"tracing_collector":    cfg.Tracing.Collector,
			"tracing_service_name": cfg.Service.Name,
		},
		"shared": map[string]interface{}{
			"jwt_secret":                cfg.TokenManager.JWTSecret,
			"gatewaysvc":                cfg.Reva.Address,
			"skip_user_groups_in_token": cfg.SkipUserGroupsInToken,
		},
		"grpc": map[string]interface{}{
			"network": cfg.GRPC.Protocol,
			"address": cfg.GRPC.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"storageprovider": map[string]interface{}{
					"driver":             cfg.Driver,
					"drivers":            UserDrivers(cfg),
					"mount_id":           cfg.MountID,
					"expose_data_server": cfg.ExposeDataServer,
					"data_server_url":    cfg.DataServerURL,
					"upload_expiration":  cfg.UploadExpiration,
				},
			},
			"interceptors": map[string]interface{}{
				"eventsmiddleware": map[string]interface{}{
					"group":     "sharing",
					"type":      "nats",
					"address":   cfg.Events.Addr,
					"clusterID": cfg.Events.ClusterID,
				},
				"prometheus": map[string]interface{}{
					"namespace": "ocis",
					"subsystem": "storage_users",
				},
			},
		},
		"http": map[string]interface{}{
			"network": cfg.HTTP.Protocol,
			"address": cfg.HTTP.Addr,
			// TODO build services dynamically
			"services": map[string]interface{}{
				"dataprovider": map[string]interface{}{
					"prefix":         cfg.HTTP.Prefix,
					"driver":         cfg.Driver,
					"drivers":        UserDrivers(cfg),
					"nats_address":   cfg.Events.Addr,
					"nats_clusterID": cfg.Events.ClusterID,
					"data_txs": map[string]interface{}{
						"simple": map[string]interface{}{
							"cache_store":    cfg.Cache.Store,
							"cache_nodes":    cfg.Cache.Nodes,
							"cache_database": cfg.Cache.Database,
							"cache_table":    "stat",
						},
						"spaces": map[string]interface{}{
							"cache_store":    cfg.Cache.Store,
							"cache_nodes":    cfg.Cache.Nodes,
							"cache_database": cfg.Cache.Database,
							"cache_table":    "stat",
						},
						"tus": map[string]interface{}{
							"cache_store":    cfg.Cache.Store,
							"cache_nodes":    cfg.Cache.Nodes,
							"cache_database": cfg.Cache.Database,
							"cache_table":    "stat",
						},
					},
				},
			},
		},
	}
	if cfg.ReadOnly {
		gcfg := rcfg["grpc"].(map[string]interface{})
		gcfg["interceptors"] = map[string]interface{}{
			"readonly": map[string]interface{}{},
		}
	}
	return rcfg
}
