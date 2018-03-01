package config

import (
	"reflect"

	"github.com/mitchellh/hashstructure"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	log "github.com/sirupsen/logrus"
)

// MonitorConfig is used to configure monitor instances.  One instance of
// MonitorConfig may be used to configure multiple monitor instances.  If a
// monitor's discovery rule does not match any discovered services, the monitor
// will not run.
type MonitorConfig struct {
	// The type of the monitor
	Type string `yaml:"type"`
	// The rule used to match up this configuration with a discovered endpoint.
	// If blank, the configuration will be run immediately when the agent is
	// started.  If multiple endpoints match this rule, multiple instances of
	// the monitor type will be created with the same configuration.
	DiscoveryRule string `yaml:"discoveryRule"`
	// A set of extra dimensions to include on datapoints emitted by the
	// monitor(s) created from this configuration
	ExtraDimensions map[string]string `yaml:"extraDimensions"`
	// The interval (in seconds) at which to emit datapoints from the
	// monitor(s) created by this configuration.  If not set (or set to 0), the
	// global agent intervalSeconds config option will be used instead.
	IntervalSeconds int `yaml:"intervalSeconds"`
	// If one or more configurations have this set to true, only those
	// configurations will be considered -- useful for testing
	Solo bool `yaml:"solo"`
	// OtherConfig is everything else that is custom to a particular monitor
	OtherConfig map[string]interface{} `yaml:",inline" neverLog:"omit"`
	// ValidationError is where a message concerning validation issues can go
	// so that diagnostics can output it.
	Hostname        string          `yaml:"-"`
	BundleDir       string          `yaml:"-"`
	ValidationError string          `yaml:"-" hash:"ignore"`
	MonitorID       types.MonitorID `yaml:"-" hash:"ignore"`
}

// Equals tests if two monitor configs are sufficiently equal to each other.
// Two monitors should only be equal if it doesn't make sense for two
// configurations to be active at the same time.
func (mc *MonitorConfig) Equals(other *MonitorConfig) bool {
	return mc.Type == other.Type && mc.DiscoveryRule == other.DiscoveryRule &&
		reflect.DeepEqual(mc.OtherConfig, other.OtherConfig)
}

// ExtraConfig returns generic config as a map
func (mc *MonitorConfig) ExtraConfig() map[string]interface{} {
	return mc.OtherConfig
}

// HasAutoDiscovery returns whether the monitor is static (i.e. doesn't rely on
// autodiscovered services and is manually configured) or dynamic.
func (mc *MonitorConfig) HasAutoDiscovery() bool {
	return mc.DiscoveryRule != ""
}

// MonitorConfigCore provides a way of getting the MonitorConfig when embedded
// in a struct that is referenced through a more generic interface.
func (mc *MonitorConfig) MonitorConfigCore() *MonitorConfig {
	return mc
}

// Hash calculates a unique hash value for this config struct
func (mc *MonitorConfig) Hash() uint64 {
	hash, err := hashstructure.Hash(mc, nil)
	if err != nil {
		log.WithError(err).Error("Could not get hash of MonitorConfig struct")
		return 0
	}
	return hash
}

// MonitorCustomConfig represents monitor-specific configuration that doesn't
// appear in the MonitorConfig struct.
type MonitorCustomConfig interface {
	MonitorConfigCore() *MonitorConfig
}
