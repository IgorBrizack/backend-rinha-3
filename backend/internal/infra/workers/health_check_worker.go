package workers

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/redis/go-redis/v9"
)

type HealthCheckerWorker struct {
	service     *services.PaymentService
	cacheClient *redis.Client
}

func NewHealthCheckerWorker(service *services.PaymentService, cacheClient *redis.Client) *HealthCheckerWorker {
	return &HealthCheckerWorker{
		service:     service,
		cacheClient: cacheClient,
	}
}

func (w *HealthCheckerWorker) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		w.checkHealth()
	}
}

func (w *HealthCheckerWorker) checkHealth() {
	ctx := context.Background()

	defaultResp, err := w.service.HealthCheckDefault()
	if err != nil {
		log.Printf("[ERROR] HealthCheckDefault failed: %v", err)
	} else {
		log.Printf("[INFO] Main Service Health: failing=%v, minResponseTime=%dms", defaultResp.Failing, defaultResp.MinResponseTime)
		w.saveToRedis(ctx, "health:main", defaultResp)
	}

	fallbackResp, err := w.service.HealthCheckFallback()
	if err != nil {
		log.Printf("[ERROR] HealthCheckFallback failed: %v", err)
	} else {
		log.Printf("[INFO] Fallback Service Health: failing=%v, minResponseTime=%dms", fallbackResp.Failing, fallbackResp.MinResponseTime)
		w.saveToRedis(ctx, "health:fallback", fallbackResp)
	}
}

func (w *HealthCheckerWorker) saveToRedis(ctx context.Context, key string, data dto.PaymentHealthCheckResponse) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}

	if err := w.cacheClient.Set(ctx, key, jsonData, 0).Err(); err != nil {
		log.Printf("[ERROR] Failed to save health status to Redis for key '%s': %v", key, err)
	} else {
		log.Printf("[REDIS] Health status saved for key '%s'", key)
	}
}
