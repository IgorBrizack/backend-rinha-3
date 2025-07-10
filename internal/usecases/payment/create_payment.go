package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	paymentservice "github.com/IgorBrizack/backend-rinha-3/internal/services"

	"github.com/redis/go-redis/v9"
)

type CreatePaymentCommand struct {
	cacheClient    *redis.Client
	paymentService *paymentservice.PaymentService
}

func NewCreatePaymentCommand(
	cacheClient *redis.Client,
	paymentService *paymentservice.PaymentService) *CreatePaymentCommand {
	return &CreatePaymentCommand{
		cacheClient:    cacheClient,
		paymentService: paymentService,
	}
}

func (c *CreatePaymentCommand) Execute(payment dto.PaymentRequest) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("erro ao serializar pagamento: %w", err)
	}

	queueName := "default_queue"

	return c.cacheClient.RPush(context.Background(), queueName, payload).Err()
}
