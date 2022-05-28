package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// GetEnvName gets the current environment name
func GetEnvName() string {
	return getFromEnv("ENV", DEFAULT_ENV)
}

// getFromEnv gets environment variable from OS with fallback
func getFromEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setViperDefaults() {
	viper.SetDefault("server.requestTimeout", "15s")
}

func configureViper(envName string) {
	viper.SetConfigName("config." + envName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./env/")
	viper.AddConfigPath("./")
	setViperDefaults()
}

func readConfig(logger *zap.Logger) *Config {
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Error reading config file", zap.Error(err))
	}
	var configuration *Config

	if err := viper.Unmarshal(&configuration); err != nil {
		logger.Fatal("Unable to decode into struct", zap.Error(err))
	}
	
	return configuration
}

// NewConfig initialize configuration
func NewConfig(logger *zap.Logger) *Config {
	envName := GetEnvName()
	logger.Info("Loading configuration for", zap.String("ENV", envName))
	configureViper(envName)
	return readConfig(logger)
}

// NewLogger create new zap logger
func NewLogger() *zap.Logger {
	envName := GetEnvName()
	var logger *zap.Logger
	var err error
	switch envName {
	case "prod", "staging":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewDevelopment()
	}

	if err != nil {
		log.Fatal("Logger failed to initialize", err)
	}
	return logger
}

var Fx = fx.Options(
	fx.Provide(NewLogger),
	fx.Provide(NewConfig),
)