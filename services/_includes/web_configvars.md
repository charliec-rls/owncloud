## Environment Variables

| Name | Type | Default Value | Description |
|------|------|---------------|-------------|
| OCIS_TRACING_ENABLED<br/>WEB_TRACING_ENABLED | bool | false | Activates tracing.|
| OCIS_TRACING_TYPE<br/>WEB_TRACING_TYPE | string |  | The type of tracing. Defaults to '', which is the same as 'jaeger'. Allowed tracing types are 'jaeger' and '' as of now.|
| OCIS_TRACING_ENDPOINT<br/>WEB_TRACING_ENDPOINT | string |  | The endpoint of the tracing agent.|
| OCIS_TRACING_COLLECTOR<br/>WEB_TRACING_COLLECTOR | string |  | The HTTP endpoint for sending spans directly to a collector, i.e. http://jaeger-collector:14268/api/traces. Only used if the tracing endpoint is unset.|
| OCIS_LOG_LEVEL<br/>WEB_LOG_LEVEL | string |  | The log level. Valid values are: 'panic', 'fatal', 'error', 'warn', 'info', 'debug', 'trace'.|
| OCIS_LOG_PRETTY<br/>WEB_LOG_PRETTY | bool | false | Activates pretty log output.|
| OCIS_LOG_COLOR<br/>WEB_LOG_COLOR | bool | false | Activates colorized log output.|
| OCIS_LOG_FILE<br/>WEB_LOG_FILE | string |  | The path to the log file. Activates logging to this file if set.|
| WEB_DEBUG_ADDR | string | 127.0.0.1:9104 | Bind address of the debug server, where metrics, health, config and debug endpoints will be exposed.|
| WEB_DEBUG_TOKEN | string |  | Token to secure the metrics endpoint.|
| WEB_DEBUG_PPROF | bool | false | Enables pprof, which can be used for profiling.|
| WEB_DEBUG_ZPAGES | bool | false | Enables zpages, which can be used for collecting and viewing in-memory traces.|
| WEB_HTTP_ADDR | string | 127.0.0.1:9100 | The bind address of the HTTP service.|
| OCIS_HTTP_TLS_ENABLED | bool | false | Activates TLS for the http based services using the server certifcate and key configured via OCIS_HTTP_TLS_CERTIFICATE and OCIS_HTTP_TLS_KEY. If OCIS_HTTP_TLS_CERTIFICATE is not set a temporary server certificate is generated - to be used with PROXY_INSECURE_BACKEND=true.|
| OCIS_HTTP_TLS_CERTIFICATE | string |  | Path/File name of the TLS server certificate (in PEM format) for the http services.|
| OCIS_HTTP_TLS_KEY | string |  | Path/File name for the TLS certificate key (in PEM format) for the server certificate to use for the http services.|
| WEB_HTTP_ROOT | string | / | Subdirectory that serves as the root for this HTTP service.|
| WEB_CACHE_TTL | int | 604800 | Cache policy in seconds for ownCloud Web assets.|
| OCIS_CORS_ALLOW_ORIGINS<br/>WEB_CORS_ALLOW_ORIGINS | []string | [https://localhost:9200] | A list of allowed CORS origins. See following chapter for more details: *Access-Control-Allow-Origin* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_METHODS<br/>WEB_CORS_ALLOW_METHODS | []string | [OPTIONS HEAD GET PUT PATCH POST DELETE MKCOL PROPFIND PROPPATCH MOVE COPY REPORT SEARCH] | A list of allowed CORS methods. See following chapter for more details: *Access-Control-Request-Method* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Method. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_HEADERS<br/>WEB_CORS_ALLOW_HEADERS | []string | [Origin Accept Content-Type Depth Authorization Ocs-Apirequest If-None-Match If-Match Destination Overwrite X-Request-Id X-Requested-With Tus-Resumable Tus-Checksum-Algorithm Upload-Concat Upload-Length Upload-Metadata Upload-Defer-Length Upload-Expires Upload-Checksum Upload-Offset X-HTTP-Method-Override] | A list of allowed CORS headers. See following chapter for more details: *Access-Control-Request-Headers* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Request-Headers. See the Environment Variable Types description for more details.|
| OCIS_CORS_ALLOW_CREDENTIALS<br/>WEB_CORS_ALLOW_CREDENTIALS | bool | false | Allow credentials for CORS. See following chapter for more details: *Access-Control-Allow-Credentials* at https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials.|
| WEB_ASSET_PATH | string |  | Serve ownCloud Web assets from a path on the filesystem instead of the builtin assets.|
| WEB_ASSET_CORE_PATH | string | /var/lib/ocis/web/assets/core | Serve ownCloud Web assets from a path on the filesystem instead of the builtin assets.|
| WEB_ASSET_APPS_PATH | string | /var/lib/ocis/web/assets/apps | Serve ownCloud Web apps assets from a path on the filesystem instead of the builtin assets.|
| WEB_UI_CONFIG_FILE | string |  | Read the ownCloud Web json based configuration from this path/file. The config file takes precedence over WEB_OPTION_xxx environment variables. See the text description for more details.|
| OCIS_URL<br/>WEB_UI_THEME_SERVER | string | https://localhost:9200 | Base URL to load themes from. Will be prepended to the theme path.|
| WEB_UI_THEME_PATH | string | /themes/owncloud/theme.json | Subpath/file to load the theme. Will be appended to the URL of the theme server.|
| OCIS_URL<br/>WEB_UI_CONFIG_SERVER | string | https://localhost:9200 | URL, where the oCIS APIs are reachable for ownCloud Web.|
| WEB_OIDC_METADATA_URL | string | https://localhost:9200/.well-known/openid-configuration | URL for the OIDC well-known configuration endpoint. Defaults to the oCIS API URL + '/.well-known/openid-configuration'.|
| OCIS_URL<br/>OCIS_OIDC_ISSUER<br/>WEB_OIDC_AUTHORITY | string | https://localhost:9200 | URL of the OIDC issuer. It defaults to URL of the builtin IDP.|
| OCIS_OIDC_CLIENT_ID<br/>WEB_OIDC_CLIENT_ID | string | web | The OIDC client ID which ownCloud Web uses. This client needs to be set up in your IDP. Note that this setting has no effect when using the builtin IDP.|
| WEB_OIDC_RESPONSE_TYPE | string | code | The OIDC response type to use for authentication.|
| WEB_OIDC_SCOPE | string | openid profile email | OIDC scopes to request during authentication to authorize access to user details. Defaults to 'openid profile email'. Values are separated by blank. More example values but not limited to are 'address' or 'phone' etc.|
| WEB_OIDC_POST_LOGOUT_REDIRECT_URI | string |  | This value needs to point to a valid and reachable web page. The web client will trigger a redirect to that page directly after the logout action. The default value is empty and redirects to the login page.|
| WEB_OPTION_HOME_FOLDER | string |  | Specifies a folder that is used when the user navigates 'home'. Navigating home gets triggered by clicking on the 'All files' menu item. The user will not be jailed in that directory, it simply serves as a default location. A static location can be provided, or variables of the user object to come up with a user specific home path can be used. This uses the twig template variable style and allows picking a value or a substring of a value of the authenticated user. Examples are '/Shares', '/{{.Id}}' and '/{{substr 0 3 .Id}}/{{.Id}'.|
| WEB_OPTION_OPEN_APPS_IN_TAB | bool | false | Configures whether apps and extensions should generally open in a new tab. Defaults to false.|
| OCIS_DISABLE_PREVIEWS<br/>WEB_OPTION_DISABLE_PREVIEWS | bool | false | Set this option to 'true' to disable previews in all the different web file listing views. This can speed up file listings in folders with many files. The only list view that is not affected by this setting is the trash bin, as it does not allow previewing at all.|
| WEB_OPTION_PREVIEW_FILE_MIMETYPES | []string | [image/gif image/png image/jpeg text/plain image/tiff image/bmp image/x-ms-bmp application/vnd.geogebra.slides] | A list of mimeTypes to specify which ones will be previewed in the UI. For example, to only preview jpg and text files, set this option to 'image/jpeg,text/plain'. See the Environment Variable Types description for more details.|
| WEB_OPTION_DISABLE_FEEDBACK_LINK | bool | false | Set this option to 'true' to disable the feedback link in the topbar. Keeping it enabled by setting the value to 'false' or with the absence of the option, allows ownCloud to get feedback from your user base through a dedicated survey website.|
| WEB_OPTION_SHARING_RECIPIENTS_PER_PAGE | int | 200 | Sets the number of users shown as recipients in the dropdown menu when sharing resources.|
| WEB_OPTION_SIDEBAR_SHARES_SHOW_ALL_ON_LOAD | bool | false | Sets the list of the (link) shares list in the sidebar to be initially expanded. Default is a collapsed state, only showing the first three shares.|
| WEB_OPTION_RUNNING_ON_EOS | bool | false | Set this option to 'true' if running on an EOS storage backend (https://eos-web.web.cern.ch/eos-web/) to enable its specific features. Defaults to 'false'.|
| WEB_OPTION_HOVERABLE_QUICK_ACTIONS | bool | false | Set this option to 'true' to hide quick actions (buttons appearing on file rows) and only show them when the user hovers over the row with his mouse. Defaults to 'false'.|
| WEB_OPTION_ROUTING_ID_BASED | bool | true | Enable or disable fileIds being added to the URL. Defaults to 'true', because otherwise spaces with name clashes cannot be resolved correctly. Note: Only disable this if you can guarantee on the server side, that spaces of the same namespace cannot have name clashes.|
| WEB_OPTION_CONTEXTHELPERS_READ_MORE | bool | true | Specifies whether the 'Read more' link should be displayed or not.|
| WEB_OPTION_LOGOUT_URL | string |  | Adds a link to the user's profile page to point him to an external page, where he can manage his session and devices. This is helpful when an external IdP is used. This option is disabled by default.|
| WEB_OPTION_LOGIN_URL | string |  | Specifies the target URL to the login page. This is helpful when an external IdP is used. This option is disabled by default. Example URL like: https://www.myidp.com/login.|
| WEB_OPTION_OPEN_LINKS_WITH_DEFAULT_APP | bool | true | Specifies whether single file link shares should be opened with the default app or not. If not opened by the default app, the Web UI just displays the file details.|
| WEB_OPTION_TOKEN_STORAGE_LOCAL | bool | true | Specifies whether the access token will be stored in the local storage when set to 'true' or in the session storage when set to 'false'. If stored in the local storage, login state will be persisted across multiple browser tabs, means no additional logins are required.|
| WEB_OPTION_DISABLED_EXTENSIONS | []string | [] | A list to disable specific Web extensions identified by their ID. The ID can e.g. be taken from the 'index.ts' file of the web extension. Example: 'com.github.owncloud.web.files.search,com.github.owncloud.web.files.print'. See the Environment Variable Types description for more details.|
| WEB_OPTION_USER_LIST_REQUIRES_FILTER | bool | false | Defines whether one or more filters must be set in order to list users in the Web admin settings. Set this option to 'true' if running in an environment with a lot of users and listing all users could slow down performance. Defaults to 'false'.|
| WEB_OPTION_CONCURRENT_REQUESTS_RESOURCE_BATCH_ACTIONS | int | 0 | Defines the maximum number of concurrent requests per file/folder/space batch action. Defaults to 4.|
| WEB_OPTION_CONCURRENT_REQUESTS_SSE | int | 0 | Defines the maximum number of concurrent requests in SSE event handlers. Defaults to 4.|
| WEB_OPTION_CONCURRENT_REQUESTS_SHARES_CREATE | int | 0 | Defines the maximum number of concurrent requests per sharing invite batch. Defaults to 4.|
| WEB_OPTION_CONCURRENT_REQUESTS_SHARES_LIST | int | 0 | Defines the maximum number of concurrent requests when loading individual share information inside listings. Defaults to 2.|
| OCIS_JWT_SECRET<br/>WEB_JWT_SECRET | string |  | The secret to mint and validate jwt tokens.|
| WEB_GATEWAY_GRPC_ADDR | string | com.owncloud.api.gateway | The bind address of the GRPC service.|