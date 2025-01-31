package server

import (
	"context"
	pb "depart/pkg/depart_v1"
	"errors"

	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	grpcServer *grpc.Server
	grpcAdress string
}

func NewGrpcServer(grpcAdress string, departmentService Service) *GrpcServer {
	grpcServer := grpc.NewServer()
	pb.RegisterDepartmentServiceServer(grpcServer, NewDepartService(departmentService))

	return &GrpcServer{grpcServer: grpcServer, grpcAdress: grpcAdress}
}

func (s *GrpcServer) Run(ctx context.Context) {
	listen, err := net.Listen("tcp", s.grpcAdress)
	if err != nil {
		log.Fatal().Err(err).Msg("listen error")
		return
	}

	log.Info().Msgf("Grpc Server listening on %s", s.grpcAdress)
	err = s.grpcServer.Serve(listen)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error().Err(err).Msg("Grpc Server fail")
	}
}

func (s *GrpcServer) Stop() {
	log.Info().Msg("grpc server stopped")
	s.grpcServer.GracefulStop()
}
