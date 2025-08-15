package payment

import (
	"context"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment Payment) error
	GetPayments(ctx context.Context, from, to *time.Time) (dto.PaymentSummaryResponse, error)
	PurgePayments(ctx context.Context) error
	CreatePaymentsBatch(ctx context.Context, payments []Payment) error
}
