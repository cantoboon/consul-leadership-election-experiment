# Consul Leadership Election Example

A very basic look at how to do Leadership Election with Consul.

We're using [HashiCorp's Application Leader Election with Sessions](https://learn.hashicorp.com/tutorials/consul/application-leader-elections)
as a guide.

To simplify the architecture, we're going to have several services in a single process. These services
will individually create sessions and attempt to modify a key.

## Dev Setup

Run Docker Compose to start Consul.

Open: [http://localhost:8500/ui/](http://localhost:8500/ui/)

## TTL Example

[cmd/ttl-example/main.go](cmd/ttl-example/main.go) shows the TTL in action. Using the lowest TTL, 10s,
we can see that after the TTL expires a new session can get the lock.