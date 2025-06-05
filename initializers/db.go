package initializers

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database connected")
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("Database is not connected")
		fmt.Println(err.Error())
	}
	DB = db
}
