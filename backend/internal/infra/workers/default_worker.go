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

func StartDefaultWorker(paymentRepository payment.Repository, client *redis.Client, paymentService *services.PaymentService) {
	queueName := "default_queue"
	const numWorkers = 3

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

				mainHealth, fallbackHealth := getHealthStatus(ctx, client)

				if mainHealth.Failing || mainHealth.MinResponseTime > int(float64(fallbackHealth.MinResponseTime)*1.2) {
					redirectToFallback(ctx, client, workerID, req)
					continue
				}

				if err := paymentService.CreatePaymentDefault(req); err != nil {
					redirectToFallback(ctx, client, workerID, req)
					continue
				}

				entity := payment.Payment{
					CorrelationID: req.CorrelationID,
					Amount:        req.Amount,
					Default:       true,
					CreatedAt:     time.Now().UTC(),
				}

				if err := paymentRepository.CreatePayment(entity); err != nil {
					fmt.Printf("[Worker %d] Erro ao salvar no banco: %v\n", workerID, err)
				}
			}
		}(i)
	}
}
