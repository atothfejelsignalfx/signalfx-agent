package monitors

import (
	"reflect"

	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/golib/event"
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/core/meta"
	"github.com/signalfx/neo-agent/core/services"
	"github.com/signalfx/neo-agent/core/writer"
	"github.com/signalfx/neo-agent/monitors/collectd"
	"github.com/signalfx/neo-agent/monitors/types"
	"github.com/signalfx/neo-agent/utils"
	log "github.com/sirupsen/logrus"
)

// MonitorManager coordinates the startup and shutdown of monitors based on the
// configuration provided by the user.  Monitors that have discovery rules can
// be injected with multiple services.  If a monitor does not have a discovery
// rule (a "static" monitor), it will be started immediately (as soon as
// Configure is called).
type MonitorManager struct {
	monitorConfigs []config.MonitorCustomConfig
	// Keep track of which services go with which monitor
	activeMonitors []*ActiveMonitor
	badConfigs     []*config.MonitorConfig
	lock           sync.Mutex
	// Map of service endpoints that have been discovered
	discoveredEndpoints map[services.ID]services.Endpoint

	dpChan      chan<- *datapoint.Datapoint
	eventChan   chan<- *event.Event
	dimPropChan chan<- *writer.DimProperties

	agentMeta *meta.AgentMeta

	idGenerator func() string
}

// NewMonitorManager creates a new instance of the MonitorManager
func NewMonitorManager(agentMeta *meta.AgentMeta) *MonitorManager {
	return &MonitorManager{
		monitorConfigs:      make([]config.MonitorCustomConfig, 0),
		activeMonitors:      make([]*ActiveMonitor, 0),
		badConfigs:          make([]*config.MonitorConfig, 0),
		discoveredEndpoints: make(map[services.ID]services.Endpoint),
		idGenerator:         utils.NewIDGenerator(),
		agentMeta:           agentMeta,
	}
}

// Configure receives a list of monitor configurations.  It will start up any
// static monitors and watch discovered services to see if any match dynamic
// monitors.
func (mm *MonitorManager) Configure(confs []config.MonitorConfig) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	// By configuring collectd with the monitor manager, we absolve the monitor
	// instances of having to know about collectd config, which makes it easier
	// to create monitor config from disparate sources such as from observers.
	if err := collectd.Instance().Configure(mm.agentMeta.CollectdConf); err != nil {
		log.WithFields(log.Fields{
			"error":          err,
			"collectdConfig": spew.Sdump(mm.agentMeta.CollectdConf),
		}).Error("Could not configure collectd")
	}

	mm.monitorConfigs = make([]config.MonitorCustomConfig, 0, len(confs))

	requireSoloTrue := anyMarkedSolo(confs)

	// All monitors are marked for deletion at first.  They can be saved and
	// reused by having a compatible config in the newly provided config
	mm.markAllRegularMonitorsAsDoomed()

	for i := range confs {
		conf := &confs[i]

		if isConfigUnique(conf, confs[:i]) {
			log.WithFields(log.Fields{
				"monitorType": conf.Type,
				"config":      conf,
			}).Error("Monitor config is duplicated")
			continue
		}

		if requireSoloTrue && !conf.Solo {
			log.Infof("Solo mode is active, skipping monitor of type %s", conf.Type)
			continue
		}

		monConfig, err := mm.processConfigForMonitor(conf)
		if err != nil {
			log.WithFields(log.Fields{
				"monitorType": conf.Type,
				"error":       err,
			}).Error("Could not process configuration for monitor")
			conf.ValidationError = err.Error()
			mm.badConfigs = append(mm.badConfigs, conf)
			continue
		}

		mm.monitorConfigs = append(mm.monitorConfigs, monConfig)
	}

	mm.deleteDoomedMonitors()
}

func (mm *MonitorManager) processConfigForMonitor(conf *config.MonitorConfig) (config.MonitorCustomConfig, error) {
	monConfig, err := getCustomConfigForMonitor(conf)
	if err != nil {
		return nil, err
	}

	configMatchedActive := false

	// Go through each actively running monitor marked for deletion and see if
	// it matches this particular config. If a monitor is already configured
	// with this same configuration, it will be saved and left untouched. This
	// makes config O(n*m) where n=<number of monitor configs> and m=<number of
	// active monitors> but given how infrequent config changes are and how few
	// active monitors there will be for a single agent, this should be
	// acceptable.
	for i := range mm.activeMonitors {
		am := mm.activeMonitors[i]
		if am.doomed {
			if reflect.DeepEqual(am.config, monConfig) {
				log.WithFields(log.Fields{
					"config": spew.Sdump(am.config),
					"id":     am.id,
				}).Debug("Saving monitor instance")

				configMatchedActive = true
				am.doomed = false
			}
		}
	}

	if configMatchedActive {
		return monConfig, nil
	}

	// No discovery rule means that the monitor should run from the start
	if conf.DiscoveryRule == "" {
		return monConfig, mm.createAndConfigureNewMonitor(monConfig, nil)
	}

	mm.makeMonitorsForMatchingEndpoints(monConfig)
	// We need to go and see if any discovered endpoints should be
	// monitored by this config, if they aren't already.
	return monConfig, nil
}

