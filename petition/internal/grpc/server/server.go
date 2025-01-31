package server

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"

	"petition/internal/entities"
	pb "petition/pkg/petition_v1"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PetitionService interface {
	Create(ctx context.Context, location string, description string) (string, error)
	Update(ctx context.Context, petition_id string, newPetition *entities.Petition) (bool, error)
	Delete(ctx context.Context, petition_id string) (bool, error)
	Get(ctx context.Context, petition_id string) (*entities.Petition, error)
}

type GrpcServer struct {
	pb.UnimplementedPetitionServer
	server               *grpc.Server
	petitionService      PetitionService
	grpcAdress           string
	gatewayServerAddress string
}

func NewGrpcServer(grpcAdress string, gatewayServerAddress string, petitionService PetitionService) *GrpcServer {
	server := grpc.NewServer()
	return &GrpcServer{
		server:               server,
		petitionService:      petitionService,
		grpcAdress:           grpcAdress,
		gatewayServerAddress: gatewayServerAddress,
	}
}

func (s *GrpcServer) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Info().Msg("Stopping grpc server")
		s.server.GracefulStop()
		log.Info().Msg("Stopped grpc server")
	}()

	pb.RegisterPetitionServer(s.server, s)

	listen, err := net.Listen("tcp", s.grpcAdress)
	if err != nil {
		log.Fatal().Err(err).Msg("listen error")
		return
	}

	log.Info().Msgf("Grpc Server listening on %s", s.grpcAdress)

	err = s.server.Serve(listen)
	if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		log.Error().Err(err).Msg("Grpc Server fail")
	}
}

func (s *GrpcServer) RunGateway(ctx context.Context) {
	conn, err := grpc.NewClient(
		s.gatewayServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to dial server")
		return
	}

	gwmux := runtime.NewServeMux()
	err = pb.RegisterPetitionHandler(ctx, gwmux, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to register gateway")
		return
	}

	gwServer := &http.Server{Addr: s.gatewayServerAddress, Handler: gwmux}

	go func() {
		<-ctx.Done()
		log.Info().Msg("shutting down grpc gateway server")
		if err := conn.Close(); err != nil {
			log.Error().Err(err).Msg("grpc gateway conn close: ")
		}
		if err := gwServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("grpc gateway server shutdown: ")
		}
		log.Info().Msg("grpc gateway server shut down")
	}()

	log.Info().Msgf("grpc gateway Server listening on %s", s.gatewayServerAddress)
	err = gwServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("grpc gateway server ListenAndServe:")
	}
}

func (s *GrpcServer) Stop() {
	log.Info().Msg("grpc server stopped")
	s.server.GracefulStop()
}

func (s *GrpcServer) CreatePetition(ctx context.Context, req *pb.CreatePetitionRequest) (*pb.CreatePetitionResponse, error) {
	log.Info().Msg("create petition")
	petitionId, err := s.petitionService.Create(ctx, req.Location, req.Description)
	if err != nil {
		return nil, err
	}
	log.Info().Msg("create petition")
	return &pb.CreatePetitionResponse{PetitionId: petitionId}, nil
}

func (s *GrpcServer) UpdatePetition(ctx context.Context, req *pb.UpdatePetitionRequest) (*pb.UpdatePetitionResponse, error) {
	log.Info().Msg("update petition")
	newPetition := &entities.Petition{}
	if req.Location != nil {
		newPetition.Location = req.Location.Value
	}
	if req.Description != nil {
		newPetition.Description = req.Description.Value
	}
	if req.Department != nil {
		newPetition.Department = req.Department.Value
	}
	if req.ReportId != nil {
		parsedUUID, err := uuid.Parse(req.ReportId.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse uuid")
			return &pb.UpdatePetitionResponse{Success: false, ErrorMessage: 2}, err
		}
		newPetition.ReportID = parsedUUID
	}
	if req.DoneAt != nil && !req.DoneAt.AsTime().IsZero() {
		newPetition.DoneAt = sql.NullTime{
			Valid: true,
			Time:  req.DoneAt.AsTime(),
		}
	}
	if req.ContentJob != nil {
		newPetition.ContentJob = req.ContentJob.Value
	}
	success, err := s.petitionService.Update(ctx, req.PetitionId, newPetition)
	if err != nil {
		log.Error().Err(err).Msg("failed to update")
		return &pb.UpdatePetitionResponse{Success: false, ErrorMessage: 0}, err
	}
	if !success {
		return &pb.UpdatePetitionResponse{Success: false, ErrorMessage: 1}, nil
	}
	return &pb.UpdatePetitionResponse{Success: true, ErrorMessage: 0}, nil
}

func (s *GrpcServer) DeletePetition(ctx context.Context, req *pb.DeletePetitionRequest) (*pb.DeletePetitionResponse, error) {
	log.Info().Msg("delete petition")
	success, err := s.petitionService.Delete(ctx, req.PetitionId)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete")
		return &pb.DeletePetitionResponse{Success: false, ErrorMessage: 0}, err
	}
	if !success {
		return &pb.DeletePetitionResponse{Success: false, ErrorMessage: 1}, nil
	}
	return &pb.DeletePetitionResponse{Success: true, ErrorMessage: 0}, nil
}

func (s *GrpcServer) GetPetition(ctx context.Context, req *pb.GetPetitionRequest) (*pb.GetPetitionResponse, error) {
	log.Info().Msg("get petition")
	petition, err := s.petitionService.Get(ctx, req.PetitionId)
	if err != nil {
		log.Error().Err(err).Msg("failed to get")
		return nil, err
	}
	if petition == nil {
		return nil, status.Errorf(codes.NotFound, "petition_id not found")
	}

	return &pb.GetPetitionResponse{
		PetitionId:  petition.PetitionID.String(),
		Location:    petition.Location,
		Description: petition.Description,
		CreatedAt:   timestamppb.New(petition.CreatedAt),
		Status:      petition.Status,
		Department:  wrapperspb.String(petition.Department),
		DoneAt:      timestamppb.New(petition.DoneAt.Time),
		ReportId:    wrapperspb.String(petition.ReportID.String()),
		ContentJob:  wrapperspb.String(petition.ContentJob),
	}, nil
}
