package commands

import (
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
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

	payments, err := c.GetPayments(params.From, params.To)
	if err != nil {
		return summary, err
	}

	result := c.calculate(payments)

	return result, nil
}

func (c *GetPaymentSummaryCommand) GetPayments(from, to time.Time) ([]payment.Payment, error) {
	payments, err := c.paymentRepository.GetPayments(from, to)

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (c *GetPaymentSummaryCommand) calculate(payments []payment.Payment) dto.PaymentSummaryResponse {
	totalDefault := decimal.NewFromInt(0)
	totalDefaultReq := 0

	totalFallback := decimal.NewFromInt(0)
	totalFallbackReq := 0

	for _, p := range payments {
		if p.Default {
			totalDefault = totalDefault.Add(p.Amount)
			totalDefaultReq++
		}
		totalFallback = totalFallback.Add(p.Amount)
		totalFallbackReq++
	}

	defaultAmount, _ := totalDefault.Float64()
	fallbackAmount, _ := totalFallback.Float64()

	return dto.PaymentSummaryResponse{
		Default: dto.PaymentProcessorSummary{
			TotalAmount:   defaultAmount,
			TotalRequests: totalDefaultReq,
		},
		Fallback: dto.PaymentProcessorSummary{
			TotalAmount:   fallbackAmount,
			TotalRequests: totalFallbackReq,
		},
	}
}
