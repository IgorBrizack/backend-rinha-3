package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Database struct {
	dbUser string
	dbPass string
	dbName string
	dbHost string
}

func NewDatabase() *Database {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Failed to load .env.")
	}

	return &Database{
		dbUser: os.Getenv("MYSQL_USER"),
		dbPass: os.Getenv("MYSQL_PASSWORD"),
		dbName: os.Getenv("MYSQL_DATABASE"),
		dbHost: os.Getenv("DB_HOST"),
	}
}

func (d *Database) DB() *gorm.DB {
	var db *gorm.DB
	var err error

	dsn := d.getDbDSN()

	const maxAttempts = 10
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}

		log.Printf("⚠️ Tentativa %d/%d: Falha ao conectar no banco: %v", attempts, maxAttempts, err)
		time.Sleep(6 * time.Second)
	}

	log.Fatalf("❌ Não foi possível conectar ao banco de dados após %d tentativas: %v", maxAttempts, err)
	return nil
}

func (d *Database) getDbDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		d.dbUser, d.dbPass, d.dbHost, d.dbName)
}
