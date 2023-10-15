# Thumbnails

The thumbnails service provides methods to generate thumbnails for various files and resolutions based on requests. It retrieves the sources at the location where the user files are stored and saves the thumbnails where system files are stored. Those locations have defaults but can be manually defined via environment variables.

## File Locations Overview

The relevant environment variables defining file locations are:

-   (1) `OCIS_BASE_DATA_PATH`
-   (2) `STORAGE_USERS_OCIS_ROOT`
-   (3) `THUMBNAILS_FILESYSTEMSTORAGE_ROOT`

(1) ... Having a default set by the Infinite Scale code, but if defined, used as base path for other services.
(2) ... Source files, defaults to (1) plus path component, but can be freely defined if required.
(3) ... Target files, defaults to (1) plus path component, but can be freely defined if required.

For details and defaults for these environment variables see the ocis admin documentation.

## Thumbnail Location

It may be beneficial to define the location of the thumbnails to be other than the default (with system files). This is due the fact that storing thumbnails can consume a lot of space over time which not necessarily needs to reside on the same partition or mount or expensive drives.

## Thumbnail Source File Types

Thumbnails can be generated from the following source file types:

-   png
-   jpg
-   gif
-   tiff
-   bmp
-   txt

The thumbnail service retrieves source files using the information provided by the backend. The Linux backend identifies source files usually based on the extension.

If a file type was not properly assigned or the type identification failed, thumbnail generation will fail and an error will be logged.

## Thumbnail Target File Types

Thumbnails can either be generated as `png`, `jpg` or `gif` files. These types are hardcoded and no other types can be requested. A requestor, like another service or a client, can request one of the available types to be generated. If more than one type is required, each type must be requested individually.

## Thumbnail Resolution

Various resolutions can be defined via `THUMBNAILS_RESOLUTIONS`. A requestor can request any arbitrary resolution and the thumbnail service will use the one closest to the requested resolution. If more than one resolution is required, each resolution must be requested individually.

Example:

Requested: 18x12
Available: 30x20, 15x10, 9x6
Returned: 15x10

## Thumbnail Processors

Image generation can be configured by defining different processors, following processors are available:

* resize
* fit
* fill
* thumbnail

to apply one of those, a query parameter has to be added to the request, e.g. `?processor=fit`

## Deleting Thumbnails

As of now, there is no automated thumbnail deletion. This is especially true when a source file gets deleted or moved. This situation will be solved at a later stage. For the time being, if you run short on physical thumbnails space, you have to manually delete the thumbnail store to free space. Thumbnails will then be recreated on request.

## Memory Considerations

Since source files need to be loaded into memory when generating thumbnails, large source files could potentially crash this service if there is insufficient memory available. For bigger instances when using container orchestration deployment methods, this service can be dedicated to its own server(s) with more memory.
