<!--- GENERATED BY gomplate from scripts/docs/monitor-page.md.tmpl --->

# collectd/disk

 This monitor collects information about the usage of
physical disks and logical disks (partitions).

See https://collectd.org/wiki/index.php/Plugin:Disk.


Monitor Type: `collectd/disk`

[Monitor Source Code](https://github.com/signalfx/signalfx-agent/tree/master/internal/monitors/collectd/disk)

**Accepts Endpoints**: No

**Multiple Instances Allowed**: **No**

## Configuration

| Config option | Required | Type | Description |
| --- | --- | --- | --- |
| `disks` | no | `list of string` | Which devices to include/exclude (**default:** `[/^loop\d+$/ /^dm-\d+$/]`) |
| `ignoreSelected` | no | `bool` | If true, the disks selected by `disks` will be excluded and all others included. (**default:** `true`) |






