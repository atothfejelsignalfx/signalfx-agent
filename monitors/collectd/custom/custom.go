// Package custom contains a custom collectd plugin monitor, for which you can
// specify your own config template and parameters.
package custom

import (
	"text/template"

	"github.com/pkg/errors"
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/collectd"
	"github.com/signalfx/neo-agent/monitors/collectd/templating"
	"github.com/signalfx/neo-agent/monitors/types"
)

const monitorType = "collectd/custom"

func init() {
	monitors.Register(monitorType, func(id types.MonitorID) interface{} {
		return &Monitor{
			MonitorCore: *collectd.NewMonitorCore(id, template.New("custom")),
		}
	}, &Config{})
}

// Config is the configuration for the collectd custom monitor
type Config struct {
	config.MonitorConfig `acceptsEndpoints:"true"`

	Host string  `yaml:"host"`
	Port uint16  `yaml:"port"`
	Name *string `yaml:"name"`

	TemplateText string `yaml:"templateText"`
}

// Validate will check the config that is specific to this monitor
func (c *Config) Validate() error {
	if c.TemplateText == "" {
		return errors.New("templateText must be set")
	}
	if _, err := templateFromText(c.TemplateText); err != nil {
		return err
	}
	return nil
}

func templateFromText(templateText string) (*template.Template, error) {
	template, err := templating.InjectTemplateFuncs(template.New("custom")).Parse(templateText)
	if err != nil {
		return nil, errors.Wrapf(err, "Template text failed to parse: \n%s", templateText)
	}
	return template, nil
}

// Monitor is the core monitor object that gets instantiated by the agent
type Monitor struct {
	collectd.MonitorCore
}

// Configure will render the custom collectd config and queue a collectd
// restart.
func (cm *Monitor) Configure(conf *Config) error {
	var err error
	cm.Template, err = templateFromText(conf.TemplateText)
	if err != nil {
		return err
	}

	return cm.SetConfigurationAndRun(conf)
}
