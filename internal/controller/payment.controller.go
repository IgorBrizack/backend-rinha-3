package controller

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	paymentservice "github.com/IgorBrizack/backend-rinha-3/internal/services"
	commands "github.com/IgorBrizack/backend-rinha-3/internal/usecases/payment"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PaymentController struct {
	cacheClient    *redis.Client
	paymentService *paymentservice.PaymentService
}

func NewPaymentController(
	cacheClient *redis.Client,
	paymentService *paymentservice.PaymentService,
) *PaymentController {
	return &PaymentController{
		cacheClient:    cacheClient,
		paymentService: paymentService,
	}
}

func (pc *PaymentController) CreatePayment(c *gin.Context) {
	var paymentRequest dto.PaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	err := commands.NewCreatePaymentCommand(pc.cacheClient, pc.paymentService).Execute(paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(200, gin.H{"message": "Payment created successfully"})
}
