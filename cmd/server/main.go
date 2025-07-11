package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	db "github.com/IgorBrizack/backend-rinha-3/internal/infra/database"
	"github.com/IgorBrizack/backend-rinha-3/internal/infra/redis"
	"github.com/IgorBrizack/backend-rinha-3/internal/infra/workers"
	"github.com/IgorBrizack/backend-rinha-3/internal/routes"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env file.")
	}

	if err := redis.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}

	database := db.NewDatabase()
	database.DB().AutoMigrate(&payment.Payment{})

	paymentService := services.NewPaymentService()
	paymentRepository := db.NewPaymentRepository(database.DB())

	// Start workers
	workers.StartDefaultWorker(paymentRepository, redis.GetClient(), paymentService)
	workers.StartFallbackWorker(paymentRepository, redis.GetClient(), paymentService)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}

	r := routes.SetupRouter()

	fmt.Printf("Running on port %s\n", port)
	r.Run(":" + port)
}
