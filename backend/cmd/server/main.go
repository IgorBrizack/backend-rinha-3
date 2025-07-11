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
	fmt.Println("[INIT] Carregando variáveis de ambiente...")
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARN] Arquivo .env não encontrado, prosseguindo com variáveis de ambiente do sistema")
	} else {
		fmt.Println("[OK] Variáveis de ambiente carregadas")
	}

	fmt.Println("[INIT] Inicializando Redis...")
	if err := redis.InitRedis(); err != nil {
		log.Fatalf("[FATAL] Falha ao inicializar Redis: %v", err)
	}
	fmt.Println("[OK] Redis inicializado com sucesso")

	fmt.Println("[INIT] Conectando ao banco de dados...")
	database := db.NewDatabase()
	if err := database.DB().AutoMigrate(&payment.Payment{}); err != nil {
		log.Fatalf("[FATAL] Falha ao migrar modelo Payment: %v", err)
	}
	fmt.Println("[OK] Banco de dados conectado e migrado")

	fmt.Println("[INIT] Inicializando repositórios e serviços...")
	paymentRepository := db.NewPaymentRepository(database.DB())
	paymentService := services.NewPaymentService()
	fmt.Println("[OK] Repositórios e serviços prontos")

	fmt.Println("[INIT] Iniciando workers...")
	workers.StartDefaultWorker(paymentRepository, redis.GetClient(), paymentService)
	fmt.Println("[OK] Worker default iniciado")
	workers.StartFallbackWorker(paymentRepository, redis.GetClient(), paymentService)
	fmt.Println("[OK] Worker fallback iniciado")

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "8081"
		fmt.Println("[WARN] Variável BACKEND_PORT não definida, usando porta padrão 8081")
	}

	fmt.Printf("[INIT] Iniciando servidor HTTP na porta %s...\n", port)
	r := routes.SetupRouter(paymentRepository)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("[FATAL] Falha ao iniciar servidor: %v", err)
	}
}
