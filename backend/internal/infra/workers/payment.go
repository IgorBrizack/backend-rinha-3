package workers

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
)

type PaymentWorker struct {
	paymentService     *services.PaymentService
	paymentRepository  payment.Repository
	paymentQueue       chan []byte
	defaultQueue       chan []byte
	fallbackQueue      chan []byte
	paymentsBatchQueue chan payment.Payment
}

func NewPaymentWorker(
	paymentService *services.PaymentService,
	paymentRepository payment.Repository,
	paymentQueue chan []byte,
	defaultQueue chan []byte,
	fallbackQueue chan []byte,
) *PaymentWorker {
	return &PaymentWorker{
		paymentService:     paymentService,
		paymentRepository:  paymentRepository,
		paymentQueue:       paymentQueue,
		defaultQueue:       defaultQueue,
		fallbackQueue:      fallbackQueue,
		paymentsBatchQueue: make(chan payment.Payment, 1000),
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

			for payload := range w.paymentQueue {
				var req dto.PaymentRequestService
				if err := json.Unmarshal(payload, &req); err != nil {
					continue
				}

				w.defaultQueue <- payload

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
					} else {
						if isDefault {
							w.fallbackQueue <- payload
						} else {
							w.defaultQueue <- payload
						}
					}
				}
			}
		}()
	}
}
