package workers

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

func PaymentWorker(client *redis.Client, paymentService *services.PaymentService, paymentRepository payment.Repository) {
	const numWorkers = 10

	const pendingQueueName = "pending_payments"
	port := os.Getenv("BACKEND_PORT")
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			ctx := context.Background()

			for {
				result, err := client.BLPop(ctx, 0*time.Second, pendingQueueName+port).Result()
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

				mainHealth, fallbackHealth := GetHealthStatus(ctx, client)

				entity := payment.Payment{
					CorrelationID: req.CorrelationID,
					Amount:        req.Amount,
					Default:       true,
					CreatedAt:     req.RequestedAt,
				}

				if mainHealth.Failing || float64(mainHealth.MinResponseTime) > float64(fallbackHealth.MinResponseTime)*1.2 {
					if err := paymentService.CreatePaymentFallback(req); err != nil {
						continue
					}
					entity.Default = false
				} else {
					if err := paymentService.CreatePaymentDefault(req); err != nil {
						continue
					}
					entity.Default = true
				}

				_ = paymentRepository.CreatePayment(ctx, entity)
			}
		}(i)
	}
}
