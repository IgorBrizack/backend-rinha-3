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
	const numWorkers = 5

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()
			ctx := context.Background()

			for {
				result, err := client.BLPop(ctx, 0*time.Second, queueName).Result()
				if err != nil {
					fmt.Printf("[FallbackWorker %d] Erro lendo da fila: %v\n", workerID, err)
					continue
				}

				if len(result) < 2 {
					continue
				}

				var req dto.PaymentRequestService
				if err := json.Unmarshal([]byte(result[1]), &req); err != nil {
					fmt.Printf("[FallbackWorker %d] Erro ao deserializar pagamento: %v\n", workerID, err)
					continue
				}

				if err := paymentService.CreatePaymentFallback(req); err != nil {
					fmt.Printf("[FallbackWorker %d] Erro ao processar pagamento fallback: %v\n", workerID, err)

					payload, errMarshal := json.Marshal(req)
					if errMarshal != nil {
						fmt.Printf("[FallbackWorker %d] Erro ao serializar para default_queue: %v\n", workerID, errMarshal)
						continue
					}

					errPush := client.RPush(ctx, "default_queue", payload).Err()
					if errPush != nil {
						fmt.Printf("[FallbackWorker %d] Erro ao reempurrar para default_queue: %v\n", workerID, errPush)
					} else {
						fmt.Printf("[FallbackWorker %d] Pagamento redirecionado para default_queue\n", workerID)
					}
					continue
				}

				entity := payment.Payment{
					CorrelationID: req.CorrelationID,
					Amount:        req.Amount,
					Default:       false,
					CreatedAt:     req.RequestedAt.UTC(), // garante que é UTC
				}

				if err := paymentRepository.CreatePayment(entity); err != nil {
					fmt.Printf("[FallbackWorker %d] Erro ao salvar no banco: %v\n", workerID, err)
				}
			}
		}(i)
	}
}
