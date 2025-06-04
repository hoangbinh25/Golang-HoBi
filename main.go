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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
