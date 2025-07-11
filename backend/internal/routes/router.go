package routes

import (
	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"github.com/gin-gonic/gin"
)

func SetupRouter(paymentRepository payment.Repository) *gin.Engine {
	router := gin.Default()

	RegisterPaymentRoutes(paymentRepository, router)

	return router
}
