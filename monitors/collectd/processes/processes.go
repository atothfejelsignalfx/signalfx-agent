package processes

//go:generate collectd-template-to-go processes.tmpl

import (
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/core/meta"
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/collectd"
	"github.com/signalfx/neo-agent/monitors/types"
)

const monitorType = "collectd/processes"

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
	Processes            []string          `yaml:"processes"`
	ProcessMatch         map[string]string `yaml:"processMatch"`
	CollectContextSwitch bool              `yaml:"collectContextSwitch" default:"false"`
	ProcFSPath           string            `yaml:"procFSPath"`
}

// Validate will check the config for correctness.
func (c *Config) Validate() error {
	return nil
}

// Monitor is the main type that represents the monitor
type Monitor struct {
	collectd.MonitorCore
	AgentMeta *meta.AgentMeta
}

// Configure configures and runs the plugin in collectd
func (am *Monitor) Configure(conf *Config) error {
	if conf.ProcFSPath == "" {
		conf.ProcFSPath = am.AgentMeta.ProcFSPath
	}

	return am.SetConfigurationAndRun(conf)
}
