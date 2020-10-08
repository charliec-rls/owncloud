---
title: "Configuration"
date: "2020-09-21T13:14:51+0200"
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/proxy
geekdocFilePath: configuration.md
---

{{< toc >}}

## Configuration

oCIS Single Binary is not responsible for configuring extensions. Instead, each extension could either be configured by environment variables, cli flags or config files.

Each extension has its dedicated documentation page (e.g. https://owncloud.github.io/extensions/ocis_proxy/configuration) which lists all possible configurations. Config files and environment variables are picked up if you use the `./bin/ocis server` command within the oCIS single binary. Command line flags must be set explicitly on the extensions subcommands.

### Configuration using config files

Out of the box extensions will attempt to read configuration details from:

```console
/etc/ocis
$HOME/.ocis
./config
```

For this configuration to be picked up, have a look at your extension `root` command and look for which default config name it has assigned. *i.e: ocis-proxy reads `proxy.json | yaml | toml ...`*.

So far we support the file formats `JSON` and `YAML`, if you want to get a full example configuration just take a look at [our repository](https://github.com/owncloud/ocis/tree/master/config), there you can always see the latest configuration format. These example configurations include all available options and the default values. The configuration file will be automatically loaded if it's placed at `/etc/ocis/ocis.yml`, `${HOME}/.ocis/ocis.yml` or `$(pwd)/config/ocis.yml`.

### Envrionment variables

If you prefer to configure the service with environment variables you can see the available variables below.

### Commandline flags

If you prefer to configure the service with commandline flags you can see the available variables below. Command line flags are only working when calling the subcommand directly.

## Root Command

proxy for Reva/oCIS

Usage: `ocis-proxy [global options] command [command options] [arguments...]`

--config-file | $PROXY_CONFIG_FILE  
: Path to config file.

--log-level | $PROXY_LOG_LEVEL  
: Set logging level. Default: `info`.

--log-pretty | $PROXY_LOG_PRETTY  
: Enable pretty logging. Default: `true`.

--log-color | $PROXY_LOG_COLOR  
: Enable colored logging. Default: `true`.

## Sub Commands

### ocis-proxy health

Check health status

Usage: `ocis-proxy health [command options] [arguments...]`

--debug-addr | $PROXY_DEBUG_ADDR  
: Address to debug endpoint. Default: `0.0.0.0:9109`.

### ocis-proxy server

Start integrated server

Usage: `ocis-proxy server [command options] [arguments...]`

--tracing-enabled | $PROXY_TRACING_ENABLED  
: Enable sending traces.

--tracing-type | $PROXY_TRACING_TYPE  
: Tracing backend type. Default: `jaeger`.

--tracing-endpoint | $PROXY_TRACING_ENDPOINT  
: Endpoint for the agent.

--tracing-collector | $PROXY_TRACING_COLLECTOR  
: Endpoint for the collector. Default: `http://localhost:14268/api/traces`.

--tracing-service | $PROXY_TRACING_SERVICE  
: Service name for tracing. Default: `proxy`.

--debug-addr | $PROXY_DEBUG_ADDR  
: Address to bind debug server. Default: `0.0.0.0:9205`.

--debug-token | $PROXY_DEBUG_TOKEN  
: Token to grant metrics access.

--debug-pprof | $PROXY_DEBUG_PPROF  
: Enable pprof debugging.

--debug-zpages | $PROXY_DEBUG_ZPAGES  
: Enable zpages debugging.

--http-addr | $PROXY_HTTP_ADDR  
: Address to bind http server. Default: `0.0.0.0:9200`.

--http-root | $PROXY_HTTP_ROOT  
: Root path of http server. Default: `/`.

--asset-path | $PROXY_ASSET_PATH  
: Path to custom assets.

--http-namespace | $PROXY_HTTP_NAMESPACE  
: Set the base namespace for the http namespace. Default: `com.owncloud`.

--transport-tls-cert | $PROXY_TRANSPORT_TLS_CERT  
: Certificate file for transport encryption.

--transport-tls-key | $PROXY_TRANSPORT_TLS_KEY  
: Secret file for transport encryption.

--tls | $PROXY_TLS  
: Use TLS (disable only if proxy is behind a TLS-terminating reverse-proxy).. Default: `true`.

--jwt-secret | $PROXY_JWT_SECRET  
: Used to create JWT to talk to reva, should equal reva's jwt-secret. Default: `Pive-Fumkiu4`.

--reva-gateway-addr | $PROXY_REVA_GATEWAY_ADDR  
: REVA Gateway Endpoint. Default: `127.0.0.1:9142`.

--oidc-issuer | $PROXY_OIDC_ISSUER  
: OIDC issuer. Default: `https://localhost:9200`.

--oidc-insecure | $PROXY_OIDC_INSECURE  
: OIDC allow insecure communication. Default: `true`.

