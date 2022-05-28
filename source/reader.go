package source

import (
	"context"
	"time"

	"github.com/koladilip/event-server/config"
	"github.com/koladilip/event-server/destination"
	"github.com/koladilip/event-server/event"
	"github.com/koladilip/event-server/store"
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func readMessages(logger *zap.Logger, config *config.Config, reader *kafka.Reader, tranformer *destination.Transformer) {

	ctx := context.Background()
	for {
		if config.Shutdown {
			logger.Info("Stop reading messages from source")
			break
		}
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.Error(err.Error())
			time.Sleep(time.Second)
		}

		sourceEvent, err := event.NewSourceEvent(m.Value)
		if err == nil {
			err = tranformer.TransformAndStore(ctx, sourceEvent)
		}
		if err != nil {
			logger.Error(err.Error())
			//TODO handle this error gracefully
		}
	}
}

func StartSourceReaders(logger *zap.Logger, config *config.Config,
	transformer *destination.Transformer) {
	reader := store.NewReader(config, store.SourceEventTopic, "source-events-reader")
	// Simulating parallel processing
	go readMessages(logger, config, reader, transformer)
	go readMessages(logger, config, reader, transformer)
}

var Fx = fx.Options(fx.Invoke(StartSourceReaders))
