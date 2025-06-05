package main

import (
	"log"
	"net/http"
	"os"

	controllers "github.com/Golang-Shoppe/controllers"
	initializers "github.com/Golang-Shoppe/initializers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	initializers.InitOAuth()
	initializers.ConnectDatabase()
	controllers.InitParseFiles()
	controllers.InitializersRoutes()

	// Khởi động server
	log.Println("Server started at http://localhost:8080")
	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Listen trên 0.0.0.0 - QUAN TRỌNG cho Render
	address := "0.0.0.0:" + port
	log.Printf("Server starting on %s", address)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
