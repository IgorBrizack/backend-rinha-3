package controller

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment/dto"
	commands "github.com/IgorBrizack/backend-rinha-3/internal/usecases/payment"
	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	paymentRepository payment.Repository
	paymentQueue      chan []byte
}

func NewPaymentController(
	paymentRepository payment.Repository,
	paymentQueue chan []byte,
) *PaymentController {
	return &PaymentController{
		paymentRepository: paymentRepository,
		paymentQueue:      paymentQueue,
	}
}

func (pc *PaymentController) CreatePayment(c *gin.Context) {
	ctx := c.Request.Context()
	var paymentRequest dto.PaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	err := commands.NewCreatePaymentCommand(pc.paymentRepository, pc.paymentQueue).Execute(ctx, paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(200, gin.H{"message": "Payment created successfully"})
}

func (pc *PaymentController) GetPaymentSummary(c *gin.Context) {
	ctx := c.Request.Context()
	var params commands.PaymentSummaryParams

	if err := c.BindQuery(&params); err != nil {
		c.JSON(400, gin.H{"error": "Invalid or missing query parameters"})
		return
	}

	cmd := commands.NewGetPaymentSummaryCommand(pc.paymentRepository)

	summary, err := cmd.Execute(ctx, params)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to retrieve payment summary"})
		return
	}

	c.JSON(200, summary)
}

func (pc *PaymentController) PurgePayments(c *gin.Context) {
	ctx := c.Request.Context()

	err := commands.NewPurgePaymentsCommand(pc.paymentRepository).Execute(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to purge payments"})
		return
	}

	c.JSON(200, gin.H{"message": "Payments purged successfully"})
}
