package destination

import (
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/koladilip/event-server/config"
	"github.com/koladilip/event-server/event"
	"github.com/koladilip/event-server/store"
	"github.com/koladilip/event-server/utils"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func sendToDestination(logger *zap.Logger, config *config.Config, baseCtx *config.BaseContext,
	reader *kafka.Reader, destination Destination) {
	for {
		m, err := reader.ReadMessage(baseCtx.Context)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(time.Second)
		}

		destEvent, err := event.NewDestinationEvent(m.Value)
		for {
			if err == nil {
				err = backoff.Retry(func() error {
					return destination.Deliver(baseCtx.Context, destEvent)
				}, backoff.NewExponentialBackOff())
			}
			if err != nil {
				logger.Error(err.Error())
				utils.WaitForRandomPeriod()
				continue
			}
			break
		}
	}
}

func StartDestinationSender(logger *zap.Logger, config *config.Config,
	baseCtx *config.BaseContext, destinations *map[string]Destination) {
	for id, destination := range *destinations {
		reader := store.NewReader(config, makeDestinationTopicName(id), fmt.Sprintf("%s-reader", id))
		// Simulating multiple systems
		go sendToDestination(logger, config, baseCtx, reader, destination)
	}
}
