package destination

import (
	"errors"
	"math/rand"

	"github.com/koladilip/event-server/event"
	"github.com/koladilip/event-server/utils"
	"go.uber.org/zap"
)

type Destination interface {
	Id() string
	Transform(event.SourceEvent) (event.DestinationEvent, error)
	Deliver(event.DestinationEvent) error
}

type Destination1 struct {
	id string
	// config goes here
	logger *zap.Logger
}

func (d *Destination1) Id() string {
	return d.id
}

func (d *Destination1) Transform(e event.SourceEvent) (event.DestinationEvent, error) {
	return event.DestinationEvent{
		DestinationID: "Destination1",
		UserID:        e.UserID,
		Payload:       "Destination1: " + e.Payload,
	}, nil
}

func (d *Destination1) Deliver(e event.DestinationEvent) error {
	utils.WaitForRandomPeriod()
	// Hold the logic to push the destination using supported protocol
	// of the destination
	if rand.Intn(100) > 70 {
		return errors.New("some error")
	} else {
		d.logger.Info("Delivered to destination1", zap.Any("event", e))
		return nil
	}
}

type Destination2 struct {
	id string
	// config goes here
	logger *zap.Logger
}

func (d *Destination2) Id() string {
	return d.id
}
func (d *Destination2) Transform(e event.SourceEvent) (event.DestinationEvent, error) {
	return event.DestinationEvent{
		DestinationID: "Destination2",
		UserID:        e.UserID,
		Payload:       "Destination2: " + e.Payload,
	}, nil
}

func (d *Destination2) Deliver(e event.DestinationEvent) error {
	utils.WaitForRandomPeriod()
	// Hold the logic to push the destination using supported protocol
	// of the destination
	if rand.Intn(100) > 30 {
		return errors.New("some error")
	} else {
		d.logger.Info("Delivered to destination2", zap.Any("event", e))
		return nil
	}
}
func NewDestinations(logger *zap.Logger) *map[string]Destination {
	return &map[string]Destination{
		"Destination1": &Destination1{id: "Destination1", logger: logger},
		"Destination2": &Destination2{id: "Destination2", logger: logger},
	}
}
