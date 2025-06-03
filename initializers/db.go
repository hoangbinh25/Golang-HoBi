package initializers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDatabase() {
	// Connect to database
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/go_ecommerce?parseTime=true")

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
