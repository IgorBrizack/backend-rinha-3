package routes

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/controller"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	redisInfra "github.com/IgorBrizack/backend-rinha-3/internal/infra/redis"
	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(paymentRespository payment.Repository, r *gin.Engine) {
	redisInfraClient := redisInfra.GetClient()

	paymentController := controller.NewPaymentController(redisInfraClient, paymentRespository)

	paymentGroup := r.Group("/payments")
	{
		paymentGroup.POST("/", paymentController.CreatePayment)
		paymentGroup.GET("/summary", paymentController.GetPaymentSummary)
	}
}
