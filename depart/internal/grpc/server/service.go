package server

import (
	"context"
	"database/sql"
	"depart/internal/entities"
	pb "depart/pkg/depart_v1"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Service interface {
	UpdateReport(ctx context.Context, petition_id uuid.UUID, report *entities.Report) (bool, error)
}

type DepartService struct {
	pb.UnimplementedDepartmentServiceServer
	service Service
}

func NewDepartService(service Service) *DepartService {
	return &DepartService{
		service: service,
	}
}

func (s *DepartService) UpdateReport(ctx context.Context, r *pb.UpdateReportRequest) (*pb.UpdateReportResponse, error) {
	log.Info().Msg("update report")
	petition_id, err := uuid.Parse(r.PetitionId)
	if err != nil {
		return &pb.UpdateReportResponse{Success: false, ErrorMessage: 2}, err
	}
	report := &entities.Report{}
	if r.ContentJob != nil {
		report.ContentJob = r.ContentJob.Value
	}
	if r.DoneAt != nil && !r.DoneAt.AsTime().IsZero() {
		report.DoneAt = sql.NullTime{
			Valid: true,
			Time:  r.DoneAt.AsTime(),
		}
	}
	success, err := s.service.UpdateReport(ctx, petition_id, report)
	if err != nil {
		return &pb.UpdateReportResponse{Success: false, ErrorMessage: 1}, err
	}
	if !success {
		return &pb.UpdateReportResponse{Success: false, ErrorMessage: 1}, nil
	}
	return &pb.UpdateReportResponse{Success: success}, nil
}
