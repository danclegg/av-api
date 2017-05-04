package base

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/byuoitav/event-router-microservice/eventinfrastructure"
	"github.com/xuther/go-message-router/common"
	"github.com/xuther/go-message-router/publisher"
)

func Publish(e eventinfrastructure.Event) error {
	var err error

	// create the event
	e.Timestamp = time.Now().Format(time.RFC3339)
	if len(os.Getenv("LOCAL_ENVIRONMENT")) > 0 {
		e.Hostname = os.Getenv("PI_HOSTNAME")
	} else {
		e.Hostname, err = os.Hostname()
	}
	if err != nil {
		return err
	}

	e.LocalEnvironment = len(os.Getenv("LOCAL_ENVIRONMENT")) > 0

	toSend, err := json.Marshal(&e)
	if err != nil {
		return err
	}

	// create the publisher
	p, err := publisher.NewPublisher("5000", 100, 10)
	if err != nil {
		return err
	}

	header := [24]byte{}
	copy(header[:], eventinfrastructure.LocalAPI)

	log.Printf("Publishing event: %+v", toSend)
	p.Write(common.Message{MessageHeader: header, MessageBody: toSend})

	return err
}
