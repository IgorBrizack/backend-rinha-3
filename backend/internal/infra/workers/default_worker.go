package workers

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

func StartDefaultWorker(client *redis.Client, paymentService *services.PaymentService) {
	queueName := "default_queue"
	const numWorkers = 10

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			ctx := context.Background()

			for {
				result, err := client.BLPop(ctx, 0*time.Second, queueName).Result()
				if err != nil {
					continue
				}

				if len(result) < 2 {
					continue
				}

				payload := result[1]

				var req dto.PaymentRequestService
				if err := json.Unmarshal([]byte(payload), &req); err != nil {
					continue
				}

				if err := paymentService.CreatePaymentDefault(req); err != nil {
					_ = client.RPush(ctx, queueName, payload).Err()
					continue
				}
			}
		}(i)
	}
}
