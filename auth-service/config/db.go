package config

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://helios_user:password@localhost:5432/helios?sslmode=disable"
	}

	for i := 0; i < 5; i++ {
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err == nil {
			log.Println("DB connection successful")
			return
		}
		log.Printf("DB connection failed, retrying in 2s... (%d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	log.Fatal("DB connection failed permanently:", err)
}
