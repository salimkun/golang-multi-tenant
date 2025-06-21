package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"multi-tenant-messaging-app/internal/model"
	"strings"

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
	var lastID string

	for _, msg := range messages {
		var payloadMap map[string]interface{}
		if err := json.Unmarshal([]byte(msg.Payload), &payloadMap); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal payload: %w", err)
		}

		result = append(result, map[string]interface{}{
			"id":         msg.ID,
			"payload":    payloadMap,
			"created_at": msg.CreatedAt,
		})
		lastID = msg.ID.String()
	}

	return result, lastID, nil
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

func (r *MessageRepository) CreateTenantPartition(tenantID string) error {
	// Ganti karakter '-' dengan '_'
	tableName := fmt.Sprintf("messages_tenant_%s", tenantID)
	tableName = strings.ReplaceAll(tableName, "-", "_")

	query := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s PARTITION OF messages
        FOR VALUES IN ('%s')
    `, tableName, tenantID)

	// Eksekusi query menggunakan GORM
	if err := r.db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to create tenant partition: %w", err)
	}

	log.Printf("Partition created for tenant %s", tenantID)
	return nil
}
