package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/redis/go-redis/v9"
)

const cacheKey = "payment_summary"

func UpdatePaymentSummary(ctx context.Context, client *redis.Client, payment dto.PaymentRequest, processorType string) error {
	var summary dto.PaymentSummaryResponse

	existing, err := client.Get(ctx, cacheKey).Result()
	if err != nil && err != redis.Nil {
		return fmt.Errorf("erro ao obter resumo do cache: %w", err)
	}

	if existing != "" {
		if err := json.Unmarshal([]byte(existing), &summary); err != nil {
			return fmt.Errorf("erro ao desserializar resumo existente: %w", err)
		}
	}

	switch processorType {
	case "default":
		summary.Default.TotalRequests++
		summary.Default.TotalAmount += payment.Amount
	case "fallback":
		summary.Fallback.TotalRequests++
		summary.Fallback.TotalAmount += payment.Amount
	default:
		return fmt.Errorf("tipo de processador inválido: %s", processorType)
	}

	updatedJSON, err := json.Marshal(summary)
	if err != nil {
		return fmt.Errorf("erro ao serializar resumo atualizado: %w", err)
	}

	if err := client.Set(ctx, cacheKey, updatedJSON, 0).Err(); err != nil {
		return fmt.Errorf("erro ao salvar resumo no Redis: %w", err)
	}

	fmt.Printf("📊 Resumo atualizado [%s]: %+v\n", processorType, summary)
	return nil
}
