package metadata

//go:generate collectd-template-to-go metadata.tmpl

import (
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/core/meta"
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/collectd"
	"github.com/signalfx/neo-agent/monitors/types"
)

const monitorType = "collectd/signalfx-metadata"

func init() {
	monitors.Register(monitorType, func(id types.MonitorID) interface{} {
		return &Monitor{
			MonitorCore: *collectd.NewMonitorCore(id, CollectdTemplate),
		}
	}, &Config{})
}

// Config is the monitor-specific config with the generic config embedded
type Config struct {
	config.MonitorConfig
	WriteServerURL string `yaml:"-"`
	ProcFSPath     string `yaml:"-"`
}

// Monitor is the main type that represents the monitor
type Monitor struct {
	collectd.MonitorCore
	AgentMeta *meta.AgentMeta
}

// Configure configures and runs the plugin in collectd
func (m *Monitor) Configure(conf *Config) error {
	conf.WriteServerURL = m.AgentMeta.CollectdConf.WriteServerURL()
	conf.ProcFSPath = m.AgentMeta.ProcFSPath

	return m.SetConfigurationAndRun(conf)
}
