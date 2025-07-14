package payment

import (
	"time"

	"github.com/shopspring/decimal"
)

type Payment struct {
	CorrelationID string          `gorm:"type:char(36);uniqueIndex;not null"`
	Amount        decimal.Decimal `gorm:"type:decimal(10,2);"`
	Default       bool            `gorm:"not null"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
}
