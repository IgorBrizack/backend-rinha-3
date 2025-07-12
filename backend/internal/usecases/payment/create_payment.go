package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"

	"github.com/redis/go-redis/v9"
)

type CreatePaymentCommand struct {
	cacheClient *redis.Client
}

func NewCreatePaymentCommand(
	cacheClient *redis.Client,
) *CreatePaymentCommand {
	return &CreatePaymentCommand{
		cacheClient: cacheClient,
	}
}

func (c *CreatePaymentCommand) Execute(payment dto.PaymentRequest) error {

	payload, err := json.Marshal(dto.PaymentRequestService{
		CorrelationID: payment.CorrelationID,
		Amount:        payment.Amount,
		RequestedAt:   time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("erro ao serializar pagamento: %w", err)
	}

	queueName := "default_queue"

	return c.cacheClient.RPush(context.Background(), queueName, payload).Err()
}
