package workers

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

type PaymentWorker struct {
	client            *redis.Client
	paymentService    *services.PaymentService
	paymentRepository payment.Repository
	paymentQueue      chan []byte
	defaultQueue      chan []byte
	fallbackQueue     chan []byte
}

func NewPaymentWorker(
	client *redis.Client,
	paymentService *services.PaymentService,
	paymentRepository payment.Repository,
	paymentQueue chan []byte,
	defaultQueue chan []byte,
	fallbackQueue chan []byte,
) *PaymentWorker {
	return &PaymentWorker{
		client:            client,
		paymentService:    paymentService,
		paymentRepository: paymentRepository,
		paymentQueue:      paymentQueue,
		defaultQueue:      defaultQueue,
		fallbackQueue:     fallbackQueue,
	}
}

func (w *PaymentWorker) Start() {
	go w.MainWorker()
	go w.startWorker(w.defaultQueue, true, w.paymentService.CreatePaymentDefault)
	go w.startWorker(w.fallbackQueue, false, w.paymentService.CreatePaymentFallback)
}

func (w *PaymentWorker) MainWorker() {
	const numWorkers = 5
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for payload := range w.paymentQueue {
				var req dto.PaymentRequestService
				if err := json.Unmarshal(payload, &req); err != nil {
					continue
				}

				mainHealth, fallbackHealth := GetHealthStatus(ctx, w.client)

				switch {
				case mainHealth.Failing:
					w.fallbackQueue <- payload
				case fallbackHealth.Failing:
					w.paymentQueue <- payload
				default:
					w.defaultQueue <- payload
				}
			}
		}()
	}
}

func (w *PaymentWorker) startWorker(
	queue chan []byte,
	isDefault bool,
	createFn func(dto.PaymentRequestService) bool,
) {
	const numWorkers = 5
	const maxRetries = 5
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for payload := range queue {
				var req dto.PaymentRequestService
				if err := json.Unmarshal(payload, &req); err != nil {
					continue
				}

				entity := payment.Payment{
					CorrelationID: req.CorrelationID,
					Amount:        req.Amount,
					Default:       isDefault,
					CreatedAt:     req.RequestedAt,
				}

				for i := 0; i < maxRetries; i++ {
					if createFn(req) {
						_ = w.paymentRepository.CreatePayment(ctx, entity)
						break
					}
					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
	}
}
