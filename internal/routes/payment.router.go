package routes

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/controller"
	redisInfra "github.com/IgorBrizack/backend-rinha-3/internal/infra/redis"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	paymentService := services.NewPaymentService()
	redisInfraClient := redisInfra.GetClient()
	paymentController := controller.NewPaymentController(redisInfraClient, paymentService)

	paymentGroup := r.Group("/payments")
	{
		paymentGroup.POST("/", paymentController.CreatePayment)
	}
}
