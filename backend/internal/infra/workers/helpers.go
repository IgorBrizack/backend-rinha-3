package workers

import (
	"context"
	"encoding/json"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/redis/go-redis/v9"
)

func GetHealthStatus(ctx context.Context, client *redis.Client) (main, fallback dto.PaymentHealthCheckResponse) {
	mainData, err := client.Get(ctx, "health:main").Result()
	if err == nil {
		_ = json.Unmarshal([]byte(mainData), &main)
	}

	fallbackData, err := client.Get(ctx, "health:fallback").Result()
	if err == nil {
		_ = json.Unmarshal([]byte(fallbackData), &fallback)
	}

	return main, fallback
}
