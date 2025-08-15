package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	payment "github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
)

type CreatePaymentCommand struct {
	paymentRepository payment.Repository
	paymentQueue      chan []byte
}

func NewCreatePaymentCommand(
	paymentRepository payment.Repository,
	paymentQueue chan []byte,

) *CreatePaymentCommand {
	return &CreatePaymentCommand{
		paymentRepository: paymentRepository,
		paymentQueue:      paymentQueue,
	}
}

func (c *CreatePaymentCommand) Execute(ctx context.Context, payment dto.PaymentRequest) error {

	timeUTC := time.Now().UTC()

	go c.sandToWorkerQueue(payment, timeUTC)

	return nil
}

func (c *CreatePaymentCommand) sandToWorkerQueue(paymentRequestData dto.PaymentRequest, timeUTC time.Time) error {

	payload, err := json.Marshal(dto.PaymentRequestService{
		CorrelationID: paymentRequestData.CorrelationID,
		Amount:        paymentRequestData.Amount,
		RequestedAt:   timeUTC,
	})
	if err != nil {
		return fmt.Errorf("erro ao serializar pagamento: %w", err)
	}

	select {
	case c.paymentQueue <- payload:
		return nil
	default:
		return fmt.Errorf("fila cheia ou indisponível")
	}
}
