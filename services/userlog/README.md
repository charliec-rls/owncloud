# Userlog Service

The `userlog` service is a mediator between the `eventhistory` service and clients who want to be informed about user related events. It provides an API to retrieve those.

## Prerequisites

Running the `userlog` service without running the `eventhistory` service is not possible.

## Storing

The `userlog` service persists information via the configured store in `USERLOG_STORE_TYPE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured redis cluster.
  -   `redis-sentinel`: Stores data in a configured redis sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.

1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The userlog service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `USERLOG_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.

## Configuring

For the time being, the configuration which user related events are of interest is hardcoded and cannot be changed.

## Retrieving

The `userlog` service provides an API to retrieve configured events. For now, this API is mostly following the [oc10 notification GET API](https://doc.owncloud.com/server/next/developer_manual/core/apis/ocs-notification-endpoint-v1.html#get-user-notifications).

## Deleting

To delete events for an user, use a `DELETE` request to `ocs/v2.php/apps/notifications/api/v1/notifications` containing the IDs to delete.

## Translations

The `userlog` service has embedded translations sourced via transifex to provide a basic set of translated languages. These embedded translations are available for all deployment scenarios. In addition, the service supports custom translations, though it is currently not possible to just add custom translations to embedded ones. If custom translations are configured, the embedded ones are not used. To configure custom translations, the `USERLOG_TRANSLATION_PATH` environment variable needs to point to a base folder that will contain the translation files. This path must be available from all instances of the userlog service, a shared storage is recommended. Translation files must be of type  [.po](https://www.gnu.org/software/gettext/manual/html_node/PO-Files.html#PO-Files) or [.mo](https://www.gnu.org/software/gettext/manual/html_node/Binaries.html). For each language, the filename needs to be `userlog.po` (or `userlog.mo`) and stored in a folder structure defining the language code. In general the path/name pattern for a translation file needs to be:

```text
{USERLOG_TRANSLATION_PATH}/{language-code}/LC_MESSAGES/userlog.po
```

The language code pattern is composed of `language[_territory]` where  `language` is the base language and `_territory` is optional and defines a country.

For example, for the language `de_DE`, one needs to place the corresponding translation files to `{USERLOG_TRANSLATION_PATH}/de_DE/LC_MESSAGES/userlog.po`.

### Translation Rules

*   If a requested language code is not available, the service tries to fall back to the base language if available. For example, if `de_DE` is not available, the service tries to fall back to translations in the `de` folder.
*   If the base language `de` is also not available, the service falls back to the system's default English (`en`), which is the source of the texts provided by the code.
