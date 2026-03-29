package config

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() {
	var err error
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://helios_user:password@localhost:5432/helios?sslmode=disable"
	}

	DB, err = sqlx.Connect("postgres", dbURL)

	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
}
