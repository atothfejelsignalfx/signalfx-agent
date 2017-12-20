package neopy

import (
	"encoding/json"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

const loggingSocketPath = "/tmp/signalfx-logs.sock"

// LogMessage represents the log message that comes back from python
type LogMessage struct {
	Message     string  `json:"message"`
	Level       string  `json:"level"`
	Logger      string  `json:"logger"`
	SourcePath  string  `json:"source_path"`
	LineNumber  string  `json:"lineno"`
	CreatedTime float64 `json:"created"`
}

// LoggingQueue wraps the zmq socket used to get log messages back from the
// python runner.
type LoggingQueue struct {
	listener net.Listener
	mutex    sync.Mutex
}

func (lq *LoggingQueue) start() error {
	var err error
	lq.listener, err = net.Listen("unixpacket", lq.socketPath())
	return err
}

func (lq *LoggingQueue) socketPath() string {
	return loggingSocketPath
}

func (lq *LoggingQueue) listenForLogMessages() {
	go func() {
		for {
			message, err := lq.listener.Accept()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Failed getting log message from NeoPy")
				continue
			}

			var msg LogMessage
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("Could not deserialize log message from NeoPy")
				continue
			}

			lq.handleLogMessage(&msg)
		}
	}()
}

func (lq *LoggingQueue) handleLogMessage(msg *LogMessage) {
	fields := log.Fields{
		"logger":      msg.Logger,
		"sourcePath":  msg.SourcePath,
		"lineno":      msg.LineNumber,
		"createdTime": msg.CreatedTime,
	}

	switch msg.Level {
	case "DEBUG":
		log.WithFields(fields).Debug(msg.Message)
	case "INFO":
		log.WithFields(fields).Info(msg.Message)
	case "WARNING":
		log.WithFields(fields).Warn(msg.Message)
	case "ERROR":
		log.WithFields(fields).Error(msg.Message)
	case "CRITICAL":
		// This will actually kill the agent, perhaps just log at error level
		// instead?
		log.WithFields(fields).Fatal(msg.Message)
	default:
		log.WithFields(fields).Errorf("No log level set for message: %s", msg.Message)
	}
}
