<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/mysql


**Accepts Endpoints**: **Yes**

**Only One Instance Allowed**: No

## Configuration

| Config option | Default | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `host` |  | **yes** | `string` |  |
| `port` | `0` | **yes** | `uint16` |  |
| `name` |  | no | `string` |  |
| `databases` | `[]` | no | `slice` |  |
| `username` |  | no | `string` | These credentials serve as defaults for all databases if not overridden |
| `password` |  | no | `string` |  |
| `reportHost` | `false` | no | `bool` |  |





The `databases` config object has the following fields:

| Config option | Default | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `name` |  | **yes** | `string` |  |
| `username` |  | no | `string` |  |
| `password` |  | no | `string` |  |