// SetDPChannel allows you to inject a new datapoint channel to the manager and
// have it propagated to all active monitors
func (mm *MonitorManager) SetDPChannel(dpChan chan<- *datapoint.Datapoint) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.dpChan = dpChan

	for i := range mm.activeMonitors {
		mm.activeMonitors[i].injectDatapointChannelIfNeeded(dpChan)
	}
}

// SetEventChannel allows you to inject a new event channel to the manager and
// have it propagated to all active monitors
func (mm *MonitorManager) SetEventChannel(eventChan chan<- *event.Event) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.eventChan = eventChan

	for i := range mm.activeMonitors {
		mm.activeMonitors[i].injectEventChannelIfNeeded(eventChan)
	}
}

// SetDimPropChannel allows you to inject a new dimension property channel to
// the manager and have it propagated to all active monitors
func (mm *MonitorManager) SetDimPropChannel(dimPropChan chan<- *writer.DimProperties) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	mm.dimPropChan = dimPropChan

	for i := range mm.activeMonitors {
		mm.activeMonitors[i].injectDimPropertiesChannelIfNeeded(dimPropChan)
	}
}

// Marks all of the monitors without self-configured endpoints as doomed.
func (mm *MonitorManager) markAllRegularMonitorsAsDoomed() {
	for i := range mm.activeMonitors {
		if !mm.activeMonitors[i].HasSelfConfiguredEndpoint() {
			mm.activeMonitors[i].doomed = true
		}
	}
}

func (mm *MonitorManager) deleteDoomedMonitors() {
	newActiveMonitors := []*ActiveMonitor{}

	for i := range mm.activeMonitors {
		am := mm.activeMonitors[i]
		if am.doomed {
			log.WithFields(log.Fields{
				"endpoint":      am.endpoint,
				"monitorType":   am.config.MonitorConfigCore().Type,
				"discoveryRule": am.config.MonitorConfigCore().DiscoveryRule,
			}).Debug("Shutting down doomed monitor")

			am.Shutdown()
		} else {
			newActiveMonitors = append(newActiveMonitors, am)
		}
	}

	mm.activeMonitors = newActiveMonitors
}

func (mm *MonitorManager) makeMonitorsForMatchingEndpoints(conf config.MonitorCustomConfig) {
	for id, endpoint := range mm.discoveredEndpoints {
		// Self configured endpoints are monitored immediately upon being
		// created and never need to be matched against discovery rules.
		if endpoint.Core().IsSelfConfigured() {
			continue
		}

		log.WithFields(log.Fields{
			"monitorType":   conf.MonitorConfigCore().Type,
			"discoveryRule": conf.MonitorConfigCore().DiscoveryRule,
			"endpoint":      endpoint,
		}).Debug("Trying to find config that matches discovered endpoint")

		if mm.isEndpointIDMonitoredByConfig(conf, id) {
			continue
		}

		if matched, err := mm.monitorEndpointIfRuleMatches(conf, endpoint); matched {
			if err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"endpointID":  endpoint.Core().ID,
					"monitorType": conf.MonitorConfigCore().Type,
				}).Error("Error monitoring endpoint that matched rule")
			} else {
				log.WithFields(log.Fields{
					"endpointID":  endpoint.Core().ID,
					"monitorType": conf.MonitorConfigCore().Type,
				}).Info("Now monitoring discovered endpoint")
			}
		}
	}
}

func (mm *MonitorManager) isEndpointIDMonitoredByConfig(conf config.MonitorCustomConfig, id services.ID) bool {
	for _, am := range mm.activeMonitors {
		if conf.MonitorConfigCore().Equals(am.config.MonitorConfigCore()) {
			return true
		}
	}
	return false
}

// Returns true is the service did match a rule in this monitor config
func (mm *MonitorManager) monitorEndpointIfRuleMatches(config config.MonitorCustomConfig, endpoint services.Endpoint) (bool, error) {
	if config.MonitorConfigCore().DiscoveryRule == "" || !services.DoesServiceMatchRule(endpoint, config.MonitorConfigCore().DiscoveryRule) {
		return false, nil
	}

	err := mm.createAndConfigureNewMonitor(config, endpoint)
	if err != nil {
		return false, err
	}

	return true, nil
}

