package routes

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/controller"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(r *gin.Engine) {
	paymentService := services.NewPaymentService()
	paymentController := controller.NewPaymentController(paymentService)

	paymentGroup := r.Group("/payments")
	{
		paymentGroup.POST("/", paymentController.CreatePayment)
	}
}
