package service

import (
	"context"
	"database/sql"
	"depart/internal/entities"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Реализация интерфейса ConsumerGroupHandler
type PetitionConsumer interface {
	StartConsuming(ctx context.Context) error
}

type Config struct {
	Brokers         []string
	TopicNames      []string
	ConsumerGroupId string
}

type DepartRepo interface {
	Create(ctx context.Context, petition_id uuid.UUID, location string, description_petition string) (uuid.UUID, error)
	Update(ctx context.Context, petition_id uuid.UUID, report *entities.Report) (bool, error)
}

type PetitionProvider interface {
	UpdatePetition(ctx context.Context, petition_id uuid.UUID, done_at sql.NullTime, report_id uuid.UUID, content_job string) (bool, error)
}

type DepartService struct {
	departRepo       DepartRepo
	petitionProvider PetitionProvider
}

func NewDepartService(departRepo DepartRepo, petitionProvider PetitionProvider) *DepartService {
	return &DepartService{
		departRepo:       departRepo,
		petitionProvider: petitionProvider,
	}
}

func (s *DepartService) CreateReport(ctx context.Context, petition_id uuid.UUID, location string, description_petition string) error {
	log.Info().Msg("create report")
	report_id, err := s.departRepo.Create(ctx, petition_id, location, description_petition)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	var done_at sql.NullTime
	var content_job string
	success, err := s.petitionProvider.UpdatePetition(ctx, petition_id, done_at, report_id, content_job)
	if err != nil {
		log.Error().Msg(err.Error())
		return err
	}
	if !success {
		log.Error().Msg("failed to update petition")
		return err
	}
	return nil
}

func (s *DepartService) UpdateReport(ctx context.Context, petition_id uuid.UUID, report *entities.Report) (bool, error) {
	log.Info().Msg("update report")
	ok, err := s.departRepo.Update(ctx, petition_id, report)
	if err != nil {
		log.Error().Msg(err.Error())
		return false, err
	}
	if ok {
		success, err := s.petitionProvider.UpdatePetition(ctx, petition_id, report.DoneAt, report.ReportID, report.ContentJob)
		if err != nil {
			log.Error().Msg(err.Error())
			return false, err
		}
		if !success {
			log.Error().Msg("failed to update petition")
			return false, err
		}
		return ok, nil
	}
	return false, nil
}
