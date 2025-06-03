package categorymodel

import (
	"fmt"
	"log"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
)

func GetAll() []models.Category {
	rows, err := initializers.DB.Query(`SELECT id, image, name, created_at, updated_at FROM categories`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	var categories []models.Category

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(
			&category.Id,
			&category.Image,
			&category.Name,
			&category.CreateAt,
			&category.UpdateAt,
		); err != nil {
			panic(err)
		}

		categories = append(categories, category)
	}
	return categories
}

func Create(category models.Category) bool {
	result, err := initializers.DB.Exec(`
	INSERT INTO categories (image, name, created_at, updated_at)
		VALUES (?, ?, ?, ?)`,
		category.Image,
		category.Name,
		category.CreateAt,
		category.UpdateAt,
	)

	if err != nil {
		fmt.Println("Error creating category:", err)
		return false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		fmt.Println("No rows affected or error:", err)
		return false
	}

	fmt.Println("Category created successfully, rows affected:", rowsAffected)
	return true
}

func Detail(id int) models.Category {
	row := initializers.DB.QueryRow(`SELECT id, image, name FROM categories WHERE id = ?`, id)

	var category models.Category
	if err := row.Scan(&category.Id, &category.Image, &category.Name); err != nil {
		fmt.Println("Error when scan", err)
	}

	return category
}

func Update(id int, category models.Category) bool {
	result, err := initializers.DB.Exec(`
	UPDATE categories SET image = ?, name = ?, updated_at = ? WHERE id = ?`,
		category.Image, category.Name, category.UpdateAt, id)
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
	rs, err := initializers.DB.Exec(`DELETE FROM categories WHERE id = ?`, id)
	if err != nil {
		log.Println("Error when deleting category", err)
		return false
	}

	rowsAff, err := rs.RowsAffected()
	if err != nil {
		log.Println("Error when checking affedted rows: ", err)
		return false
	}

	if rowsAff == 0 {
		log.Println("No product found with ID: ", err)
		return false
	}

	log.Println("Category deleted successfully with Id: ", id)

	return true
}
