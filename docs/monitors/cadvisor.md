<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# cadvisor

This monitor pulls metrics directly from cadvisor, which
conventionally runs on port 4194, but can be configured to anything.


**Accepts Endpoints**: No

**Only One Instance Allowed**: No

## Configuration

| Config option | Default | Required | Type | Description |
| --- | --- | --- | --- | --- |
| `cadvisorURL` | `http://localhost:4194` | no | `string` | Where to find cAdvisor |





## Metrics

| Name | Type | Description |
| ---  | ---  | ---         |
| container_last_seen | timestamp | Last time a container was seen by the exporter |
| container_cpu_user_seconds_total | counter | Cumulative user cpu time consumed in nanoseconds |

## Dimensions

| Name | Description |
| ---  | ---         |
| kubernetes_namespace | The K8s namespace the container is part of |
| kubernetes_pod_name | The pod instance under which this container runs |
| kubernetes_pod_uid | The UID of the pod instance under which this container runs |
| container_spec_name | The container's name as it appears in the pod spec |
| container_name | The container's name as it appears in the pod spec, the same as container_spec_name but retained for backwards compatibility. |
| container_id | The ID of the running container |
| container_image | The container image name |


