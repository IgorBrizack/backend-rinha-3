package controller

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	commands "github.com/IgorBrizack/backend-rinha-3/internal/usecases/payment"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PaymentController struct {
	cacheClient       *redis.Client
	paymentRepository payment.Repository
}

func NewPaymentController(
	cacheClient *redis.Client,
	paymentRepository payment.Repository,
) *PaymentController {
	return &PaymentController{
		cacheClient:       cacheClient,
		paymentRepository: paymentRepository,
	}
}

func (pc *PaymentController) CreatePayment(c *gin.Context) {
	var paymentRequest dto.PaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	err := commands.NewCreatePaymentCommand(pc.cacheClient).Execute(paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(200, gin.H{"message": "Payment created successfully"})
}

func (pc *PaymentController) GetPaymentSummary(c *gin.Context) {
	var params commands.PaymentSummaryParams

	if err := c.BindQuery(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid or missing query parameters"})
		return
	}

	cmd := commands.NewGetPaymentSummaryCommand(pc.cacheClient, pc.paymentRepository)

	summary, err := cmd.Execute(params)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve payment summary"})
		return
	}

	c.JSON(200, summary)
}
