package main

import (
	"context"
	"orchestrator/internal/config"
	"orchestrator/internal/grpc"
	"orchestrator/internal/kafka"
	"orchestrator/internal/service"
	"orchestrator/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	serviceName = "orchestrator"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)

	var cfg config.ConfigData
	err := cleanenv.ReadConfig("configs/scheduler.env", &cfg)
	if err != nil {
		mainLogger.Error(ctx, "Failed to read config")
	}
	mainLogger.Info(ctx, "Configs read")

	srv := service.NewModerationService(cfg.AIapi, &kafka.KafkaService{})
	mainLogger.Info(ctx, "Service layer created successfully")

	grpcserver, err := grpc.New(ctx, cfg.GRPCServerPort, srv)
	if err != nil {
		mainLogger.Error(ctx, "Failed to create grpc server")
	}
	mainLogger.Info(ctx, "Service layer created successfully")

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcserver.Start(ctx); err != nil {
			mainLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	grpcserver.Stop(ctx)
	mainLogger.Info(ctx, "Server stopped")
}
