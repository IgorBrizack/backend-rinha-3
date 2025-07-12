package payment

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	ID            string          `gorm:"type:char(36);primaryKey;default:(UUID())"`
	CorrelationID string          `gorm:"type:char(36);uniqueIndex;not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(10,2);"`
	Default       bool            `gorm:"not null"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
}
