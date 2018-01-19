package genericjmx

import (
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/types"
	yaml "gopkg.in/yaml.v2"
)

const monitorType = "collectd/genericjmx"

// Monitor is the main type that represents the monitor
type Monitor struct {
	*JMXMonitorCore
}

func init() {
	err := yaml.Unmarshal([]byte(defaultMBeanYAML), &DefaultMBeans)
	if err != nil {
		panic("YAML for GenericJMX MBeans is invalid: " + err.Error())
	}

	monitors.Register(monitorType, func(id types.MonitorID) interface{} {
		return Monitor{
			NewJMXMonitorCore(id, DefaultMBeans, "java"),
		}
	}, &Config{})
}
