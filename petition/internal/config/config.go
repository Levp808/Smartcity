package config

import (
	"github.com/rs/zerolog"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HttpServerAddress          string        `default:"0.0.0.0:8080" split_words:"true"`
	GrpcServerAddress          string        `default:"0.0.0.0:50049" split_words:"true"`
	GrpcGatewayServerAddress   string        `default:"0.0.0.0:8084" split_words:"true"`
	OrchestratorGrpcServiceUrl string        `default:"0.0.0.0:50051" split_words:"true"`
	DbUser                     string        `default:"postgres" split_words:"true"`
	DbPassword                 string        `default:"password" split_words:"true"`
	DbHost                     string        `default:"0.0.0.0:5432" split_words:"true"`
	DbName                     string        `default:"petition" split_words:"true"`
	LogLevel                   zerolog.Level `default:"1" split_words:"true"`
}

func BuildConfig() *Config {
	var config Config

	err := envconfig.Process("petition", &config)
	if err != nil {
		log.Panic().Err(err).Msg("cant build config")
	}

	log.Debug().Any("App config", config).Msg("config")

	return &config
}
