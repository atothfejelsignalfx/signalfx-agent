package neopy

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func listenForMessages(path string, cb func(string)) error {
	listener, err := net.Listen("unixpacket", path)
	if err != nil {
		return err
	}

	go func() {
		for {
			message, err := listener.Accept()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"path":  path,
				}).Error("Failed getting message from NeoPy")
				continue
			}
			cb(message)
		}
	}()
	return nil
}
