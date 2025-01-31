package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"petition/internal/config"
	"petition/internal/grpc/clients/orchestrator"
	"petition/internal/grpc/server"
	"petition/internal/repository"
	"petition/internal/service"

	"petition/pkg/postgres"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).With().Caller().Logger()

	config := config.BuildConfig()
	zerolog.SetGlobalLevel(config.LogLevel)

	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	dbPool, err := postgres.InitPostgresDbConnection(ctx, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot initialize connection to postgres")
		return
	}

	orchestrator, err := orchestrator.NewOrchestratorGrpcClient(ctx, config.OrchestratorGrpcServiceUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed init grpc orchestrator client")
		return
	}

	petitionRepo := repository.NewPetitionPgRepository(dbPool)
	petitionService := service.NewPetitionService(petitionRepo, orchestrator)
	grpcServer := server.NewGrpcServer(config.GrpcServerAddress, config.GrpcGatewayServerAddress, petitionService)
	grpcServer.Run(ctx)

	<-graceCh
	grpcServer.Stop()
	dbPool.Close()
}
