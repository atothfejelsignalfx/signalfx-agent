<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/haproxy


**Accepts Endpoints**: **Yes**

**Only One Instance Allowed**: No

## Configuration

| Config option | Default | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `host` |  | **yes** | `string` |  |
| `port` | `0` | **yes** | `uint16` |  |
| `name` |  | no | `string` |  |
| `proxiesToMonitor` | `[]` | no | `slice` |  |
| `excludedMetrics` | `[]` | no | `slice` |  |
| `enhancedMetrics` | `false` | no | `bool` |  |











