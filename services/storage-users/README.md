# Storage-Users Service

Purpose and description to be added

## CLI Commands

### Manage Unfinished Uploads

<!-- referencing: [oCIS FS] clean up aborted uploads https://github.com/owncloud/ocis/issues/2622 -->

When using Infinite Scale as user storage, a directory named `storage/users/uploads` can be found in the Infinite Scale data folder. This is an intermediate directory based on [TUS](https://tus.io) which is an open protocol for resumable uploads. Each upload consists of a _blob_ and a _blob.info_ file. Note that the term _blob_ is just a placeholder.

*   **If an upload succeeds**, the _blob_ file will be moved to the target and the _blob.info_ file will be deleted.

*   **In case of incomplete uploads**, the _blob_ and _blob.info_ files will continue to receive data until either the upload succeeds in time or the upload expires based on the `STORAGE_USERS_UPLOAD_EXPIRATION` variable, see the table below for details.

*   **In case of expired uploads**, the _blob_ and _blob.info_ files will _not_ be removed automatically. Thus a lot of data can pile up over time wasting storage space.

*   **In the rare case of a failure**, after the upload succeeded but the file was not moved to its target location, which can happen when postprocessing fails, the situation is the same as with expired uploads.

Example cases for expired uploads

*   When a user uploads a big file but the file exceeds the user-quota, the upload can't be moved to the target after it has finished. The file stays at the upload location until it is manually cleared.
*   If the bandwidth is limited and the file to transfer can't be transferred completely before the upload expiration time is reached, the file expires and can't be processed. 

There are two commands available to manage unfinished uploads

```bash
ocis storage-users uploads <command>
```

```plaintext
COMMANDS:
   list     Print a list of all incomplete uploads
   clean    Clean up leftovers from expired uploads
```

#### Command Examples

Command to identify incomplete uploads

```bash
ocis storage-users uploads list
```

```plaintext
Incomplete uploads:
 - 455bd640-cd08-46e8-a5a0-9304908bd40a (file_example_PPT_1MB.ppt, Size: 1028608, Expires: 2022-08-17T12:35:34+02:00)
```

Command to clear expired uploads
```bash
ocis storage-users uploads clean
```

```plaintext
Cleaned uploads:
- 455bd640-cd08-46e8-a5a0-9304908bd40a (Filename: file_example_PPT_1MB.ppt, Size: 1028608, Expires: 2022-08-17T12:35:34+02:00)
```

### Purge Expired Space Trash-Bins Items

<!-- referencing: https://github.com/owncloud/ocis/pull/5500 -->

This command is about purging old trash-bin items of `project` spaces (spaces that have been created manually) and `personal` spaces.

```bash
ocis storage-users trash-bin <command>
```

```plaintext
COMMANDS:
   purge-expired     Purge all expired items from the trashbin
```

The configuration for the `purge-expired` command is done by using the following environment variables.

*   `STORAGE_USERS_PURGE_TRASH_BIN_USER_ID` is used to obtain space trash-bin information and takes the system admin user as the default which is the `OCIS_ADMIN_USER_ID` but can be set individually. It should be noted, that the `OCIS_ADMIN_USER_ID` is only assigned automatically when using the single binary deployment and must be manually assigned in all other deployments. The command only considers spaces to which the assigned user has access and delete permission.

*   `STORAGE_USERS_PURGE_TRASH_BIN_PERSONAL_DELETE_BEFORE` has a default value of `30 days`, which means the command will delete all files older than `30 days`. The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is equivalent to disable and prevents the deletion of `personal space` trash-bin files.

*   `STORAGE_USERS_PURGE_TRASH_BIN_PROJECT_DELETE_BEFORE` has a default value of `30 days`, which means the command will delete all files older than `30 days`. The value is human-readable, valid values are `24h`, `60m`, `60s` etc. `0` is equivalent to disable and prevents the deletion of `project space` trash-bin files.

## Caching

The `storage-users` service caches file metadata via the configured store in `STORAGE_USERS_CACHE_STORE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `redis`: Stores data in a configured Redis cluster.
  -   `redis-sentinel`: Stores data in a configured Redis Sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.

1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The `storage-users` service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using the `redis-sentinel` cache, the Redis master to use is configured using as part of the `STORAGE_USERS_CACHE_NODES`  in the form of `<sentinel-host>:<sentinel-port>/<redis-master>`, e.g. `10.10.0.200:26379/mymaster`.