package neopy

import (
	"encoding/json"
	"net"
	"sync"

	"github.com/signalfx/golib/datapoint"
	"github.com/signalfx/neo-agent/monitors"
	log "github.com/sirupsen/logrus"
)

const datapointsSocketPath = "/tmp/signalfx-datapoints.sock"
const datapointsTopic = "datapoints"

// DatapointMessage represents the message sent by python with a datapoint
type DatapointMessage struct {
	MonitorID monitors.MonitorID `json:"monitor_id"`
	// Will be deserialized by the golib method by itself
	Datapoint *datapoint.Datapoint
}

// DatapointsQueue wraps the zmq socket used to get datapoints back from python
type DatapointsQueue struct {
	listener net.Listener
	mutex    sync.Mutex
	dpChan   chan *DatapointMessage
}

func (dq *DatapointsQueue) start() error {
	var err error
	dq.listener, err = net.Listen("unixpacket", dq.socketPath())
	if err != nil {
		return err
	}

	dq.dpChan = make(chan *DatapointMessage)
	return nil
}

func (dq *DatapointsQueue) SocketPath() string {
	return datapointsSocketPath
}

func (dq *DatapointsQueue) AcceptMessage(message string) {
	log.Debug("waiting for datapoints")
	message, err := dq.listener.Accept()
	log.Debug("got datapoint")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed getting datapoints from NeoPy")
		return
	}

	var msg struct {
		MonitorID monitors.MonitorID `json:"monitor_id"`
		// Will be deserialized by the golib method by itself
		Datapoint *json.RawMessage
	}
	err = json.Unmarshal(message, &msg)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"message": message,
		}).Error("Could not deserialize datapoint message from NeoPy")
		return
	}

	var dp datapoint.Datapoint
	err = dp.UnmarshalJSON(*msg.Datapoint)
	if err != nil {
		log.WithFields(log.Fields{
			"dpJSON": string(*msg.Datapoint),
			"error":  err,
		}).Error("Could not deserialize datapoint from NeoPy")
		return
	}

	log.WithFields(log.Fields{
		"msg": msg,
	}).Debug("Datapoint Received from NeoPy")

	ch <- &DatapointMessage{
		MonitorID: msg.MonitorID,
		Datapoint: &dp,
	}
}
