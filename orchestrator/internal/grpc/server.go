package grpc

import (
	"context"
	"fmt"
	"net"
	"orchestrator/internal/service"
	"orchestrator/pkg/grpc/pb"
	"orchestrator/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

func New(ctx context.Context, port int, srv *service.ModerationService) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, fmt.Sprintf("Failed to listen: %v", err))
	}
	opts := []grpc.ServerOption{}

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterOrchestratorServiceServer(grpcServer, NewOrchestratorService(ctx, srv))

	return &Server{grpcServer: grpcServer, listener: lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	logger.GetLoggerFromContext(ctx).Info(ctx, "starting gRPC server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
	return s.grpcServer.Serve(s.listener)
}

func (s *Server) Stop(ctx context.Context) {
	logger.GetLoggerFromContext(ctx).Info(ctx, "stopping gRPC server")
	s.grpcServer.GracefulStop()
}
