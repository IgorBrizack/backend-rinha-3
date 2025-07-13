package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	payment "github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/infra/workers"

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
	queue := c.selectQueueToSend(ctx)

	if err := c.saveInDB(ctx, payment, timeUTC, queue); err != nil {
		fmt.Print("Erro ao salvar no banco")
		return err
	}

	if err := c.sandToWorkerQueue(payment, timeUTC, queue); err != nil {
		fmt.Print("Erro ao enviar para worker")
		return err
	}

	return nil
}

func (c *CreatePaymentCommand) selectQueueToSend(ctx context.Context) string {
	qdefault := "default_queue"
	qfallback := "fallback_queue"

	mainHealth, fallbackHealth := workers.GetHealthStatus(ctx, c.cacheClient)

	if mainHealth.Failing && fallbackHealth.Failing {
		return qdefault
	}

	if mainHealth.Failing || mainHealth.MinResponseTime > int(float64(fallbackHealth.MinResponseTime)*2) {
		return qfallback
	}

	return qdefault
}

func (c *CreatePaymentCommand) saveInDB(ctx context.Context, paymentRequestData dto.PaymentRequest, timeUTC time.Time, queue string) error {

	var isDefault bool
	if queue == "default_queue" {
		isDefault = true
	}

	entity := payment.Payment{
		CorrelationID: paymentRequestData.CorrelationID,
		Amount:        paymentRequestData.Amount,
		Default:       isDefault,
		CreatedAt:     timeUTC,
	}

	if err := c.paymentRepository.CreatePayment(ctx, entity); err != nil {
		fmt.Printf("Erro ao salvar no banco: %v\n", err)
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
