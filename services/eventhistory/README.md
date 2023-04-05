# Eventhistory Service

The `eventhistory` consumes all events from the configured event system like NATS, stores them and allows other services to retrieve them via an eventid.

## Prerequisites

Running the eventhistory service without an event sytem like NATS is not possible.

## Consuming

The `eventhistory` services consumes all events from the configured event sytem.

## Storing

The `eventhistory` service stores each consumed event via the configured store in `EVENTHISTORY_STORE_TYPE`. Possible stores are:
  -   `memory`: Basic in-memory store and the default.
  -   `ocmem`: Advanced in-memory store allowing max size.
  -   `redis`: Stores data in a configured redis cluster.
  -   `redis-sentinel`: Stores data in a configured redis sentinel cluster.
  -   `etcd`: Stores data in a configured etcd cluster.
  -   `nats-js`: Stores data using key-value-store feature of [nats jetstream](https://docs.nats.io/nats-concepts/jetstream/key-value-store)
  -   `noop`: Stores nothing. Useful for testing. Not recommended in productive enviroments.

1.  Note that in-memory stores are by nature not reboot persistent.
2.  Though usually not necessary, a database name and a database table can be configured for event stores if the event store supports this. Generally not applicapable for stores of type `in-memory`. These settings are blank by default which means that the standard settings of the configured store applies.
3.  Events stay in the store for 2 weeks by default. Use `EVENTHISTORY_RECORD_EXPIRY` to adjust this value.
4.  The eventhistory service can be scaled if not using `in-memory` stores and the stores are configured identically over all instances.
5.  When using `redis-sentinel`, the Redis master to use is configured via `EVENTHISTORY_STORE_NODES` in the form of `<sentinel-host>:<sentinel-port>/<redis-master>` like `10.10.0.200:26379/mymaster`.

## Retrieving

Other services can call the `eventhistory` service via a grpc call to retrieve events. The request must contain the eventid that should be retrieved.
