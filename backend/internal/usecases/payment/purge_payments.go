package commands

import (
	"context"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
)

type PurgePaymentsCommand struct {
	repo payment.Repository
}

func NewPurgePaymentsCommand(repo payment.Repository) *PurgePaymentsCommand {
	return &PurgePaymentsCommand{
		repo: repo,
	}
}

func (c *PurgePaymentsCommand) Execute(ctx context.Context) error {
	return c.repo.PurgePayments(ctx)
}
