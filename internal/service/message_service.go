package service

import (
	"encoding/json"
	"fmt"
	"multi-tenant-messaging-app/internal/repository"

	"github.com/streadway/amqp"
)

//go:generate mockgen -source=message_service.go -destination=mocks/message_service_mock.go -package=mocks
type MessageServiceInterface interface {
	FetchMessages(tenantID string, cursor string, limit int) ([]map[string]interface{}, string, error)
	PublishToTenantQueue(tenantID string, payload map[string]interface{}) error
}

type MessageService struct {
	repo       *repository.MessageRepository
	rabbitConn *amqp.Connection
}

func NewMessageService(repo *repository.MessageRepository, rabbitConn *amqp.Connection) *MessageService {
	return &MessageService{repo: repo, rabbitConn: rabbitConn}
}

func (s *MessageService) FetchMessages(tenantID string, cursor string, limit int) ([]map[string]interface{}, string, error) {
	return s.repo.FetchMessages(tenantID, cursor, limit)
}

func (s *MessageService) PublishToTenantQueue(tenantID string, payload map[string]interface{}) error {
	// Konversi payload ke JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Buka channel RabbitMQ
	ch, err := s.rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open RabbitMQ channel: %w", err)
	}
	defer ch.Close()

	// Tentukan nama queue berdasarkan tenant ID
	queueName := fmt.Sprintf("tenant_%s_queue", tenantID)

	// Kirim pesan ke queue RabbitMQ
	err = ch.Publish(
		"",        // Exchange
		queueName, // Routing key (queue name)
		false,     // Mandatory
		false,     // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payloadJSON,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
