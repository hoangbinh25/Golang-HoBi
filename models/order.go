package models

import (
	"time"

	"github.com/Golang-Shoppe/initializers"
)

type Order struct {
	IdUser    int
	OrderDate time.Time
	Total     int64
	Status    string
}

func CreateOrder(order Order) error {
	_, err := initializers.DB.Exec(`
        INSERT INTO orders (idUser, order_date, total_amount, status)
        VALUES (?, ?, ?, ?)`,
		order.IdUser, order.OrderDate, order.Total, order.Status,
	)
	return err
}
