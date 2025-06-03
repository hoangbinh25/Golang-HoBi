package models

import (
	"log"

	"github.com/Golang-Shoppe/initializers"
)

type CartItem struct {
	Id          int    `json:"id"`
	IdUser      int    `json:"idUser"`
	ProductId   int    `json:"product_id"`
	ProductName string `json:"name"`
	Image       string `json:image`
	Price       int64  `json:"price"`
	Quantity    int    `json:"quantity"`
	AddedAt     string `json:"added_at"`
}

func GetCartById(id int) (CartItem, error) {
	var cart CartItem

	row := initializers.DB.QueryRow(
		`SELECT
			p.idUser,
			p.product_id,
			p.name,
			p.image,
			c.quantity,
			c.price
		FROM cart_items c
		JOIN products p ON p.product_id = c.product_id 
		WHERE idUser = ?`, id)
	err := row.Scan(
		&cart.IdUser,
		&cart.ProductId,
		&cart.ProductName,
		&cart.Image,
		&cart.Price,
		&cart.Quantity,
	)

	if err != nil {
		log.Println("Error when fetching cart: ", err)
		return cart, err
	}

	return cart, nil
}
