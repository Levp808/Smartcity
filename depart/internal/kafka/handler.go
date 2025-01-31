package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type DepartmentService interface {
	CreateReport(ctx context.Context, petition_id uuid.UUID, location string, description_petition string) error
}

type petitionMsg struct {
	ID          string `json:"id"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

func (msg *petitionMsg) validate() error {
	if msg.ID == "" {
		log.Error().Msg("invalid message: petition_id field is required")
		return errors.New("invalid message: petition_id is required")
	}
	if msg.Location == "" {
		log.Error().Msg("invalid message: location field is required")
		return errors.New("invalid message: location is required")
	}
	if msg.Description == "" {
		log.Error().Msg("invalid message: description field is required")
		return errors.New("invalid message: description is required")
	}
	return nil
}

type PetitionHandler struct {
	DepartmentService DepartmentService
}

func NewPetitionHandler(departmentService DepartmentService) *PetitionHandler {
	return &PetitionHandler{DepartmentService: departmentService}
}

func (h *PetitionHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PetitionHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *PetitionHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			var data petitionMsg
			err := json.Unmarshal(message.Value, &data)
			if err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("failed unmarshalling message: %w", err)
			}
			if err := data.validate(); err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("received invalid message: %w", err)
			}
			dataID, err := uuid.Parse(data.ID)
			if err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("failed parsing uuid: %w", err)
			}
			err = h.DepartmentService.CreateReport(session.Context(), dataID, data.Location, data.Description)
			if err != nil {
				session.MarkMessage(message, "")
				return fmt.Errorf("error from department service: %w", err)
			}
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}
