package marathon

//go:generate collectd-template-to-go marathon.tmpl

import (
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/collectd"
	"github.com/signalfx/neo-agent/monitors/types"
)

const monitorType = "collectd/marathon"

func init() {
	monitors.Register(monitorType, func(id types.MonitorID) interface{} {
		return &Monitor{
			*collectd.NewMonitorCore(id, CollectdTemplate),
		}
	}, &Config{})
}

// Config is the monitor-specific config with the generic config embedded
type Config struct {
	config.MonitorConfig `acceptsEndpoints:"true"`

	Host     string  `yaml:"host"`
	Port     uint16  `yaml:"port"`
	Name     string  `yaml:"name"`
	Username string  `yaml:"username"`
	Password *string `yaml:"password"`
}

// Monitor is the main type that represents the monitor
type Monitor struct {
	collectd.MonitorCore
}

// Configure configures and runs the plugin in collectd
func (am *Monitor) Configure(conf *Config) error {
	return am.SetConfigurationAndRun(conf)
}
