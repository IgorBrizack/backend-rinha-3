package payment

import (
	"context"
	"time"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment Payment) error
	GetPayments(ctx context.Context, from, to *time.Time) ([]Payment, error)
	GetByCorrelationID(ctx context.Context, correlationID string) (Payment, error)
	UpdatePayment(ctx context.Context, payment Payment) error
}
