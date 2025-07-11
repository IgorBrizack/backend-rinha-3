package commands

import (
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/redis/go-redis/v9"
)

type PaymentSummaryParams struct {
	From time.Time `form:"from" json:"from" query:"from"`
	To   time.Time `form:"to" json:"to" query:"to"`
}

type GetPaymentSummaryCommand struct {
	cacheClient       *redis.Client
	paymentRepository payment.Repository
}

func NewGetPaymentSummaryCommand(
	cacheClient *redis.Client,
	paymentRepository payment.Repository,
) *GetPaymentSummaryCommand {
	return &GetPaymentSummaryCommand{
		cacheClient:       cacheClient,
		paymentRepository: paymentRepository,
	}
}

func (c *GetPaymentSummaryCommand) Execute(params PaymentSummaryParams) (dto.PaymentSummaryResponse, error) {
	var summary dto.PaymentSummaryResponse

	return summary, nil
}
