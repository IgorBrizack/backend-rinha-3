package services

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (s *PaymentService) CreatePaymentDefault(payment dto.PaymentRequest) error {
	payload, err := json.Marshal(payment)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", s.main_url+"/payments", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("erro ao criar pagamento: status %d", resp.StatusCode)
	}

	return nil
}
