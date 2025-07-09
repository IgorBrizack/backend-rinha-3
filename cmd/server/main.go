package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IgorBrizack/backend-rinha-3/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env file.")
	}

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
	}

	r := routes.SetupRouter()

	fmt.Printf("Running on port %s\n", port)
	r.Run(":" + port)
}
