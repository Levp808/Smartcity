package config

import (
	"github.com/rs/zerolog"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	GrpcServerAddress        string        `default:"0.0.0.0:50048" split_words:"true"`
	GrpcGatewayServerAddress string        `default:"0.0.0.0:8085" split_words:"true"`
	PetitionGrpcServiceUrl   string        `default:"0.0.0.0:50049" split_words:"true"`
	DbUser                   string        `default:"postgres" split_words:"true"`
	DbPassword               string        `default:"password" split_words:"true"`
	DbHost                   string        `default:"0.0.0.0:5431" split_words:"true"`
	DbName                   string        `default:"report" split_words:"true"`
	LogLevel                 zerolog.Level `default:"1" split_words:"true"`
	KafkaBrokers             []string      `default:"localhost:9092" split_words:"true"`
	TopicNames               []string      `default:"CommDepart" split_words:"true"`
	ConsumerGroupId          string        `default:"CommDepart-group" split_words:"true"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("depart", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Debug().Any("App config", config).Msg("config")

	return &config
}
