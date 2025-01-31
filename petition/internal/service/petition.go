package service

import (
	"context"
	"errors"
	"petition/internal/entities"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type PetitionRepo interface {
	Create(ctx context.Context, p *entities.Petition) (string, error)
	Update(ctx context.Context, petition_id string, p *entities.Petition) (bool, error)
	Get(ctx context.Context, petition_id string) (*entities.Petition, error)
	Delete(ctx context.Context, petition_id string) (bool, error)
}

type OrchestratorProvider interface {
	GetDepartment(ctx context.Context, petition_id string, location string, description string) (string, error)
}

type PetitionService struct {
	petitionRepo PetitionRepo
	orchestrator OrchestratorProvider
}

func NewPetitionService(petitionRepo PetitionRepo, orchestrator OrchestratorProvider) *PetitionService {
	return &PetitionService{
		petitionRepo: petitionRepo,
		orchestrator: orchestrator}
}

func (s *PetitionService) Create(ctx context.Context, location string, description string) (string, error) {
	if description == "" || location == "" {
		return "", errors.New("description or location is empty")
	}
	petition := &entities.Petition{
		Location:    location,
		Description: description}

	petition_id, err := s.petitionRepo.Create(ctx, petition)
	if err != nil {
		log.Error().Msg("failed to create petition")
		return "", err
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	department, err := s.orchestrator.GetDepartment(ctxWithTimeout, petition_id, location, description)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error().Msg("request to orchestrator timed out")
		} else {
			log.Error().Err(err).Msg("failed to get status and department")
		}
		s.petitionRepo.Delete(ctx, petition_id)
		return "", err
	}

	if department == "" {
		log.Error().Msg("failed to get department")
		s.petitionRepo.Delete(ctx, petition_id)
		return "", err
	}

	petition_uuid, _ := uuid.Parse(petition_id)

	petition = &entities.Petition{
		PetitionID: petition_uuid,
		Status:     "moderated",
		Department: department,
	}
	_, err = s.Update(ctx, petition_id, petition)
	if err != nil {
		log.Error().Msg("failed to update petition")
		return "", err
	}
	return petition_id, nil
}

func (s *PetitionService) Update(ctx context.Context, petition_id string, newPetition *entities.Petition) (bool, error) {
	success, err := s.petitionRepo.Update(ctx, petition_id, newPetition)
	if err != nil {
		return false, err
	}
	if !success {
		log.Error().Msg("petition_id not found")
		return false, nil
	}
	return success, nil
}

func (s *PetitionService) Delete(ctx context.Context, petition_id string) (bool, error) {
	success, err := s.petitionRepo.Delete(ctx, petition_id)
	if err != nil {
		return false, err
	}
	if !success {
		log.Error().Msg("petition_id not found")
		return false, nil
	}
	return success, nil
}

func (s *PetitionService) Get(ctx context.Context, petition_id string) (*entities.Petition, error) {
	petition, err := s.petitionRepo.Get(ctx, petition_id)
	if err != nil {
		return nil, err
	}
	if petition == nil {
		log.Error().Msg("petition_id not found")
		return nil, nil
	}
	return petition, nil
}
