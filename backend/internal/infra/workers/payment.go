package workers

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

func PaymentWorker(client *redis.Client, paymentService *services.PaymentService, paymentRepository payment.Repository, paymentQueue chan []byte) {
	const numWorkers = 20
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			ctx := context.Background()

			for {
				payload, ok := <-paymentQueue
				if !ok {
					break // Canal fechado, encerra worker
				}

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
					_ = paymentRepository.CreatePayment(ctx, entity)
					continue
				}

				if fallbackHealth.Failing {
					continue
				} else {
					if err := paymentService.CreatePaymentDefault(req); err != nil {
						continue
					}
					entity.Default = true
					_ = paymentRepository.CreatePayment(ctx, entity)
					continue
				}

			}
		}(i)
	}
}
