---
title: "Testing"
date: 2018-05-02T00:00:00+00:00
weight: 37
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: testing.md
---

{{< toc >}}

To run tests in the test suite you have two options. You may go the easy way and just run the test suite in docker. But for some tasks you could also need to install the test suite natively, which requires a little more setup since PHP and some dependencies need to be installed.

Both ways to run tests with the test suites are described here.

## Testing with test suite in docker

Let's see what is available. Invoke the following command from within the root of the oCIS repository.

```
make -C tests/acceptance/docker help
```

Basically we have two sources for feature tests and test suites:

- [oCIS feature test and test suites](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features)
- [tests and test suites transferred from ownCloud, they have prefix coreApi](https://github.com/owncloud/ocis/tree/master/tests/acceptance/features)

At the moment both can be applied to oCIS since the api of oCIS is designed to be compatible with ownCloud.

As a storage backend, we offer oCIS native storage, also called "ocis". This stores files directly on disk. Along with that we also provide `S3` storage driver.

You can invoke two types of test suite runs:

- run a full test suite, which consists of multiple feature tests
- run a single feature or single scenario in a feature

### Run full test suite

The names of the full test suite make targets have the same naming as in the CI pipeline. The available local ocis specific test suites are `apiAccountsHashDifficulty`, `apiArchiver`, `apiContract`, `apiGraph`, `apiSpaces`, `apiSpacesShares`, `apiCors`, `apiAsyncUpload`. They can be run with `ocis` storage and `S3` storage.

> Note: In order to see the tests log attach `show-test-logs` in the command

For example `make -C tests/acceptance/docker localApiTests-apiAccountsHashDifficulty-ocis` runs the same tests as the `localApiTests-apiAccountsHashDifficulty-ocis` CI pipeline, which runs the oCIS test suite "apiAccountsHashDifficulty" against an oCIS with ocis storage.

For example `make -C tests/acceptance/docker localApiTests-apiAccountsHashDifficulty-s3ng` runs the oCIS test suite "apiAccountsHashDifficulty" against an oCIS with s3 storage.

> Note: To run the tests from `apiAsyncUpload` suite you need to provide extra environment variable `POSTPROCESSING_DELAY`

For example `make -C tests/acceptance/docker Core-API-Tests-ocis-storage-3` runs the same tests as the `Core-API-Tests-ocis-storage-3` CI pipeline, which runs the third (out of ten) test suite transferred from ownCloud against an oCIS with ocis storage.

For example `make -C tests/acceptance/docker Core-API-Tests-s3ng-storage-3` runs the third (out of ten) test suite transferred from ownCloud against an oCIS with s3 storage.

### Run single feature test

The single feature tests can also be run against the different storage backends. Therefore, multiple make targets with the schema test-<test source>-feature-<storage backend> exist. To select a single feature you have to add an additional `BEHAT_FEATURE=...` parameter when invoking the make command:

```
make -C tests/acceptance/docker test-ocis-feature-ocis-storage BEHAT_FEATURE='tests/acceptance/features/apiAccountsHashDifficulty/addUser.feature:21'
```

This must be pointing to a valid feature definition.

To run a single scenario in a feature, then mention the line number of the scenario:

```
make -C tests/acceptance/docker test-ocis-feature-ocis-storage BEHAT_FEATURE='tests/acceptance/features/apiAccountsHashDifficulty/addUser.feature:21'
```

Similarly, with S3 storage,
- run a whole feature
```
make -C tests/acceptance/docker test-ocis-feature-s3ng-storage BEHAT_FEATURE='tests/acceptance/features/apiAccountsHashDifficulty/addUser.feature'
```

- run a single scenario
```
make -C tests/acceptance/docker test-ocis-feature-s3ng-storage BEHAT_FEATURE='tests/acceptance/features/apiAccountsHashDifficulty/addUser.feature:21'
```

In the same way for the tests transferred from oc10 can be run as
- run a whole feature
```
make -C tests/acceptance/docker test-core-feature-ocis-storage BEHAT_FEATURE='tests/acceptance/features/coreApiAuth/webDavAuth.feature'
```

- run a single scenario
```
make -C tests/acceptance/docker test-core-feature-ocis-storage BEHAT_FEATURE='tests/acceptance/features/coreApiAuth/webDavAuth.feature:13'
```

> Note: the tests transferred from oc10 start with coreApi

### oCIS image to be tested (or: skip build and take existing image)

By default, the tests will be run against the docker image built from your current working state of the oCIS repository. For some purposes it might also be handy to use an oCIS image from Docker Hub. Therefore, you can provide the optional flag `OCIS_IMAGE_TAG=...` which must contain an available docker tag of the [owncloud/ocis registry on Docker Hub](https://hub.docker.com/r/owncloud/ocis) (e.g. 'latest').

```
make -C tests/acceptance/docker localApiTests-apiAccountsHashDifficulty-ocis OCIS_IMAGE_TAG=latest
```

### Test log output

While a test is running or when it is finished, you can attach to the logs generated by the tests.

```
make -C tests/acceptance/docker show-test-logs
```

{{< hint info >}}
The log output is opened in `less`. You can navigate up and down with your cursors. By pressing "F" you can follow the latest line of the output.
{{< /hint >}}

### Cleanup

During testing we start a redis and oCIS docker container. These will not be stopped automatically. You can stop them with:

```
make -C tests/acceptance/docker clean
```

## Testing with test suite natively installed

We have two sets of tests:
- `test-acceptance-from-core-api` set was transferred from [core](https://github.com/owncloud/core) repository
The suite name of all tests transferred from the core starts with "core"

- `test-acceptance-api` set was created for ocis. Mainly for testing spaces features


### Run ocis

Create an up-to-date ocis binary by [building oCIS]({{< ref "build" >}})

To start ocis:

```bash
IDM_ADMIN_PASSWORD=admin \
ocis/bin/ocis init --insecure true

OCIS_INSECURE=true PROXY_ENABLE_BASIC_AUTH=true \
ocis/bin/ocis server
```

`PROXY_ENABLE_BASIC_AUTH` will allow the acceptance tests to make requests against the provisioning api (and other endpoints) using basic auth.

### Run the test-acceptance-from-core-api tests

```bash
make test-acceptance-from-core-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
```
Note: This command only works for suites that start with core

### Run the test-acceptance-api tests

```bash
make test-acceptance-api \
TEST_SERVER_URL=https://localhost:9200 \
TEST_WITH_GRAPH_API=true \
TEST_OCIS=true \
```

Make sure to adjust the settings `TEST_SERVER_URL` according to your environment.

To run a single feature add `BEHAT_FEATURE=<feature file>`

example: `BEHAT_SUITE=tests/acceptance/features/apiGraph/createUser.feature`

To run a single test add `BEHAT_FEATURE=<file.feature:(line number)>`

example: `BEHAT_SUITE=tests/acceptance/features/apiGraph/createUser.feature:12`

To run a single suite add `BEHAT_SUITE=<test suite>`

example: `BEHAT_SUITE=apiGraph`

To run tests with a different storage driver set `STORAGE_DRIVER` to the correct value. It can be set to `OCIS` or `OWNCLOUD` and uses `OWNCLOUD` as the default value.

### use existing tests for BDD

As a lot of scenarios from `test-acceptance-from-core-api` are written for oC10, we can use those tests for Behaviour driven development in ocis.
Every scenario that does not work in oCIS with "ocis" storage, is listed in `tests/acceptance/expected-failures-API-on-OCIS-storage.md` with a link to the related issue.

Those scenarios are run in the ordinary acceptance test pipeline in CI. The scenarios that fail are checked against the
expected failures. If there are any differences then the CI pipeline fails.

The tests are not currently run in CI with the OWNCLOUD or EOS storage drivers, so there are no expected-failures files for those.

If you want to work on a specific issue

1. locally run each of the tests marked with that issue in the expected failures file.

    E.g.:

    ```bash
    make test-acceptance-from-core-api \
    TEST_SERVER_URL=https://localhost:9200 \
    TEST_OCIS=true \
    TEST_WITH_GRAPH_API=true \
    STORAGE_DRIVER=OCIS \
    BEHAT_FEATURE='tests/acceptance/features/coreApiVersions/fileVersions.feature:147'
    ```

2. the tests will fail, try to understand how and why they are failing
3. fix the code
4. go back to 1. and repeat till the tests are passing.
5. remove those tests from the expected failures file
6. make a PR that has the fixed code, and the relevant lines removed from the expected failures file.

## Running tests for parallel deployment

### Setup the parallel deployment environment

Instruction on setup is available [here](https://owncloud.dev/ocis/deployment/oc10_ocis_parallel/#local-setup)

Edit the `.env` file and uncomment this line:

```bash
COMPOSE_FILE=docker-compose.yml:testing/docker-compose-additions.yml
```

Start the docker stack with the following command:

```bash
docker-compose up -d
```

### Getting the test helpers

All the test helpers are located in the core repo.

```bash
git clone https://github.com/owncloud/core.git
```

### Run the acceptance tests

Run the acceptance tests with the following command from the root of the oCIS repository:

```bash
make test-paralleldeployment-api \
TEST_SERVER_URL="https://cloud.owncloud.test" \
TEST_OC10_URL="http://localhost:8080" \
TEST_PARALLEL_DEPLOYMENT=true \
TEST_OCIS=true \
TEST_WITH_LDAP=true \
PATH_TO_CORE="<path_to_core>" \
SKELETON_DIR="<path_to_core>/apps/testing/data/apiSkeleton"
```

Replace `<path_to_core>` with the actual path to the root directory of core repo that you have cloned earlier.

In order to run a single test, use the `BEHAT_FEATURE` environment variable.

```bash
make test-paralleldeployment-api \
... \
BEHAT_FEATURE="tests/parallelDeployAcceptance/features/apiShareManagement/acceptShares.feature"
```

## Running test suite with email service (inbucket)

### Setup inbucket

Run the following command to setup inbucket

```bash
docker run -d --name inbucket -p 9000:9000 -p 2500:2500 -p 1100:1100 inbucket/inbucket
```

### Run oCIS with following environment variables

Documentation for environment variables is available [here](https://owncloud.dev/services/notifications/#environment-variables)

```bash
OCIS_INSECURE=true \
PROXY_ENABLE_BASIC_AUTH=true \
NOTIFICATIONS_SMTP_HOST=localhost \
NOTIFICATIONS_SMTP_INSECURE=true \
NOTIFICATIONS_SMTP_PORT=2500 \
OCIS_URL=https://localhost:9200 \
<path_to_ocis>/ocis/bin/ocis server
```

### Run the acceptance test

Run the acceptance test with the following command:
```bash
make test-acceptance-api \
TEST_SERVER_URL="https://localhost:9200" \
TEST_OCIS=true \
TEST_WITH_GRAPH_API=true \
EMAIL_HOST="localhost" \
EMAIL_PORT=9000 \
BEHAT_FEATURE="tests/acceptance/features/apiEmailNotification/emailNotification.feature"
```
