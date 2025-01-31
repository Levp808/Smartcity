package service

import (
	"context"
	"orchestrator/ai_assist"
	"orchestrator/internal/models"
	"orchestrator/pkg/logger"

	"go.uber.org/zap"
)

type KafkaGateway interface {
	SendMessage(ctx context.Context, petition *models.Petition, topic string) error
}

type ModerationService struct {
	keyApi       string
	kafkaGateway KafkaGateway
}

func NewModerationService(key string, kafka KafkaGateway) *ModerationService {
	return &ModerationService{
		keyApi:       key,
		kafkaGateway: kafka,
	}
}

func (s *ModerationService) Moderation(ctx context.Context, petition *models.Petition) (string, error) {

	logger.GetLoggerFromContext(ctx).Info(ctx, "Request to AI assist: ", zap.String("description", petition.Description))
	depart, err := ai_assist.Connection(ctx, s.keyApi, petition.Description)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
		return "", err
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Response from AI assist: ", zap.String("response", depart))
	var topic string
	if depart == "Дорожная служба" {
		topic = "RoadDepart"
	}
	if depart == "Коммунальная служба" {
		topic = "CommDepart"
	}
	err = s.kafkaGateway.SendMessage(ctx, petition, topic)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Info(ctx, "Kafka gateway send message error")
		return "", err
	}
	return depart, err
}
