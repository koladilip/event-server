package event

import (
	"encoding/json"
	"errors"
)

type SourceEvent struct {
	UserID string `json:"userId"`
	// we can event type for supporting multiple different types of events
	// EventType string `json:"eventType"`
	Payload string `json:"payload"`
}

func NewSourceEvent(eventJson []byte) (SourceEvent, error) {
	sourceEvent := SourceEvent{}
	err := json.Unmarshal(eventJson, &sourceEvent)
	if err != nil || !sourceEvent.IsValid() {
		err = errors.New("Error in paring the message")
	}
	return sourceEvent, err
}

func (e *SourceEvent) IsValid() bool {
	return e.UserID != "" && e.Payload != ""
}

type DestinationEvent struct {
	DestinationID string `json:"destinationId"`
	UserID        string `json:"userId"`
	// we can event type for supporting multiple different types of events
	// EventType string `json:"eventType"`
	Payload string `json:"payload"`
}

func NewDestinationEvent(eventJson []byte) (DestinationEvent, error) {
	destEvent := DestinationEvent{}
	err := json.Unmarshal(eventJson, &destEvent)
	if err != nil || !destEvent.IsValid() {
		err = errors.New("Error in paring the message")
	}
	return destEvent, err
}

func (e *DestinationEvent) IsValid() bool {
	return e.DestinationID != "" &&
		e.UserID != "" && e.Payload != ""
}
