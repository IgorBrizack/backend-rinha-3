package dto

type PaymentProcessorSummary struct {
	TotalRequests int     `json:"totalRequests"`
	TotalAmount   float64 `json:"totalAmount"`
}

type PaymentSummaryResponse struct {
	Default  PaymentProcessorSummary `json:"default"`
	Fallback PaymentProcessorSummary `json:"fallback"`
}

type PaymentRequest struct {
	ID     string  `json:"id"`
	Amount float64 `json:"amount"`
}

type PaymentHealthCheckResponse struct {
	Failing         bool `json:"failing"`
	MinResponseTime int  `json:"minResponseTime"`
}
