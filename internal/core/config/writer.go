package config

import (
	"net/url"

	"github.com/mitchellh/hashstructure"
	"github.com/signalfx/signalfx-agent/internal/core/dpfilters"
	log "github.com/sirupsen/logrus"
)

// WriterConfig holds configuration for the datapoint writer.
type WriterConfig struct {
	// These are soft limits and affect how much memory will be initially
	// allocated for datapoints, not the maximum memory allowed. Capacity
	// options get applied at startup and subsequent changes require an agent
	// restart.
	DatapointBufferCapacity uint `yaml:"datapointBufferCapacity" default:"1000"`
	// A hard limit on the number of buffered datapoints.  If this is hit, new
	// datapoints will be dropped until the buffer has room.
	DatapointBufferHardMax int `yaml:"datapointBufferHardMax" default:"100000"`
	// The anaologue of `datapointBufferCapacity` for events
	EventBufferCapacity uint `yaml:"eventBufferCapacity" default:"1000"`
	// A hard limit on the number of buffered events, beyond which all events
	// will be dropped until the buffer length drops below this.
	EventBufferHardMax int `yaml:"eventBufferHardMax" default:"10000"`
	// The agent does not send datapoints immediately upon a monitor generating
	// them, but buffers them and sends them in batches.  The lower this
	// number, the less delay for datapoints to appear in SignalFx.
	DatapointSendIntervalSeconds int `yaml:"datapointSendIntervalSeconds" default:"1"`
	// The analogue of `datapointSendIntervalSeconds` for events
	EventSendIntervalSeconds int `yaml:"eventSendIntervalSeconds" default:"1"`
	// How long to cache dimension objects in the SignalFx client.  Probably
	// won't ever need to be changed, but can be adjusted down if you have
	// extremely fast changing dimensions to reduce memory consumption.
	DimensionCacheResetIntervalSeconds int `yaml:"dimensionCacheResetIntervalSeconds" default:"900"`
	// If the log level is set to `debug` and this is true, all datapoints
	// generated by the agent will be logged.
	LogDatapoints bool `yaml:"logDatapoints"`
	// The analogue of `logDatapoints` for events.
	LogEvents bool `yaml:"logEvents"`
	// The following are propagated from elsewhere
	HostIDDims          map[string]string    `yaml:"-"`
	IngestURL           *url.URL             `yaml:"-"`
	APIURL              *url.URL             `yaml:"-"`
	SignalFxAccessToken string               `yaml:"-"`
	GlobalDimensions    map[string]string    `yaml:"-"`
	Filter              *dpfilters.FilterSet `yaml:"-"`
}

// Hash calculates a unique hash value for this config struct
func (wc *WriterConfig) Hash() uint64 {
	hash, err := hashstructure.Hash(wc, nil)
	if err != nil {
		log.WithError(err).Error("Could not get hash of WriterConfig struct")
		return 0
	}
	return hash
}
