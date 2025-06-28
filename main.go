package main

import (
	"log"
	"net/http"
	"os"

	controllers "github.com/Golang-Shoppe/controllers/routes"

	initializers "github.com/Golang-Shoppe/initializers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy file .env hoặc lỗi khi đọc file")
	}

	initializers.InitOAuth()
	initializers.ConnectDatabase()
	initializers.InitParseFiles()
	controllers.InitializersRoutes()

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at http://0.0.0.0:" + port)

	// Use 0.0.0.0 for deployment
	if err := http.ListenAndServe("0.0.0.0:"+port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
