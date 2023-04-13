# Proxy

The proxy service is an API-Gateway for the ownCloud Infinite Scale microservices. Every HTTP request goes through this service. Authentication, logging and other preprocessing of requests also happens here. Mechanisms like request rate limitting or intrusion prevention are **not** included in the proxy service and must be setup in front like with an external reverse proxy.

The proxy service is the only service communicating to the outside and needs therefore usual protections against DDOS, Slow Loris or other attack vectors. All other services are not exposed to the outside, but also need protective measures when it comes to distributed setups like when using container orchestration over various physical servers.

## Authentication

The following request authentication schemes are implemented:

-   Basic Auth (Only use in development, **never in production** setups!)
-   OpenID Connect
-   Signed URL
-   Public Share Token

## Automatic Quota Assignments

It is possible to automatically assign a specific quota to new users depending on their role.
To do this, you need to configure a mapping between roles defined by their ID and the quota in bytes.
The assignment can only be done via a `yaml` configuration and not via environment variables.
See the following `proxy.yaml` config snippet for a configuration example.

```yaml
role_quotas:
    <role ID1>: <quota1>
    <role ID2>: <quota2>
```

## Automatic Role Assignments

When users login, they do automatically get a role assigned. The automatic role assignment can be
configured in different ways. The `PROXY_ROLE_ASSIGNMENT_DRIVER` environment variable (or the `driver`
setting in the `role_assignment` section of the configuration file select which mechanism to use for
the automatic role assignment.

When set to `default`, all users which do not have a role assigned at the time for the first login will
get the role 'user' assigned. (This is also the default behavior if `PROXY_ROLE_ASSIGNMENT_DRIVER`
is unset.

When `PROXY_ROLE_ASSIGNMENT_DRIVER` is set to `oidc` the role assignment for a user will happen
based on the values of an OpenID Connect Claim of that user. The name of the OpenID Connect Claim to
be used for the role assignment can be configured via the `PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM`
environment variable. It is also possible to define a mapping of claim values to role names defined
in ownCloud Infinite Scale via a `yaml` configuration. See the following `proxy.yaml` snippet for an
example.

```yaml
role_assignment:
    driver: oidc
    oidc_role_mapper:
        role_claim: ocisRoles
        role_mapping:
            - role_name: admin
              claim_value: myAdminRole
            - role_name: spaceadmin
              claim_value: mySpaceAdminRole
            - role_name: user
              claim_value: myUserRole
            - role_name: guest:
              claim_value: myGuestRole
```

This would assign the role `admin` to users with the value `myAdminRole` in the claim `ocisRoles`.
The role `user` to users with the values `myUserRole` in the claims `ocisRoles` and so on.

Claim values that are not mapped to a specific ownCloud Infinite Scale role will be ignored.

Note: An ownCloud Infinite Scale user can only have a single role assigned. If the configured
`role_mapping` and a user's claim values result in multiple possible roles for a user, the order in
which the role mappings are defined in the configuration is important. The first role in the
`role_mappings` where the `claim_value` matches a value from the user's roles claim will be assigned
to the user. So if e.g. a user's `ocisRoles` claim has the values `myUserRole` and
`mySpaceAdminRole` that user will get the ocis role `spaceadmin` assigned (because `spaceadmin`
appears before `user` in the above sample configuration).

If a user's claim values don't match any of the configured role mappings an error will be logged and
the user will not be able to login.

The default `role_claim` (or `PROXY_ROLE_ASSIGNMENT_OIDC_CLAIM`) is `roles`. The default `role_mapping` is:

```yaml
- role_name: admin
  claim_value: ocisAdmin
- role_name: spaceadmin
  claim_value: ocisSpaceAdmin
- role_name: user
  claim_value: ocisUser
- role_name: guest:
  claim_value: ocisGuest
```

## Recommendations for Production Deployments

In a production deployment, you want to have basic authentication (`PROXY_ENABLE_BASIC_AUTH`) disabled which is the default state. You also want to setup a firewall to only allow requests to the proxy service or the reverse proxy if you have one. Requests to the other services should be blocked by the firewall.

## Caching

The `proxy` service can use a configured store via `PROXY_STORE_TYPE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured redis cluster.
  -   `redis-sentinel`: Stores data in a configured redis sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.

1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  The proxy service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
4.  When using `redis-sentinel`, the Redis master to use is configured via `PROXY_OIDC_USERINFO_CACHE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.
