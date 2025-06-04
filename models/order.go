package models

import (
	"time"

	"github.com/Golang-Shoppe/initializers"
)

type Order struct {
	IdOrder     int
	Username    string
	IdUser      int
	OrderDate   time.Time
	TotalAmount int64
	Status      string
}

func CreateOrder(order Order) error {
	_, err := initializers.DB.Exec(`
        INSERT INTO orders (idUser, order_date, total_amount, status)
        VALUES (?, ?, ?, ?)`,
		order.IdUser, order.OrderDate, order.TotalAmount, order.Status,
	)
	return err
}

func (o Order) FormattedOrderDate() string {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	t := time.Date(
		o.OrderDate.Year(),
		o.OrderDate.Month(),
		o.OrderDate.Day(),
		o.OrderDate.Hour(),
		o.OrderDate.Minute(),
		o.OrderDate.Second(),
		o.OrderDate.Nanosecond(),
		loc,
	)
	return t.Format("15:04:05 02/01/2006")
}
