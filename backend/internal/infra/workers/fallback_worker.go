package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

func StartFallbackWorker(paymentRepository payment.Repository, client *redis.Client, paymentService *services.PaymentService) {
	queueName := "fallback_queue"

	go func() {
		ctx := context.Background()

		for {
			result, err := client.BLPop(ctx, 0*time.Second, queueName).Result()
			if err != nil {
				fmt.Println("Erro lendo da fila fallback:", err)
				continue
			}

			if len(result) < 2 {
				continue
			}

			var req dto.PaymentRequestService
			if err := json.Unmarshal([]byte(result[1]), &req); err != nil {
				fmt.Println("Erro ao deserializar pagamento:", err)
				continue
			}

			// Tenta enviar para o serviço fallback
			if err := paymentService.CreatePaymentFallback(req); err != nil {
				fmt.Println("Erro ao processar pagamento fallback:", err)

				payload, errMarshal := json.Marshal(req)
				if errMarshal != nil {
					fmt.Println("Erro ao serializar pagamento para default:", errMarshal)
					continue
				}

				errPush := client.RPush(ctx, "default_queue", payload).Err()
				if errPush != nil {
					fmt.Println("Erro ao empurrar pagamento para default_queue:", errPush)
				} else {
					fmt.Println("Pagamento redirecionado para default_queue")
				}
			}

			entity := payment.Payment{
				CorrelationID: req.CorrelationID,
				Amount:        req.Amount,
				Default:       false,
				CreatedAt:     req.RequestedAt,
			}

			if err := paymentRepository.CreatePayment(entity); err != nil {
				fmt.Println("Erro ao criar pagamento no banco de dados:", err)
			}

		}
	}()
}
