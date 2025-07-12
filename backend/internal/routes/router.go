package routes

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/controller"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	redisInfra "github.com/IgorBrizack/backend-rinha-3/internal/infra/redis"
	"github.com/gin-gonic/gin"
)

func SetupRouter(paymentRepository payment.Repository) *gin.Engine {
	router := gin.Default()

	redisInfraClient := redisInfra.GetClient()

	paymentController := controller.NewPaymentController(redisInfraClient, paymentRepository)

	router.POST("/payments", paymentController.CreatePayment)

	router.GET("/payments-summary", paymentController.GetPaymentSummary)

	return router
}
