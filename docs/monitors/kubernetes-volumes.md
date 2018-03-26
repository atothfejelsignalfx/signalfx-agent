<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# kubernetes-volumes

 This monitor sends usage stats about volumes
mounted to Kubernetes pods (e.g. free space/inodes).  This information is
gotten from the Kubelet /stats/summary endpoint.  The normal `collectd/df`
monitor generally will not report Persistent Volume usage metrics because
those volumes are not seen by the agent since they can be mounted
dynamically and older versions of K8s don't support mount propagation of
those mounts to the agent container.


Monitor Type: `kubernetes-volumes`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/kubernetes/volumes)

**Accepts Endpoints**: No

**Multiple Instances Allowed**: Yes

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `kubeletAPI` | no | `object (see below)` | Kubelet client configuration |


The **nested** `kubeletAPI` config object has the following fields:

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `url` | no | `string` | URL of the Kubelet instance.  This will default to `https://<current node hostname>:10250` if not provided. |
| `authType` | no | `string` | Can be `none` for no auth, `tls` for TLS client cert auth, or `serviceAccount` to use the pod's default service account token to authenticate. (**default:** `none`) |
| `skipVerify` | no | `bool` | Whether to skip verification of the Kubelet's TLS cert (**default:** `false`) |
| `caCertPath` | no | `string` | Path to the CA cert that has signed the Kubelet's TLS cert, unnecessary if `skipVerify` is set to false. |
| `clientCertPath` | no | `string` | Path to the client TLS cert to use if `authType` is set to `tls` |
| `clientKeyPath` | no | `string` | Path to the client TLS key to use if `authType` is set to `tls` |
| `logResponses` | no | `bool` | Whether to log the raw cadvisor response at the debug level for debugging purposes. (**default:** `false`) |




## Dimensions

The following dimensions may occur on metrics emitted by this monitor.  Some
dimensions may be specific to certain metrics.

| Name | Description |
| ---  | ---         |
| `kubernetes_namespace` | The namespace of the pod that has this volume |
| `kubernetes_pod_name` | The name of the pod that has this volume |
| `kubernetes_pod_uid` | The UID of the pod that has this volume |
| `volume` | The volume name as given in the pod spec under `volumes` |


