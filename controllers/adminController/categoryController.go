package admincontroller

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/categorymodel"
)

func Index(w http.ResponseWriter, r *http.Request) {
	categories := categorymodel.GetAll()
	data := map[string]any{
		"categories": categories,
	}

	temp, err := template.ParseFiles("views/category/index.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, data)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/category/create.html")
		if err != nil {
			http.Error(w, "Error when parsing file .html", http.StatusBadRequest)
			return
		}

		categories := categorymodel.GetAll()

		data := map[string]any{
			"categories": categories,
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		categoryName := r.FormValue("categoryName")

		currentTime := time.Now()

		// upload ảnh
		file, handler, err := r.FormFile("uploadImage")
		if err != nil {
			http.Error(w, "Image upload error", http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadDir := "assets/uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err := os.MkdirAll(uploadDir, 0755)
			if err != nil {
				log.Println("Error creating upload directory", err)
				http.Error(w, "Failed to create directory", http.StatusInternalServerError)
				return
			}
		}

		imagePath := uploadDir + handler.Filename
		dst, err := os.Create(imagePath)

		if err != nil {
			log.Println("Error when creating file", err)
			http.Error(w, "Failed to save image", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			log.Println("Error when copying file", err)
			http.Error(w, "Failed to upload image", http.StatusInternalServerError)
			return
		}
		log.Println("Image saved successfully", imagePath)

		category := models.Category{
			Name:     categoryName,
			Image:    &imagePath,
			CreateAt: currentTime,
			UpdateAt: currentTime,
		}
		if categorymodel.Create(category) {
			http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
		} else {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
		}
	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		category := categorymodel.Detail(id)
		data := map[string]any{
			"category": category,
		}

		temp, err := template.ParseFiles("views/category/edit.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {

		idString := r.FormValue("categoryId")
		id, err := strconv.Atoi(idString)

		if err != nil {
			fmt.Println("Id Category", id)
			http.Error(w, "Invalid Id", http.StatusBadRequest)
		}

		Name := r.FormValue("categoryName")

		// Xử lý upload ảnh mới nếu có
		var imagePath string
		file, handler, err := r.FormFile("uploadImage")
		if err == nil {
			defer file.Close()
			uploadDir := "assets/uploads/"
			imagePath = uploadDir + handler.Filename
			dst, err := os.Create(imagePath)
			if err == nil {
				defer dst.Close()
				io.Copy(dst, file)
			}
		}

		category := models.Category{
			Id:       uint(id),
			Name:     Name,
			UpdateAt: time.Now(),
		}

		if imagePath != "" {
			category.Image = &imagePath
		}

		if categorymodel.Update(int(category.Id), category) {
			http.Redirect(w, r, "/admin/categories", http.StatusSeeOther)
		} else {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
		}
	}
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get id
	idString := r.URL.Query().Get("id")
	if idString == "" {
		http.Error(w, "Missing category id", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Invalid category id", http.StatusBadRequest)
		return
	}

	if !categorymodel.Delete(id) {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Delete successfully"))
}
