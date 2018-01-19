package haproxy

//go:generate collectd-template-to-go haproxy.tmpl

import (
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/collectd"
	"github.com/signalfx/neo-agent/monitors/types"
)

const monitorType = "collectd/haproxy"

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

	Host string `yaml:"host"`
	Port uint16 `yaml:"port"`
	Name string `yaml:"name"`

	ProxiesToMonitor []string `yaml:"proxiesToMonitor"`
	ExcludedMetrics  []string `yaml:"excludedMetrics" default:"[]"`
	EnhancedMetrics  *bool    `yaml:"enhancedMetrics" default:"false"`
}

// Monitor is the main type that represents the monitor
type Monitor struct {
	collectd.MonitorCore
}

// Configure configures and runs the plugin in collectd
func (mm *Monitor) Configure(conf *Config) error {
	return mm.SetConfigurationAndRun(conf)
}
