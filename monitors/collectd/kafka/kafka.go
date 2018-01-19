package kafka

import (
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/collectd/genericjmx"
	"github.com/signalfx/neo-agent/monitors/types"
	yaml "gopkg.in/yaml.v2"
)

const monitorType = "collectd/kafka"

var serviceName = "kafka"

// Monitor is the main type that represents the monitor
type Monitor struct {
	*genericjmx.JMXMonitorCore
}

func init() {
	var defaultMBeans genericjmx.MBeanMap
	err := yaml.Unmarshal([]byte(defaultMBeanYAML), &defaultMBeans)
	if err != nil {
		panic("YAML for GenericJMX MBeans is invalid: " + err.Error())
	}
	defaultMBeans = defaultMBeans.MergeWith(genericjmx.DefaultMBeans)

	monitors.Register(monitorType, func(id types.MonitorID) interface{} {
		return Monitor{
			genericjmx.NewJMXMonitorCore(id, defaultMBeans, serviceName),
		}
	}, &genericjmx.Config{})
}
