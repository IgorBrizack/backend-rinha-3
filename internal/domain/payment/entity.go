package payment

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CorrelationID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Amount        float64   `gorm:"type:numeric(10,2);not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}
