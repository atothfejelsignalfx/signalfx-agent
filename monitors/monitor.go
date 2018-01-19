package monitors

import (
	"github.com/pkg/errors"
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/core/services"
	"github.com/signalfx/neo-agent/monitors/types"
	"github.com/signalfx/neo-agent/utils"
	log "github.com/sirupsen/logrus"
)

// MonitorFactory is a niladic function that creates an unconfigured instance
// of a monitor.
type MonitorFactory func(types.MonitorID) interface{}

// MonitorFactories holds all of the registered monitor factories
var MonitorFactories = map[string]MonitorFactory{}

// These are blank (zero-value) instances of the configuration struct for a
// particular monitor type.
var configTemplates = map[string]config.MonitorCustomConfig{}

// InjectableMonitor should be implemented by a dynamic monitor that needs to
// know when services are added and removed.
type InjectableMonitor interface {
	AddService(services.Endpoint)
	RemoveService(services.Endpoint)
}

// Register a new monitor type with the agent.  This is intended to be called
// from the init function of the module of a specific monitor
// implementation. configTemplate should be a zero-valued struct that is of the
// same type as the parameter to the Configure method for this monitor type.
func Register(_type string, factory MonitorFactory, configTemplate config.MonitorCustomConfig) {
	if _, ok := MonitorFactories[_type]; ok {
		panic("Monitor type '" + _type + "' already registered")
	}
	MonitorFactories[_type] = factory
	configTemplates[_type] = configTemplate
}

// DeregisterAll unregisters all monitor types.  Primarily intended for testing
// purposes.
func DeregisterAll() {
	for k := range MonitorFactories {
		delete(MonitorFactories, k)
	}

	for k := range configTemplates {
		delete(configTemplates, k)
	}
}

func newUninitializedMonitor(_type string, id types.MonitorID) interface{} {
	if factory, ok := MonitorFactories[_type]; ok {
		return factory(id)
	}

	log.WithFields(log.Fields{
		"monitorType": _type,
	}).Error("Monitor type not supported")
	return nil
}

// Creates a new, unconfigured instance of a monitor of _type.  Returns nil if
// the monitor type is not registered.
func newMonitor(_type string, id types.MonitorID) interface{} {
	mon := newUninitializedMonitor(_type, id)
	if initMon, ok := mon.(Initializable); ok {
		if err := initMon.Init(); err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"monitorType": _type,
			}).Error("Could not initialize monitor")
			return nil
		}
	}
	return mon

}

// Initializable represents a monitor that has a distinct InitMonitor method.
// This should be called once after the monitor is created and before any of
// its other methods are called.  It is useful for things that are not
// appropriate to do in the monitor factory function.
type Initializable interface {
	Init() error
}

// Shutdownable should be implemented by all monitors that need to clean up
// resources before being destroyed.
type Shutdownable interface {
	Shutdown()
}

// Takes a generic MonitorConfig and pulls out monitor-specific config to
// populate a clone of the config template that was registered for the monitor
// type specified in conf.  This will also validate the config and return nil
// if validation fails.
func getCustomConfigForMonitor(conf *config.MonitorConfig) (config.MonitorCustomConfig, error) {
	confTemplate, ok := configTemplates[conf.Type]
	if !ok {
		return nil, errors.Errorf("Unknown monitor type %s", conf.Type)
	}
	monConfig := utils.CloneInterface(confTemplate).(config.MonitorCustomConfig)

	if err := config.FillInConfigTemplate("MonitorConfig", monConfig, conf); err != nil {
		return nil, err
	}

	// These methods will set state inside the config such that conf.IsValid
	// will return true or false
	if err := validateConfig(monConfig); err != nil {
		return monConfig, err
	}

	return monConfig, nil
}

func anyMarkedSolo(confs []config.MonitorConfig) bool {
	for i := range confs {
		if confs[i].Solo {
			return true
		}
	}
	return false
}
