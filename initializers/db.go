package initializers

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile("./certs/ca.pem") // File từ Aiven
	if err != nil {
		log.Fatal("Failed to read CA cert:", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		log.Fatal("Failed to append CA cert")
	}

	err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal("Failed to register TLS config:", err)
	}

	// Connect to database
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Println("⚠️ DB_DSN not found, using fallback local DSN")
		dsn = "root:123456@tcp(localhost:3306)/go_ecommerce?parseTime=true"
	}
	log.Println("Using DSN:", dsn)

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Ping DB failure:", err)
	}

	log.Println("Connected successfully")

	if err != nil {
		fmt.Println(err.Error())
	}
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
