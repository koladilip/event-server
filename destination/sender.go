package destination

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"koladilip.github.io/event-server/config"
	"koladilip.github.io/event-server/event"
	"koladilip.github.io/event-server/store"
)

func sendToDestination(logger *zap.Logger, config *config.Config, reader *kafka.Reader, destination Destination) {
	ctx := context.Background()
	for {
		if config.Shutdown {
			logger.Info("Stop sending messages to destination:" + destination.Id())
			break
		}
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(time.Second)
		}

		destEvent, err := event.NewDestinationEvent(m.Value)
		if err == nil {
			err = backoff.Retry(func() error {
				return destination.Deliver(destEvent)
			}, backoff.NewExponentialBackOff())
		}
		if err != nil {
			logger.Error(err.Error())
			//TODO handle this error gracefully
		}
	}
}

func StartDestinationSender(logger *zap.Logger, config *config.Config, destinations *map[string]Destination) {
	for id, destination := range *destinations {
		reader := store.NewReader(config, makeDestinationTopicName(id), fmt.Sprintf("%s-reader", id))
		// Simulating multiple systems
		go sendToDestination(logger, config, reader, destination)
	}
}
