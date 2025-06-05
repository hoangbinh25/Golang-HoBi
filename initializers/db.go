package initializers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	log.Println("Connecting to database...")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Cấu hình connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Test connection với timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Database connected successfully")
	DB = db
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := DB.PingContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintf(w, "Database connection failed: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK - Database connected")
}
