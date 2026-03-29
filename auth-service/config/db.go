package config

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
    connStr := "postgres://postgres:password@localhost:5432/helios?sslmode=disable"

    db, err := sql.Open("postgres", connStr)
    if err != nil {
        panic(err)
    }

    fmt.Println("Connected to DB 🚀")

    return db
}