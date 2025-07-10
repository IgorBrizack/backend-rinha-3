package main

import (
	"fmt"
	"log"
	"os"

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

	paymentService := services.NewPaymentService()

	// Start workers
	workers.StartDefaultWorker(redis.GetClient(), paymentService)
	workers.StartFallbackWorker(redis.GetClient(), paymentService)

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}

	r := routes.SetupRouter()

	fmt.Printf("Running on port %s\n", port)
	r.Run(":" + port)
}
