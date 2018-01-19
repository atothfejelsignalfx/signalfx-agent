package internalmetrics

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"time"

	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/neo-agent/core/config"
	"github.com/signalfx/neo-agent/core/meta"
	"github.com/signalfx/neo-agent/monitors"
	"github.com/signalfx/neo-agent/monitors/types"
	"github.com/signalfx/neo-agent/utils"
	log "github.com/sirupsen/logrus"
)

const (
	monitorType = "internal-metrics"
)

// Config for internal metric monitoring
type Config struct {
	config.MonitorConfig
}

// Monitor for collecting internal metrics from the unix socket that dumps
// them.
type Monitor struct {
	DPs       chan<- *datapoint.Datapoint
	AgentMeta *meta.AgentMeta
	stop      func()
}

func init() {
	monitors.Register(monitorType, func(id types.MonitorID) interface{} { return &Monitor{} }, &Config{})
}

// Configure and kick off internal metric collection
func (m *Monitor) Configure(conf *Config) error {
	m.Shutdown()

	m.stop = utils.RunOnInterval(func() {
		c, err := net.Dial("unix", m.AgentMeta.InternalMetricsSocketPath)
		if err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"monitorType": monitorType,
				"path":        m.AgentMeta.InternalMetricsSocketPath,
			}).Error("Could not connect to internal metric socket")
			return
		}

		c.SetReadDeadline(time.Now().Add(5 * time.Second))
		jsonIn, err := ioutil.ReadAll(c)
		c.Close()
		if err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"monitorType": monitorType,
				"path":        m.AgentMeta.InternalMetricsSocketPath,
			}).Error("Could not read metrics from internal metric socket")
			return
		}

		dps := make([]*datapoint.Datapoint, 0)
		err = json.Unmarshal(jsonIn, &dps)

		for _, dp := range dps {
			m.DPs <- dp
		}
	}, time.Duration(conf.IntervalSeconds)*time.Second)

	return nil
}

// Shutdown the internal metric collection
func (m *Monitor) Shutdown() {
	if m.stop != nil {
		m.stop()
	}
}
