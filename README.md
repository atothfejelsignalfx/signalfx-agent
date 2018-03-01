# SignalFx Agent

The SignalFx Agent is a metric agent written in Go for monitoring nodes and
application services in a variety of different environments.

## Concepts

The agent has three main components:

1) _Observers_ that discover applications and services running on the host
2) _Monitors_ that collect metrics from the host and applications
3) The _Writer_ that sends the metrics collected by monitors to SignalFx.

### Observers

Observers are what watch the various environments we support to discover running
services and automatically configure the agent to send metrics for those
services.

See [Observer Config](./docs/observer-config.md) for a list of supported
observers and their configuration.

### Monitors

Monitors are what collect metrics from the host system and services.  They are
configured under the `monitors` list in the agent config.  For monitors of
applications, you can configure a discovery rule on the monitor so that a
separate instance of the monitor is created for each discovered instance of
applications that match that discovery rule.  See [Auto
Discovery](./docs/auto-discovery.md) for more information.

Many of the monitors rely on a third-party "super monitor",
[collectd](https://collectd.org), under the covers to do a lot of the metric
collection, although we also have monitors apart from Collectd.  They are
configured in the same way, however.

See [Monitor Config](./docs/monitor-config.md) for a list of supported monitors
and their configuration.


## Configuration

The agent is configured primarily from a YAML file (by default,
`/etc/signalfx/agent.yaml`, but this can be overridden by the `-config` command
line flag).  

For the full schema of the config, see [Config Schema](./docs/config-schema.md).

See [Remote Configuration](./docs/remote-config.md) for information on how to
configure the agent from remote sources such as other files on the filesystem
or KV stores such as Etcd.

## Installation

The agent is available in both a containerized and standalone form for Linux.
Whatever form you use, the dependencies are completely bundled along with the
agent, including a Java JRE runtime and a Python runtime, so there are no
additional dependencies required.  This means that the agent should work on any
relatively modern Linux distribution (kernel version 2.6+).

### Bundles
We offer the agent in the following forms:

#### Docker Image
The agent is available as a Docker image at
[quay.io/signalfx/signalfx-agent](https://quay.io/signalfx/signalfx-agent). The
image is tagged using the same agent version scheme.

#### Debian Package
We provide a Debian package repository that you can make use of with the
following commands:

```sh
curl -sSL https://dl.signalfx.com/debian.gpg > /etc/apt/trusted.gpg.d/signalfx.gpg
echo 'deb http://dl.signalfx.com/debs/signalfx-agent/main /' > /etc/apt/sources.list.d/signalfx-agent.list
apt-get update
apt-get install -y signalfx-agent
```

#### RPM Package
We provide a RHEL/RPM package repository that you can make use of with the
following commands:

```sh
cat <<EOH > /etc/yum.repos.d/signalfx-agent.repo
[signalfx-agent]
name=SignalFx Agent Repository
baseurl=http://dl.signalfx.com/rpms/signalfx-agent/main
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://dl.signalfx.com/yum-rpm.key
enabled=1
EOH

yum install -y signalfx-agent
```

#### Standalone tar.gz
For Linux distros without a supported package, we offer the agent bundle as a
.tar.gz that can be deployed to the target host.  This bundle is available for
download on the [Github Releases
Page](https://github.com/signalfx/signalfx-agent/releases) for each new
release.

To use the bundle, simply 

1) Unarchive it to any directory on the target system you like
2) Ensure a valid configuration file is available somewhere on the target
system
3) Run the agent by invoking `signalfx-agent/bin/signalfx-agent -config <path
to config.yaml>`.

### Deployment Tools
We support the following deployment/configuration management tools to automate the
installation process:

#### Installer Script
For non-containerized environments, there is a convenience script that you can
run on your host to install the agent package.  This is useful for testing and
trails, but for full-scale deployments you will probably want to use a
configuration management system like Chef or Puppet.

#### Chef
We offer a Chef cookbook to install and configure the agent.  See [the cookbook
source](./deployments/chef) and INSERT URL TO SUPERMARKET.

#### Puppet
We also offer a Puppet manifest to install and configure the agent.  See [the
manifest source](./deployments/puppet) and INSERT MORE.

#### Kubernetes
See our [Kubernetes Quickstart
Guide](https://docs.signalfx.com/en/latest/integrations/kubernetes-quickstart.html)
for more information.

### Privileges

When using the [`host` observer](./docs/observers/host.md), the agent requires
the [Linux
capabilities](http://man7.org/linux/man-pages/man7/capabilities.7.html)
`DAC_READ_SEARCH` and `SYS_PTRACE`, both of which are necessary to allow the
agent to determine which processes are listening on network ports on the host.
Otherwise, there is nothing built into the agent that requires privileges.
When using a package to install the agent, the agent binary is given those
capabilities in the package post-install script, but the agent is run as the
`signalfx-agent` user.  If you are not using the `host` observer, then you can
strip those capabilities from the agent binary if so desired.

You should generally not run the agent as `root` unless you can't use
capabilities for some reason.

## Logging
Currently the agent only supports logging to stdout, which will generally be
redirected by the init scripts we provide to either a file at
`/var/log/signalfx-agent.log` or to the systemd journal on newer distros. The
default log level is `info`, which will log anything noteworthy in the agent
without spamming the logs too much.  Most of the `info` level logs are on
startup and upon service discovery changes.  `debug` will create very verbose
log output and should only be used when trying to resolve a problem with the
agent.

## Proxy Support

To use a HTTP(S) proxy, set the environment variable `HTTP_PROXY` and/or
`HTTPS_PROXY` in the container configuration to proxy either protocol.  The
SignalFx ingest and API servers both use HTTPS.  The agent will automatically
manipulate the `NO_PROXY` envvar to not use the proxy for local services.

## Diagnostics
The agent serves diagnostic information on a unix domain socket at
`/var/run/signalfx.sock`.  The socket takes no input, but simply dumps it's
current status back upon connection.  The command `signalfx-agent status` (or
the special symlink `agent-status`) will read this socket and dump out its
contents.

## Development

If you wish to contribute to the agent, see the [Developer's
Guide](./docs/development.md).

