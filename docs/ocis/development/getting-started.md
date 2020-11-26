---
title: "Getting Started"
date: 2020-07-07T20:35:00+01:00
weight: 15
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: getting-started.md
---

{{< toc >}}

## Requirements

We want contribution to oCIS and the creation of extensions to be as easy as possible.
So we are trying to reflect this the used tooling. It should be kept simple and quick to be set up.

Besides of standard development tools like git and a text editor, you need following software for development:

- Go >= v1.13 ([install instructions](https://golang.org/doc/install))
- Yarn ([install instructions](https://classic.yarnpkg.com/en/docs/install))
- docker ([install instructions](https://docs.docker.com/get-docker/))
- docker-compose ([install instructions](https://docs.docker.com/compose/install/))

If you find tools needed besides the mentioned above, please feel free to open an issue or open a PR.

## Repository structure

This repository follows the [golang-standard project-layout](https://github.com/golang-standards/project-layout).

oCIS consists of multiple micro services, also called extensions. We started by having standalone repositories for each of them but quickly noticed, that this adds a time consuming overhead for developers. So we ended up with a monorepo housing all the extensions in one repository.

Each of the extensions live in a subfolder (eg. `accounts` or `settings`) in this repository, technically creating independant Go modules.

The `ocis` folder does also contain a Go module but is no extension at all. Instead this module is used to import all extensions and furthermore implement commands to start the extensions. With the resulting oCIS binary you can start single extensions or even all extensions at the same time.

The `docs` folder contains the source for the [oCIS documentation](https://owncloud.github.io/ocis/).

The `deployments` folder contains documented deployment configurations and templates.

The `scripts` folder contains scripts to perform various build, install, analysis, etc operations.

## Starting points

Depending on what you want do develop there are different starting points. These will be described below.

### Developing oCIS

If you want to contribute to oCIS:

- see [contribution guidelines](https://github.com/owncloud/ocis#contributing)
- make sure the tooling is set up by [building oCIS]({{< relref "build.md" >}}) and [building the docs]({{< relref "build-docs.md" >}})
- create or pick an [open issue](https://github.com/owncloud/ocis/issues) to develop on
- open a PR and get things done

### Developing extensions

If you want to develop an extension, start here: [Extensions]({{< relref "extensions.md">}})
