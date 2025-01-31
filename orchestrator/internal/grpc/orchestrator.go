package grpc

import (
	"context"
	"orchestrator/internal/models"
	"orchestrator/internal/service"
	"orchestrator/pkg/grpc/pb"
	"orchestrator/pkg/logger"

	"go.uber.org/zap"
)

type PetitionService struct {
	service    *service.ModerationService // Сервисный слой
	serviceCtx context.Context
	pb.UnimplementedOrchestratorServiceServer
}

func NewOrchestratorService(ctx context.Context, service *service.ModerationService) *PetitionService {
	return &PetitionService{
		serviceCtx: ctx,
		service:    service,
	}
}

// Moderation - gRPC метод
func (s *PetitionService) Moderation(ctx context.Context, req *pb.ModerationRequest) (*pb.ModerationResponse, error) {

	petition := &models.Petition{
		PetitionID:  req.Id,
		Location:    req.Location,
		Description: req.Description,
	}
	logger.GetLoggerFromContext(s.serviceCtx).Info(s.serviceCtx, "GRPC Request: ", zap.Any("petition", *petition))

	// Вызов сервисного слоя
	service, err := s.service.Moderation(s.serviceCtx, petition)
	if err != nil {
		return nil, nil
	}

	// преобразование в pb.ModerationResponse
	resp := &pb.ModerationResponse{
		Depart: service,
	}
	logger.GetLoggerFromContext(s.serviceCtx).Info(s.serviceCtx, "GRPC Response: ", zap.Any("response", *resp))

	return resp, nil
}
