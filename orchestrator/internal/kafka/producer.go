package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"orchestrator/internal/models"
	"orchestrator/internal/service"
	"orchestrator/pkg/logger"

	"github.com/segmentio/kafka-go"
)

type KafkaService struct {
	service.KafkaGateway
}

func (s *KafkaService) SendMessage(ctx context.Context, petition *models.Petition, topic string) error {

	brokers := []string{"localhost:9092"}

	w := &kafka.Writer{
		Addr:        kafka.TCP(brokers...),
		Topic:       topic,
		Balancer:    &kafka.RoundRobin{},
		Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	}
	defer w.Close()

	jsonData, err := json.Marshal(petition)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
	}

	msg := kafka.Message{
		Key:   []byte{},
		Value: jsonData,
	}

	err = w.WriteMessages(context.Background(), msg)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Error(ctx, err.Error())
	}
	logger.GetLoggerFromContext(ctx).Info(ctx, "Message sent to queue successfully")

	return nil
}

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}