// EndpointAdded should be called when a new service is discovered
func (mm *MonitorManager) EndpointAdded(endpoint services.Endpoint) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	ensureProxyingDisabledForService(endpoint)
	mm.discoveredEndpoints[endpoint.Core().ID] = endpoint

	// If the endpoint has a monitor type specified, then it is expected to
	// have all of its configuration already set in the endpoint and discovery
	// rules will be ignored.
	if endpoint.Core().IsSelfConfigured() {
		mm.monitorSelfConfiguredEndpoint(endpoint)
		return
	}

	mm.findConfigForMonitorAndRun(endpoint)
}

func (mm *MonitorManager) monitorSelfConfiguredEndpoint(endpoint services.Endpoint) {
	monitorType := endpoint.Core().MonitorType
	conf := &config.MonitorConfig{
		Type: monitorType,
	}

	monConfig, err := getCustomConfigForMonitor(conf)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"monitorType": monitorType,
		}).Error("Could not create monitor config for self-configured endpoint")
	}

	err = mm.createAndConfigureNewMonitor(monConfig, endpoint)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"endpoint": endpoint,
		}).Error("Could not create monitor for self-configured endpoint")
	}
}

func (mm *MonitorManager) findConfigForMonitorAndRun(endpoint services.Endpoint) {
	monitoring := false

	for _, config := range mm.monitorConfigs {
		matched, err := mm.monitorEndpointIfRuleMatches(config, endpoint)
		monitoring = matched || monitoring
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"config":   config,
				"endpoint": endpoint,
			}).Error("Could not monitor new endpoint")
		}
	}

	if !monitoring {
		log.WithFields(log.Fields{
			"endpoint": endpoint,
		}).Debug("Endpoint added that doesn't match any discovery rules")
	}
}

// endpoint may be nil for static monitors
func (mm *MonitorManager) createAndConfigureNewMonitor(config config.MonitorCustomConfig, endpoint services.Endpoint) error {
	id := types.MonitorID(mm.idGenerator())

	instance := newMonitor(config.MonitorConfigCore().Type, id)
	if instance == nil {
		return errors.Errorf("Could not create new monitor of type %s", config.MonitorConfigCore().Type)
	}

	am := &ActiveMonitor{
		id:        id,
		instance:  instance,
		endpoint:  endpoint,
		agentMeta: mm.agentMeta,
	}

	am.injectDatapointChannelIfNeeded(mm.dpChan)
	am.injectEventChannelIfNeeded(mm.eventChan)
	am.injectDimPropertiesChannelIfNeeded(mm.dimPropChan)

	if err := am.configureMonitor(config); err != nil {
		return err
	}
	mm.activeMonitors = append(mm.activeMonitors, am)

	log.WithFields(log.Fields{
		"monitorType":   config.MonitorConfigCore().Type,
		"discoveryRule": config.MonitorConfigCore().DiscoveryRule,
	}).Debug("Creating new monitor")

	return nil
}

func (mm *MonitorManager) monitorsForEndpointID(id services.ID) (out []*ActiveMonitor) {
	for i := range mm.activeMonitors {
		if mm.activeMonitors[i].endpointID() == id {
			out = append(out, mm.activeMonitors[i])
		}
	}
	return // Named return value
}

func (mm *MonitorManager) isServiceMonitored(id services.ID) bool {
	return len(mm.monitorsForEndpointID(id)) > 0
}

// EndpointRemoved should be called by observers when a service endpoint was
// removed.
func (mm *MonitorManager) EndpointRemoved(endpoint services.Endpoint) {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	delete(mm.discoveredEndpoints, endpoint.Core().ID)

	monitors := mm.monitorsForEndpointID(endpoint.Core().ID)
	for _, am := range monitors {
		am.doomed = true
	}
	mm.deleteDoomedMonitors()

	log.WithFields(log.Fields{
		"endpoint": endpoint,
	}).Info("No longer monitoring service")
}

func (mm *MonitorManager) isEndpointMonitored(endpoint services.Endpoint) bool {
	monitors := mm.monitorsForEndpointID(endpoint.Core().ID)
	return len(monitors) > 0
}

// Shutdown will shutdown all managed monitors and deinitialize the manager.
func (mm *MonitorManager) Shutdown() {
	mm.lock.Lock()
	defer mm.lock.Unlock()

	for i := range mm.activeMonitors {
		mm.activeMonitors[i].doomed = true
	}
	mm.deleteDoomedMonitors()

	mm.activeMonitors = nil
	mm.discoveredEndpoints = nil
}
