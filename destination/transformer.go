package destination

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/koladilip/event-server/event"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Transformer struct {
	logger       *zap.Logger
	writer       *kafka.Writer
	destinations *map[string]Destination
}

func makeDestinationTopicName(destinationId string) string {
	return fmt.Sprintf("destination.%s.events", destinationId)
}

func (transformer *Transformer) TransformAndStore(ctx context.Context,
	sourceEvent event.SourceEvent) error {
	messages := []kafka.Message{}
	for id, destination := range *transformer.destinations {
		destinationEvent, err := destination.Transform(sourceEvent)
		var destinationPayload []byte
		if err == nil {
			destinationPayload, err = json.Marshal(destinationEvent)
		}
		if err != nil {
			transformer.logger.Error("Error in transforming event",
				zap.Any("sourceEvent", sourceEvent), zap.String("Destination", id))
			return err
		}
		messages = append(messages, kafka.Message{
			Topic: makeDestinationTopicName(id),
			Key:   []byte(sourceEvent.UserID),
			Value: destinationPayload,
		})
	}

	return transformer.writer.WriteMessages(ctx,
		messages...,
	)
}

func NewTransformer(logger *zap.Logger, writer *kafka.Writer, destinations *map[string]Destination) *Transformer {
	return &Transformer{
		logger:       logger,
		writer:       writer,
		destinations: destinations,
	}
}

var Fx = fx.Options(
	fx.Provide(NewDestinations),
	fx.Provide(NewTransformer),
	fx.Invoke(StartDestinationSender),
)
