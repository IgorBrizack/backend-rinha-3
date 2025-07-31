package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/redis/go-redis/v9"
)

type paymentRepository struct {
	redisClient *redis.Client
}

func NewPaymentRepository(redisClient *redis.Client) payment.Repository {
	return &paymentRepository{
		redisClient: redisClient,
	}
}

func (r *paymentRepository) ExistsByCorrelationID(ctx context.Context, correlationID string) (bool, error) {
	key := "correlation:" + correlationID
	exists, err := r.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *paymentRepository) CreatePayment(ctx context.Context, p payment.Payment) error {
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	// Armazenamento por chave única para acesso direto
	key := fmt.Sprintf("payment:%s", p.CorrelationID)
	if err := r.redisClient.Set(ctx, key, payload, 0).Err(); err != nil {
		return fmt.Errorf("failed to store payment: %w", err)
	}

	// Indexação para buscas por período
	if err := r.redisClient.ZAdd(ctx, "payment_index", redis.Z{
		Score:  float64(p.CreatedAt.Unix()),
		Member: key,
	}).Err(); err != nil {
		return fmt.Errorf("failed to index payment in sorted set: %w", err)
	}

	return nil
}

func (r *paymentRepository) GetPayments(ctx context.Context, from, to *time.Time) ([]payment.Payment, error) {
	var min, max string

	if from != nil {
		min = fmt.Sprintf("%d", from.Unix())
	} else {
		min = "-inf"
	}

	if to != nil {
		max = fmt.Sprintf("%d", to.Unix())
	} else {
		max = "+inf"
	}

	// Busca as chaves do Sorted Set por score (timestamp)
	keys, err := r.redisClient.ZRangeByScore(ctx, "payment_index", &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query redis sorted set: %w", err)
	}

	payments := make([]payment.Payment, 0, len(keys))
	for _, key := range keys {
		payload, err := r.redisClient.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				continue // chave foi removida ou nunca existiu
			}
			return nil, fmt.Errorf("failed to get payment by key %s: %w", key, err)
		}

		var p payment.Payment
		if err := json.Unmarshal([]byte(payload), &p); err != nil {
			continue
		}
		payments = append(payments, p)
	}

	return payments, nil
}

// Constrói a mesma chave que foi usada para salvar.
func (r *paymentRepository) GetByCorrelationID(ctx context.Context, correlationID string) (payment.Payment, error) {
	var p payment.Payment
	key := fmt.Sprintf("payment:%s", correlationID)

	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return p, nil // não encontrado
		}
		return p, fmt.Errorf("failed to get payment: %w", err)
	}

	if err := json.Unmarshal([]byte(val), &p); err != nil {
		return p, fmt.Errorf("failed to unmarshal payment: %w", err)
	}

	return p, nil
}

func (r *paymentRepository) UpdatePayment(ctx context.Context, p payment.Payment) error {
	// Constrói a chave com base no CorrelationID
	key := fmt.Sprintf("payment:%s", p.CorrelationID)

	// Serializa a struct de pagamento para JSON
	payload, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal payment: %w", err)
	}

	// Salva o valor no Redis (sobrescrevendo o existente)
	if err := r.redisClient.Set(ctx, key, payload, 0).Err(); err != nil {
		return fmt.Errorf("failed to update payment in redis: %w", err)
	}

	return nil
}
