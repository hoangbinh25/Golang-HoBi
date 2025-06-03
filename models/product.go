package models

import "time"

type Product struct {
	ProductId    uint
	Name         string
	Description  string
	OldPrice     int64
	Price        int64
	Image        string
	Quantity     int
	SoldQuantity int
	Category     Category
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
