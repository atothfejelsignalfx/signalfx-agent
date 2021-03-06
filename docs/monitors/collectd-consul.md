<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/consul

 Monitors the Consul data store by using the
[Consul collectd Python
plugin](https://github.com/signalfx/collectd-consul), which collects metrics
from Consul instances by hitting these endpoints:
- [/agent/self](https://www.consul.io/api/agent.html#read-configuration)
- [/agent/metrics](https://www.consul.io/api/agent.html#view-metrics)
- [/catalog/nodes](https://www.consul.io/api/catalog.html#list-nodes)
- [/catalog/node/:node](https://www.consul.io/api/catalog.html#list-services-for-node)
- [/status/leader](https://www.consul.io/api/status.html#get-raft-leader)
- [/status/peers](https://www.consul.io/api/status.html#list-raft-peers)
- [/coordinate/datacenters](https://www.consul.io/api/coordinate.html#read-wan-coordinates)
- [/coordinate/nodes](https://www.consul.io/api/coordinate.html#read-lan-coordinates)
- [/health/state/any](https://www.consul.io/api/health.html#list-checks-in-state)


Monitor Type: `collectd/consul`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/consul)

**Accepts Endpoints**: **Yes**

**Multiple Instances Allowed**: Yes

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `host` | **yes** | `string` |  |
| `port` | **yes** | `integer` |  |
| `aclToken` | no | `string` |  |
| `useHTTPS` | no | `bool` |  (**default:** `false`) |
| `enhancedMetrics` | no | `bool` |  (**default:** `false`) |
| `caCertificate` | no | `string` |  |
| `clientCertificate` | no | `string` |  |
| `clientKey` | no | `string` |  |
| `signalFxAccessToken` | no | `string` |  |




## Metrics

This monitor emits the following metrics.  Note that configuration options may
cause only a subset of metrics to be emitted.

| Name | Type | Description |
| ---  | ---  | ---         |
| `consul.dns.stale_queries` | gauge | Number of times an agent serves a DNS query with stale information |
| `consul.memberlist.msg.suspect` | gauge | Number of suspect messages received per interval |
| `consul.serf.member.flap` | gauge | Tracks flapping agents |
| `gauge.consul.catalog.nodes.total` | gauge | Number of nodes in the Consul datacenter |
| `gauge.consul.catalog.nodes_by_service` | gauge | Number of nodes providing a given service |
| `gauge.consul.catalog.services.total` | gauge | Total number of services registered with Consul in the given datacenter |
| `gauge.consul.catalog.services_by_node` | gauge | Number of services registered with a node |
| `gauge.consul.consul.dns.domain_query.AGENT.avg` | gauge | Average time to complete a forward DNS query |
| `gauge.consul.consul.dns.domain_query.AGENT.max` | gauge | Max time to complete a forward DNS query |
| `gauge.consul.consul.dns.domain_query.AGENT.min` | gauge | Min time to complete a forward DNS query |
| `gauge.consul.consul.dns.ptr_query.AGENT.avg` | gauge | Average time to complete a Reverse DNS query |
| `gauge.consul.consul.dns.ptr_query.AGENT.max` | gauge | Max time to complete a Reverse DNS query |
| `gauge.consul.consul.dns.ptr_query.AGENT.min` | gauge | Min time to complete a Reverse DNS query |
| `gauge.consul.consul.leader.reconcile.avg` | gauge | Leader time to reconcile the differences between Serf membership and Consul's store |
| `gauge.consul.health.nodes.critical` | gauge | Number of nodes for which health checks are reporting Critical state |
| `gauge.consul.health.nodes.passing` | gauge | Number of nodes for which health checks are reporting Passing state |
| `gauge.consul.health.nodes.warning` | gauge | Number of nodes for which health checks are reporting Warning state |
| `gauge.consul.health.services.critical` | gauge | Number of services for which health checks are reporting Critical state |
| `gauge.consul.health.services.passing` | gauge | Number of services for which health checks are reporting Passing state |
| `gauge.consul.health.services.warning` | gauge | Number of services for which health checks are reporting Warning state |
| `gauge.consul.is_leader` | gauge | Metric to map consul server's in leader or follower state |
| `gauge.consul.network.dc.latency.avg` | gauge | Average network latency between 2 datacenters |
| `gauge.consul.network.dc.latency.max` | gauge | Maximum network latency between 2 datacenters |
| `gauge.consul.network.dc.latency.min` | gauge | Minimum network latency between 2 datacenters |
| `gauge.consul.network.node.latency.avg` | gauge | Average network latency between given node and other nodes in the datacenter |
| `gauge.consul.network.node.latency.max` | gauge | Minimum network latency between given node and other nodes in the datacenter |
| `gauge.consul.network.node.latency.min` | gauge | Minimum network latency between given node and other nodes in the datacenter |
| `gauge.consul.peers` | gauge | Number of Raft peers in Consul datacenter |
| `gauge.consul.raft.apply` | gauge | Number of raft transactions |
| `gauge.consul.raft.commitTime.avg` | gauge | Average of the time it takes to commit an entry on the leader |
| `gauge.consul.raft.commitTime.max` | gauge | Max of the time it takes to commit an entry on the leader |
| `gauge.consul.raft.commitTime.min` | gauge | Minimum of the time it takes to commit an entry on the leader |
| `gauge.consul.raft.leader.dispatchLog.avg` | gauge | Average of the time it takes for the leader to write log entries to disk |
| `gauge.consul.raft.leader.dispatchLog.max` | gauge | Maximum of the time it takes for the leader to write log entries to disk |
| `gauge.consul.raft.leader.dispatchLog.min` | gauge | Minimum of the time it takes for the leader to write log entries to disk |
| `gauge.consul.raft.leader.lastContact.avg` | gauge | Mean of the time since the leader was last able to contact follower nodes |
| `gauge.consul.raft.leader.lastContact.max` | gauge | Max of the time since the leader was last able to contact follower nodes |
| `gauge.consul.raft.leader.lastContact.min` | gauge | Min of the time since the leader was last able to contact follower nodes |
| `gauge.consul.raft.replication.appendEntries.rpc.AGENT.avg` | gauge | Mean time taken to complete the AppendEntries RPC |
| `gauge.consul.raft.replication.appendEntries.rpc.AGENT.max` | gauge | Max time taken to complete the AppendEntries RPC |
| `gauge.consul.raft.replication.appendEntries.rpc.AGENT.min` | gauge | Min time taken to complete the AppendEntries RPC |
| `gauge.consul.raft.state.candidate` | gauge | Tracks the number of times given node enters the candidate state |
| `gauge.consul.raft.state.leader` | gauge | Tracks the number of leadership transitions per interval |
| `gauge.consul.runtime.alloc_bytes` | gauge | Number of bytes allocated to Consul process on the node |
| `gauge.consul.runtime.heap_objects` | gauge | Number of heap objects allocated to Consul |
| `gauge.consul.runtime.num_goroutines` | gauge | Number of GO routines run by Consul process |
| `gauge.consul.serf.events` | gauge | Number of serf events processed |
| `gauge.consul.serf.member.join` | gauge | Tracks successful node joins |
| `gauge.consul.serf.member.left` | gauge | Tracks successful node leaves |
| `gauge.consul.serf.queue.Event.avg` | gauge | Average number of serf events in queue yet to be processed |
| `gauge.consul.serf.queue.Event.max` | gauge | Maximum number of serf events in queue yet to be processed during the interval |
| `gauge.consul.serf.queue.Event.min` | gauge | Minimum number of serf events in queue yet to be processed during the interval |
| `gauge.consul.serf.queue.Query.avg` | gauge | Average number of serf queries in queue yet to be processed during the interval |
| `gauge.consul.serf.queue.Query.max` | gauge | Maximum number of serf queries in queue yet to be processed during the interval |
| `gauge.consul.serf.queue.Query.min` | gauge | Minimum number of serf queries in queue yet to be processed during the interval |

## Dimensions

The following dimensions may occur on metrics emitted by this monitor.  Some
dimensions may be specific to certain metrics.

| Name | Description |
| ---  | ---         |
| `consul_mode` | Whether this consul instance is running as a server or client |
| `consul_node` | The name of the consul node |
| `datacenter` | The name of the consul datacenter |



