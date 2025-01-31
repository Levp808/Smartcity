package main

import (
	"context"
	"depart/internal/config"
	"depart/internal/grpc/clients/petition"
	"depart/internal/grpc/server"
	"depart/internal/kafka"
	"depart/internal/repository"
	"depart/internal/service"
	"depart/pkg/postgres"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.BuildConfig()
	err := cleanenv.ReadConfig("", cfg)
	if err != nil {
		log.Fatal()
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()
	zerolog.SetGlobalLevel(cfg.LogLevel)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	dbPool, err := postgres.InitPostgresDbConnection(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	departmentRepo := repository.NewReportPgRepository(dbPool)
	petition, err := petition.NewPetitionGrpcClient(ctx, cfg.PetitionGrpcServiceUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize petition client")
	}

	departmentSevice := service.NewDepartService(departmentRepo, petition)
	grpcServer := server.NewGrpcServer(cfg.GrpcServerAddress, departmentSevice)

	kafkaConfig := kafka.Config{
		Brokers:         cfg.KafkaBrokers,
		TopicNames:      cfg.TopicNames,
		ConsumerGroupId: cfg.ConsumerGroupId,
	}

	kafkaEventListener, err := kafka.NewPetitionListener(ctx, kafkaConfig, departmentSevice)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize kafka listener")
		panic(err)
	}

	go grpcServer.Run(ctx)
	go kafkaEventListener.Run(ctx)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	<-graceCh
	grpcServer.Stop()
	dbPool.Close()
}
