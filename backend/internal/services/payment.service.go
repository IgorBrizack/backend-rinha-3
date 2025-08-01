package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
)

type PaymentService struct {
	main_url     string
	fallback_url string
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		main_url:     "http://payment-processor-default:8080",
		fallback_url: "http://payment-processor-fallback:8080",
	}
}

func (s *PaymentService) CreatePaymentDefault(payment dto.PaymentRequestService) bool {
	return s.sendPayment(s.main_url, payment)
}

func (s *PaymentService) CreatePaymentFallback(payment dto.PaymentRequestService) bool {
	return s.sendPayment(s.fallback_url, payment)
}

func (s *PaymentService) sendPayment(url string, payment dto.PaymentRequestService) bool {
	payload, err := json.Marshal(payment)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", url+"/payments", bytes.NewBuffer(payload))
	if err != nil {
		return false
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (s *PaymentService) HealthCheckDefault() (dto.PaymentHealthCheckResponse, error) {
	return s.healthCheck(s.main_url)
}

func (s *PaymentService) HealthCheckFallback() (dto.PaymentHealthCheckResponse, error) {
	return s.healthCheck(s.fallback_url)
}

func (s *PaymentService) healthCheck(url string) (dto.PaymentHealthCheckResponse, error) {
	client := &http.Client{Timeout: 1 * time.Second}
	req, err := http.NewRequest("GET", url+"/payments/service-health", nil)
	if err != nil {
		return dto.PaymentHealthCheckResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return dto.PaymentHealthCheckResponse{Failing: true}, nil // Considera como falha
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.PaymentHealthCheckResponse{Failing: true}, nil
	}

	var result dto.PaymentHealthCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return dto.PaymentHealthCheckResponse{}, err
	}

	return result, nil
}
