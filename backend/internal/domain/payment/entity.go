package payment

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID            string         `gorm:"type:char(36);primaryKey;default:(UUID())"` // char(36) é usado para UUID em MySQL
	CorrelationID string         `gorm:"type:char(36);uniqueIndex;not null"`
	Amount        float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"` // se quiser soft delete
}
