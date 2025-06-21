package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"multi-tenant-messaging-app/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	if db == nil {
		panic("Database connection is nil") // Validasi koneksi database
	}
	return &MessageRepository{db: db}
}

func (r *MessageRepository) FetchMessages(tenantID string, cursor string, limit int) ([]map[string]interface{}, string, error) {
	var messages []model.Message
	query := r.db.Where("tenant_id = ?", tenantID).Order("id").Limit(limit)

	if cursor != "" {
		query = query.Where("id > ?", cursor)
	}

	if err := query.Find(&messages).Error; err != nil {
		return nil, "", fmt.Errorf("failed to fetch messages: %w", err)
	}

	var result []map[string]interface{}
	var nextCursor string

	for i, msg := range messages {
		var payloadMap map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &payloadMap); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		result = append(result, map[string]interface{}{
			"id":         msg.ID,
			"payload":    payloadMap,
			"created_at": msg.CreatedAt,
		})
		nextCursor = msg.ID.String()
		// Simpan next_cursor sebagai ID dari pesan terakhir jika ini adalah iterasi terakhir
		if i == len(messages)-1 {
			nextCursor = msg.ID.String()
		}
	}

	return result, nextCursor, nil
}

func (r *MessageRepository) SaveMessage(tenantID string, payload map[string]interface{}) error {
	log.Printf("Saving message for tenant %s: %v", tenantID, payload)

	// Konversi payload ke JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Hasilkan UUID untuk kolom id
	message := model.Message{
		ID:       uuid.New(), // Generate UUID
		TenantID: tenantID,
		Payload:  string(payloadJSON),
	}

	// Simpan pesan ke database
	if err := r.db.Create(&message).Error; err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	log.Printf("Message saved for tenant %s", tenantID)
	return nil
}
