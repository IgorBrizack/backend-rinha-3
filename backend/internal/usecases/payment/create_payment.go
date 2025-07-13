package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	payment "github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"

	"github.com/redis/go-redis/v9"
)

type CreatePaymentCommand struct {
	cacheClient       *redis.Client
	paymentRepository payment.Repository
}

func NewCreatePaymentCommand(
	cacheClient *redis.Client,
	paymentRepository payment.Repository,
) *CreatePaymentCommand {
	return &CreatePaymentCommand{
		cacheClient:       cacheClient,
		paymentRepository: paymentRepository,
	}
}

func (c *CreatePaymentCommand) Execute(ctx context.Context, payment dto.PaymentRequest) error {

	timeUTC := time.Now().UTC()
	port := os.Getenv("BACKEND_PORT")

	if err := c.sandToWorkerQueue(payment, timeUTC, "pending_payments"+port); err != nil {
		fmt.Print("Erro ao enviar para worker")
		return err
	}

	return nil
}

func (c *CreatePaymentCommand) sandToWorkerQueue(paymentRequestData dto.PaymentRequest, timeUTC time.Time, queue string) error {

	payload, err := json.Marshal(dto.PaymentRequestService{
		CorrelationID: paymentRequestData.CorrelationID,
		Amount:        paymentRequestData.Amount,
		RequestedAt:   timeUTC,
	})
	if err != nil {
		return fmt.Errorf("erro ao serializar pagamento: %w", err)
	}

	return c.cacheClient.RPush(context.Background(), queue, payload).Err()
}
