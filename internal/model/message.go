package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	TenantID  string    `gorm:"index;not null"`
	Payload   string    `gorm:"type:jsonb;not null"` // JSONB untuk PostgreSQL
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
