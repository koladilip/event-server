package store

import (
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"koladilip.github.io/event-server/config"
)

func NewReader(config *config.Config, topic string, groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{config.Kafka.Endpoint},
		Topic:   topic,
		GroupID: groupId,
	})
}

func NewWriter(config *config.Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Kafka.Endpoint),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

var Fx = fx.Options(fx.Provide(NewWriter))
