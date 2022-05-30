package config

import "context"

const DEFAULT_ENV = "local"

// KafkaConfig is a set of configuration params for kafka
type KafkaConfig struct {
	Endpoint string
}

// ServerConfig is a set of configuration params for API
type ServerConfig struct {
	Port string
}

// Config for API
type Config struct {
	Kafka    KafkaConfig
	Server   ServerConfig
	Shutdown bool
}

type BaseContext struct {
	Context context.Context
	Cancel  context.CancelFunc
}
