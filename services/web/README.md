# Web

The web service embeds and serves the static files for the [Infinite Scale Web Client](https://github.com/owncloud/web).
Note that clients will respond with a connection error if the web service is not available.

The web service also provides a minimal API for branding functionality like changing the logo shown.

## Custom Compiled Web Assets

If you want to use your custom compiled web client assets instead of the embedded ones, then you can do that by setting the `WEB_ASSET_PATH` variable to point to your compiled files. See [ownCloud Web / Getting Started](https://owncloud.dev/clients/web/getting-started/) and [ownCloud Web / Setup with oCIS](https://owncloud.dev/clients/web/backend-ocis/) for more details.

## Web UI Configuration

Note that single configuration settings of the embedded web UI can be defined via `WEB_OPTION_xxx` environment variables. If a json based configuration file is used via the `WEB_UI_CONFIG_FILE` environment variable, these configurations take precedence over single options set.

### Web UI Options

Besides theming, the behavior of the web UI can be configured via options. See the environment variables `WEB_OPTION_xxx` for more details.

### Web UI Config File

When defined via the `WEB_UI_CONFIG_FILE` environment variable, the configuration of the web UI can be made with a [json based](https://github.com/owncloud/web/tree/master/config) file.

### Embedding Web

Web can be consumed by another application in a stripped down version called “Embed mode”. This mode is supposed to be used in the context of selecting or sharing resources. For more details see the developer documentation [ownCloud Web / Embed Mode](https://owncloud.dev/clients/web/embed-mode/). See the environment variables: `WEB_OPTION_MODE` and `WEB_OPTION_EMBED_TARGET` to configure the embedded mode.

# Web Apps

The administrator of the environment is capable of providing custom web applications to the users.
This feature is useful for organizations that want to provide third party or custom apps to their users.

It's important to note that the feature at the moment is only capable of providing static (js, mjs, e.g.) web applications
and does not support injection of dynamic web applications (custom dynamic backends).

## Loading Applications

Loading applications is achieved in three ways:

* By loading a web application from a user-provided path, by setting the `WEB_ASSET_APPS_PATH` environment variable.
* By loading a web application from the default ocis home directory, e.g. `$OCIS_BASE_DATA_PATH/web/assets/apps`.
* By embedding a web application into the ocis binary which is a build-time option,
  e.g. `ocis_src_path/services/web/assets/apps` followed by a build.

The list of available applications is composed of the build in extensions and the custom applications
provided by the administrator, e.g. `WEB_ASSET_APPS_PATH` or `$OCIS_BASE_DATA_PATH/web/assets/apps`.

For example, if ocis would contain a build in extension named `image-viewer-dfx` and the administrator provides a custom application named `image-viewer-obj` in the `WEB_ASSET_APPS_PATH` directory, the user will be able to access both applications from the web ui.

## Application Structure

Applications always have to follow a strict structure, which is as follows:

* each application must be in its own directory
* each application directory must contain a `manifest.json` file

Everything else is skipped and not considered as an application.

The `manifest.json` file contains the following fields:

* `entrypoint` - required - the entrypoint of the application, e.g. `index.js`, the path is relative to the parent directory
* `config` - optional - a list of key-value pairs that are passed to the global web application configuration

## Application Configuration

It's important to note that an application manifest should never be changed manually;
if a custom configuration is needed, the administrator should provide the required configuration inside the
`$OCIS_BASE_DATA_PATH/config/apps.yaml` file.

The `apps.yaml` file must contain a list of key-value pairs which gets merged with the `config` field.
For example, if the `image-viewer-obj` application contains the following configuration:

```json
{
  "entrypoint": "index.js",
  "config": {
    "maxWith": 1280,
    "maxHeight": 1280
  }
}
```

and the `apps.yaml` file contains the following configuration:

```yaml
image-viewer-obj:
  config:
    maxHeight: 640
    maxSize: 512
```

the final configuration for web will be:

```json
{
  "external_apps": [
    {
      "id": "image-viewer-obj",
      "path": "index.js",
      "config": {
        "maxWith": 1280,
        "maxHeight": 640,
        "maxSize": 512
      }
    }
  ]
}
```

besides the configuration from the `manifest.json` file, the `apps.yaml` file can also contain the following fields:

* `disabled` - optional - defaults to `false` - if set to `true`, the application will not be loaded

The local provided configuration yaml will always override the shipped application manifest configuration.

## Fallback Mechanism

Besides the configuration and application registration, there is one further important aspect to know;
in the process of loading the application assets, the system uses a fallback mechanism to load the assets.

This is incredibly useful for cases where just a single asset should be overwritten, e.g., a logo or similar.

Consider the following, ocis is shipped with a default extension named `image-viewer-dfx` which contains a logo,
but the administrator wants to provide a custom logo for the `image-viewer-dfx` application.

This can be achieved by providing a custom logo in the `WEB_ASSET_APPS_PATH` directory,
e.g. `WEB_ASSET_APPS_PATH/image-viewer-dfx/logo.png`.
Every other asset is loaded from the build in extension, but the logo is loaded from the custom directory.

The same applies for the `manifest.json` file, if the administrator wants to provide a custom `manifest.json` file.

## Miscellaneous

Please note that ocis needs a restart to load new applications or changes to the `apps.yaml` file.
