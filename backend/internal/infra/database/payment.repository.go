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

func (r *paymentRepository) CreatePayment(ctx context.Context, p payment.Payment) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	key := "payment:" + p.ID

	if err := r.redisClient.Set(ctx, key, data, 0).Err(); err != nil {
		return err
	}

	score := float64(p.CreatedAt.Unix())
	if err := r.redisClient.ZAdd(ctx, "payments_index", redis.Z{
		Score:  score,
		Member: key,
	}).Err(); err != nil {
		return err
	}

	return nil
}

func (r *paymentRepository) GetPayments(ctx context.Context, from, to *time.Time) ([]payment.Payment, error) {
	var start, end string
	if from != nil {
		start = fmt.Sprintf("%d", from.Unix())
	} else {
		start = "-inf"
	}
	if to != nil {
		end = fmt.Sprintf("%d", to.Unix())
	} else {
		end = "+inf"
	}

	ids, err := r.redisClient.ZRangeByScore(ctx, "payments_index", &redis.ZRangeBy{
		Min: start,
		Max: end,
	}).Result()
	if err != nil {
		return nil, err
	}

	var result []payment.Payment
	for _, key := range ids {
		data, err := r.redisClient.Get(ctx, key).Bytes()
		if err != nil {
			continue // skip not found
		}
		var p payment.Payment
		if err := json.Unmarshal(data, &p); err != nil {
			continue
		}
		result = append(result, p)
	}

	return result, nil
}
