package main

import (
	"log"
	"net/http"

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
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
