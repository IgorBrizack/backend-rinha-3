package controller

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

func (pc *PaymentController) CreatePayment(c *gin.Context) {
	var paymentRequest dto.PaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	err := pc.paymentService.CreatePaymentDefault(paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(200, gin.H{"message": "Payment created successfully"})
}
