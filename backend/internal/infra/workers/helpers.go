package workers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/redis/go-redis/v9"
)

func getHealthStatus(ctx context.Context, client *redis.Client) (main, fallback dto.PaymentHealthCheckResponse) {
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

func redirectToFallback(ctx context.Context, client *redis.Client, workerID int, req dto.PaymentRequestService) {
	payload, err := json.Marshal(req)
	if err != nil {
		return
	}

	err = client.RPush(ctx, "fallback_queue", payload).Err()
	if err != nil {
		fmt.Printf("[Worker %d] Erro ao empurrar para fallback_queue: %v\n", workerID, err)
	} else {
		fmt.Printf("[Worker %d] Pagamento redirecionado para fallback_queue\n", workerID)
	}
}

func redirectToQueue(ctx context.Context, client *redis.Client, queueName string, workerID int, req dto.PaymentRequestService) {
	payload, err := json.Marshal(req)
	if err != nil {
		return
	}

	err = client.RPush(ctx, queueName, payload).Err()
	if err != nil {
		fmt.Printf("[Worker %d] Erro ao empurrar para %s: %v\n", workerID, queueName, err)
	} else {
		fmt.Printf("[Worker %d] Pagamento redirecionado para %s\n", workerID, queueName)
	}
}
