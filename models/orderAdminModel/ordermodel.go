package orderadminmodel

import (
	"fmt"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
)

func GetAllOrders() ([]models.Order, error) {
	rows, err := initializers.DB.Query(`
		SELECT
			od.id,
			u.name, 
			od.idUser,
			od.order_date,
			od.total_amount,
			od.status 
		FROM orders od
		JOIN users u ON od.idUser = u.idUser`)
	if err != nil {
		fmt.Println("Error when query: ", err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var o models.Order

		err := rows.Scan(
			&o.IdOrder,
			&o.Username,
			&o.IdUser,
			&o.OrderDate,
			&o.TotalAmount,
			&o.Status,
		)
		if err != nil {
			fmt.Println("Error when scan: ", err)
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
