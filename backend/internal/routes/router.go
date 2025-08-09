package routes

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/controller"
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/gin-gonic/gin"
)

func SetupRouter(paymentRepository payment.Repository, paymentQueue chan []byte) *gin.Engine {
	router := gin.Default()

	paymentController := controller.NewPaymentController(paymentRepository, paymentQueue)

	router.POST("/payments", paymentController.CreatePayment)

	router.GET("/payments-summary", paymentController.GetPaymentSummary)

	router.DELETE("/payments", paymentController.PurgePayments)

	return router
}
