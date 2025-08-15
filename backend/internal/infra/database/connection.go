package database

import (
	"fmt"
	"time"

	"github.com/IgorBrizack/backend-rinha-3/internal/domain/payment"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabaseConnection(host, port, user, password, dbname string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *gorm.DB
	var err error
	newLogger := logger.Default.LogMode(logger.Info)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err == nil {
			sqlDB, _ := db.DB()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				// Executa migrations
				if migrateErr := runMigrations(db); migrateErr != nil {
					return nil, fmt.Errorf("erro ao rodar migrations: %w", migrateErr)
				}
				return db, nil
			}
		}
		fmt.Printf("Tentativa %d: aguardando Postgres...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("não foi possível conectar ao banco: %w", err)
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&payment.Payment{}, // adicione todos os modelos aqu
	)
}

func Close(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		return
	}
	dbSQL.Close()
}
