package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

func StartFallbackWorker(paymentRepository payment.Repository, client *redis.Client, paymentService *services.PaymentService) {
	queueName := "fallback_queue"
	const numWorkers = 2

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

				var req dto.PaymentRequestService
				if err := json.Unmarshal([]byte(result[1]), &req); err != nil {
					continue
				}

				mainStatus, fallbackStatus := getHealthStatus(ctx, client)
				if fallbackStatus.Failing || (mainStatus.MinResponseTime > 0 && fallbackStatus.MinResponseTime > int(float64(mainStatus.MinResponseTime)*1.2)) {
					redirectToQueue(ctx, client, "default_queue", workerID, req)
					continue
				}

				if err := paymentService.CreatePaymentFallback(req); err != nil {
					redirectToQueue(ctx, client, "default_queue", workerID, req)
					continue
				}

				entity := payment.Payment{
					CorrelationID: req.CorrelationID,
					Amount:        req.Amount,
					Default:       false,
					CreatedAt:     req.RequestedAt.UTC(),
				}

				if err := paymentRepository.CreatePayment(entity); err != nil {
					fmt.Printf("[FallbackWorker %d] Erro ao salvar no banco: %v\n", workerID, err)
				}
			}
		}(i)
	}
}
