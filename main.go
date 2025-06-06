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

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at http://0.0.0.0:" + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}
