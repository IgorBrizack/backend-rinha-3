package main

import (
	"fmt"
	"log"
	"os"

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

	db := database.NewDatabase().DB()

	r := routes.SetupRouter(db)

	fmt.Printf("Running on port %s\n", port)
	r.Run(":" + port)
}
