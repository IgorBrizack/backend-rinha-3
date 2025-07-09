package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	userController := controller.NewController(db)

	userGroup := r.Group("/users")
	{
		userGroup.GET("/", userController.GetUsers)
		userGroup.POST("/", userController.CreateUser)
	}
}
