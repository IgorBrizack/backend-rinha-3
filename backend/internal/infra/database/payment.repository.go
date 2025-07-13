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

	if err := r.redisClient.LPush(ctx, "payment_queue", payload).Err(); err != nil {
		return fmt.Errorf("failed to push to redis list: %w", err)
	}

	if err := r.redisClient.ZAdd(ctx, "payment_index", redis.Z{
		Score:  float64(p.CreatedAt.Unix()),
		Member: payload,
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

	result, err := r.redisClient.ZRangeByScore(ctx, "payment_index", &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query redis sorted set: %w", err)
	}

	payments := make([]payment.Payment, 0, len(result))
	for _, item := range result {
		var p payment.Payment
		if err := json.Unmarshal([]byte(item), &p); err != nil {
			continue
		}
		payments = append(payments, p)
	}

	return payments, nil
}
