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

func StartDefaultWorker(paymentRepository payment.Repository, client *redis.Client, paymentService *services.PaymentService) {
	queueName := "default_queue"

	go func() {
		ctx := context.Background()

		for {
			result, err := client.BLPop(ctx, 0*time.Second, queueName).Result()
			if err != nil {
				fmt.Println("Erro lendo da fila default:", err)
				continue
			}

			if len(result) < 2 {
				continue
			}

			var req dto.PaymentRequest
			if err := json.Unmarshal([]byte(result[1]), &req); err != nil {
				fmt.Println("Erro ao deserializar pagamento:", err)
				continue
			}

			// Primeiro, tenta processar com o serviço externo
			if err := paymentService.CreatePaymentDefault(req); err != nil {
				fmt.Println("Erro ao processar pagamento default:", err)

				payload, errMarshal := json.Marshal(req)
				if errMarshal != nil {
					fmt.Println("Erro ao serializar pagamento para fallback:", errMarshal)
					continue
				}

				errPush := client.RPush(ctx, "fallback_queue", payload).Err()
				if errPush != nil {
					fmt.Println("Erro ao empurrar pagamento para fallback_queue:", errPush)
				} else {
					fmt.Println("Pagamento redirecionado para fallback_queue")
				}
				continue
			}

			entity := payment.Payment{
				CorrelationID: req.CorrelationID,
				Amount:        req.Amount,
				CreatedAt:     time.Now(),
			}

			if err := paymentRepository.CreatePayment(entity); err != nil {
				fmt.Println("Erro ao criar pagamento no banco de dados:", err)
			}
		}
	}()
}
