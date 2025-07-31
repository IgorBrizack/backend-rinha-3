package commands

import (
	"context"

	payment "github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
)

type UpdatePaymentCommand struct {
	paymentRepository payment.Repository
}

func NewUpdatePaymentCommand(
	paymentRepository payment.Repository,
) *UpdatePaymentCommand {
	return &UpdatePaymentCommand{
		paymentRepository: paymentRepository,
	}
}

func (c *UpdatePaymentCommand) Execute(ctx context.Context, payment payment.Payment) error {
	return c.paymentRepository.UpdatePayment(ctx, payment)

}
