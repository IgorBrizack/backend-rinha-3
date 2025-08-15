package commands

import (
	"context"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
)

type PaymentSummaryParams struct {
	From *time.Time `form:"from" json:"from" query:"from"`
	To   *time.Time `form:"to" json:"to" query:"to"`
}

type GetPaymentSummaryCommand struct {
	paymentRepository payment.Repository
}

func NewGetPaymentSummaryCommand(
	paymentRepository payment.Repository,
) *GetPaymentSummaryCommand {
	return &GetPaymentSummaryCommand{
		paymentRepository: paymentRepository,
	}
}

func (c *GetPaymentSummaryCommand) Execute(ctx context.Context, params PaymentSummaryParams) (dto.PaymentSummaryResponse, error) {
	var summary dto.PaymentSummaryResponse

	payments, err := c.paymentRepository.GetPayments(ctx, params.From, params.To)
	if err != nil {
		return summary, err
	}

	return payments, nil
}
