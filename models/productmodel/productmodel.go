package productmodel

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models/categorymodel"
)

type Product struct {
	ProductId    uint
	Name         string
	Description  string
	OldPrice     int64
	Price        int64
	Image        string
	Quantity     int
	SoldQuantity int
	Category     categorymodel.Category
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func GetAll() []Product {
	rows, err := initializers.DB.Query(`
	SELECT 
		products.product_id,
		products.name,
		products.description,
		products.old_price,
		products.price,
		products.image,
		products.quantity,
		products.sold_quantity,
		categories.name as idCategory,
		products.createdAt,
		products.updatedAt
	FROM products
	JOIN categories ON products.idCategory = categories.id`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.Description,
			&product.OldPrice,
			&product.Price,
			&product.Image,
			&product.Quantity,
			&product.SoldQuantity,
			&product.Category.Name,
			&product.CreatedAt,
			&product.UpdatedAt)

		if err != nil {
			panic(err)
		}

		products = append(products, product)
	}

	return products
}

func GetByID(id int) (Product, error) {
	var product Product

	// Thực hiện truy vấn lấy thông tin sản phẩm
	row := initializers.DB.QueryRow(`
	SELECT 
		product_id,
		name,
		description,
		old_price, 
		price,
		image, 
		quantity,
		sold_quantity, 
		idCategory 
	FROM products 
	WHERE product_id = ?`, id)
	err := row.Scan(
		&product.ProductId,
		&product.Name,
		&product.Description,
		&product.OldPrice,
		&product.Price,
		&product.Image,
		&product.Quantity,
		&product.SoldQuantity,
		&product.Category.Id)
	if err == sql.ErrNoRows {
		return product, fmt.Errorf("Product not found")
	} else if err != nil {
		return product, err
	}
	return product, nil
}

func Create(product Product) bool {
	result, err := initializers.DB.Exec(`
	INSERT INTO products
	(
	name, 
	description,
	old_price, 
	price, 
	image, 
	quantity,
	sold_quantity,
	idCategory,
	createdAt, 
	updatedAt)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		product.Name,
		product.Description,
		product.OldPrice,
		product.Price,
		product.Image,
		product.Quantity,
		product.SoldQuantity,
		product.Category.Id,
		product.CreatedAt,
		product.UpdatedAt)

	if err != nil {
		panic(err)
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	return lastInsertId > 0
}

func Detail(id int) Product {
	row := initializers.DB.QueryRow(`
	SELECT 
		products.product_id,
		products.name,
		products.description,
		products.old_price,
		products.price,
		products.image,
		products.quantity,
		products.sold_quantity,
		categories.name as idCategory,
		products.createdAt,
		products.updatedAt
	FROM products
	JOIN categories ON products.idCategory = categories.id`)

	var product Product
	err := row.Scan(
		&product.ProductId,
		&product.Name,
		&product.Description,
		&product.OldPrice,
		&product.Price,
		&product.Image,
		&product.Quantity,
		&product.SoldQuantity,
		&product.Category.Id)
	if err != nil {
		panic(err)
	}

	return product
}

func Update(id int, product Product) bool {
	result, err := initializers.DB.Exec(`UPDATE products SET
	name = ?,
	description = ?,
	old_price = ?, 
	price = ?,
	image = ?, 
	quantity = ?,
	sold_quantity = ?,
	idCategory = ?, 
	updatedAt = ?
	WHERE product_id = ?`,
		product.Name,
		product.Description,
		product.OldPrice,
		product.Price,
		product.Image,
		product.Quantity,
		product.SoldQuantity,
		product.Category.Id,
		product.UpdatedAt,
		id)

	if err != nil {
		panic(err)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	return rowsAff > 0
}

func Delete(id int) bool {
	rs, err := initializers.DB.Exec(`DELETE FROM products WHERE product_id = ?`, id)
	if err != nil {
		log.Println("Error when deleting product", err)
		return false
	}

	rowsAff, err := rs.RowsAffected()
	if err != nil {
		log.Println("Error when checking affected rows", err)
		return false
	}

	if rowsAff == 0 {
		log.Println("No product found with ID: ", err)
		return false
	}

	log.Println("Product deleted successfully, ID: ", id)

	return true
}

func GetProductsByCategory(cateId int) ([]Product, error) {
	rows, err := initializers.DB.Query(`
	SELECT 	product_id, 
			name, 
			description, 
			price, 
			image, 
			quantity, 
			sold_quantity, 
			idCategory 
		FROM products 
		WHERE idCategory = ?`, cateId)
	if err != nil {
		fmt.Println("Not found Category.Id")
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		p.Category = categorymodel.Category{}
		err := rows.Scan(
			&p.ProductId,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.Image,
			&p.Quantity,
			&p.SoldQuantity,
			&p.Category.Id)
		if err != nil {
			fmt.Println("Error when scan: ", err)
			return nil, err
		}

		products = append(products, p)
	}
	return products, nil

}

func (p Product) InStock() int {
	return p.Quantity - p.SoldQuantity
}
