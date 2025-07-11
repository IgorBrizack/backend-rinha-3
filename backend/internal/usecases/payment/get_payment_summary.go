package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/redis/go-redis/v9"
)

type GetPaymentSummaryCommand struct {
	cacheClient *redis.Client
}

func NewGetPaymentSummaryCommand(
	cacheClient *redis.Client,
) *GetPaymentSummaryCommand {
	return &GetPaymentSummaryCommand{
		cacheClient: cacheClient,
	}
}

func (c *GetPaymentSummaryCommand) Execute() (dto.PaymentSummaryResponse, error) {
	var summary dto.PaymentSummaryResponse

	val, err := c.cacheClient.Get(context.Background(), "payment_summary").Result()
	if err == redis.Nil {
		return summary, nil
	} else if err != nil {
		return summary, fmt.Errorf("erro ao acessar cache: %w", err)
	}

	if err := json.Unmarshal([]byte(val), &summary); err != nil {
		return summary, fmt.Errorf("erro ao desserializar dados do cache: %w", err)
	}

	return summary, nil
}
