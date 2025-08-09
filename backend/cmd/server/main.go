package main

import (
	"fmt"
	"log"
	"os"

	"github.com/IgorBrizack/backend-rinha-3/internal/infra/database"
	"github.com/IgorBrizack/backend-rinha-3/internal/infra/workers"
	"github.com/IgorBrizack/backend-rinha-3/internal/routes"
	"github.com/IgorBrizack/backend-rinha-3/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("[INIT] Carregando variáveis de ambiente...")
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARN] Arquivo .env não encontrado, prosseguindo com variáveis de ambiente do sistema")
	} else {
		fmt.Println("[OK] Variáveis de ambiente carregadas")
	}

	conn, err := database.NewDatabaseConnection(
		"postgres", "5432", "postgres", "postgres", "app_db",
	)
	if err != nil {
		log.Fatalf("[FATAL] Falha ao conectar ao banco de dados: %v", err)
	}
	defer database.Close(conn)

	paymentQueue := make(chan []byte, 50000)
	defaultQueue := make(chan []byte, 50000)
	fallbackQueue := make(chan []byte, 50000)

	paymentRepository := database.NewPaymentRepository(conn)
	paymentService := services.NewPaymentService()

	workers.NewPaymentWorker(
		paymentService,
		paymentRepository,
		paymentQueue,
		defaultQueue,
		fallbackQueue).Start()

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
		fmt.Println("[WARN] Variável BACKEND_PORT não definida, usando porta padrão 8081")
	}

	fmt.Printf("[INIT] Iniciando servidor HTTP na porta %s...\n", port)

	r := routes.SetupRouter(paymentRepository, paymentQueue)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[FATAL] Falha ao iniciar servidor: %v", err)
	}
}
