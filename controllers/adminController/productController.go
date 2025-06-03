package admincontroller

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/categorymodel"
	"github.com/Golang-Shoppe/models/productmodel"
)

func ProductHome(w http.ResponseWriter, r *http.Request) {
	products := productmodel.GetAll()

	data := map[string]any{
		"products": products,
	}
	temp, err := template.ParseFiles("views/product/index.html")
	if err != nil {
		panic(err)
	}
	temp.Execute(w, data)

}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles("views/product/create.html")
		if err != nil {
			panic(err)
		}

		categories := categorymodel.GetAll()
		data := map[string]any{
			"categories": categories,
		}

		temp.Execute(w, data)
		return
	}

	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20) // Giới hạn 10MB
		if err != nil {
			http.Error(w, "Unable to process form", http.StatusBadRequest)
			return
		}

		// Lấy dữ liệu từ form
		name := r.FormValue("productName")
		description := r.FormValue("description")

		categoryStr := r.FormValue("category_id")
		if categoryStr == "" {
			http.Error(w, "Category is required", http.StatusBadRequest)
			return
		}

		categoryId, err := strconv.Atoi(categoryStr)
		if err != nil {
			log.Println("Error when parsing category_id", err)
			http.Error(w, "Invalid category", http.StatusBadRequest)
			return
		}

		oldPriceStr := r.FormValue("oldPrice")
		if oldPriceStr == "" {
			http.Error(w, "oldPrice is empty", http.StatusBadRequest)
			return
		}

		oldPrice, err := strconv.ParseInt(oldPriceStr, 10, 64)
		if err != nil {
			log.Println("Error when parsing oldPrice: ", err)
			return
		}

		price, err := strconv.ParseInt(r.FormValue("currentPrice"), 10, 64)
		if err != nil {
			log.Println("Error when parsing price", err)
			http.Error(w, "Invalid price", http.StatusBadRequest)
			return
		}

		quantity, err := strconv.Atoi(r.FormValue("quantity"))
		if err != nil {
			log.Println("Error when parsing quantity", err)
			http.Error(w, "Invalid quantity", http.StatusBadRequest)
			return
		}

		currentTime := time.Now()

		// Xử lý upload ảnh
		file, handler, err := r.FormFile("uploadImage")
		if err != nil {
			http.Error(w, "Image upload error", http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadDir := "assets/uploads/"
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

		// Tạo sản phẩm
		product := models.Product{
			Name:        name,
			Description: description,
			OldPrice:    oldPrice,
			Price:       price,
			Image:       imagePath,
			Quantity:    quantity,
			Category:    models.Category{Id: uint(categoryId)},
			CreatedAt:   currentTime,
			UpdatedAt:   currentTime,
		}

		if productmodel.Create(product) {
			http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		} else {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
		}
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Lấy ID từ URL
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		// Lấy thông tin sản phẩm từ database
		product, err := productmodel.GetByID(id)
		if err != nil {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		// Lấy danh sách danh mục sản phẩm (nếu có)
		categories := categorymodel.GetAll()

		// Truyền dữ liệu vào template
		data := map[string]any{
			"id":         product.ProductId,
			"product":    product,
			"categories": categories,
		}

		// Render file edit.html
		temp, err := template.ParseFiles("views/product/edit.html")
		if err != nil {
			log.Println("Template parsing error:", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to process form", http.StatusBadRequest)
			return
		}

		// Lấy dữ liệu từ form
		idStr := r.FormValue("productId")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		name := r.FormValue("productName")
		description := r.FormValue("description")
		categoryStr := r.FormValue("category_id")
		oldPrice, _ := strconv.ParseInt(r.FormValue("oldPrice"), 10, 64)
		price, _ := strconv.ParseInt(r.FormValue("currentPrice"), 10, 64)
		quantity, _ := strconv.Atoi(r.FormValue("quantity"))

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

		// Cập nhật sản phẩm
		product := models.Product{
			ProductId:   uint(id),
			Name:        name,
			Description: description,
			OldPrice:    oldPrice,
			Price:       price,
			Quantity:    quantity,
			UpdatedAt:   time.Now(),
		}

		if imagePath != "" {
			product.Image = imagePath
		}

		if categoryStr != "" {
			categoryId, _ := strconv.Atoi(categoryStr)
			product.Category = models.Category{Id: uint(categoryId)}
		}

		if productmodel.Update(int(product.ProductId), product) {
			http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
		} else {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
		}
	}
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id") // Lấy id từ URL
	if idStr == "" {
		http.Error(w, "Missing product ID", http.StatusBadRequest)
		return
	}

	// Gọi hàm xóa sản phẩm
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if !productmodel.Delete(id) {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted successfully"))
}

func DetailProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Println("Method not allowed")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Lấy ID từ query string (?id=1)
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		log.Println("Invalid product ID:", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	session, _ := initializers.Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)
	isLoggedIn := (username != "")
	product, err := productmodel.GetByID(id)

	data := map[string]any{
		"username":   username,
		"product":    product,
		"isLoggedIn": isLoggedIn,
	}

	// Truy vấn sản phẩm từ database
	if err != nil {
		log.Println("Product not found:", id)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Kiểm tra template file có tồn tại không
	templatePath := "views/product/detail.html"
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		log.Println("Template file not found:", templatePath)
		http.Error(w, "Template file not found", http.StatusInternalServerError)
		return
	}

	// Parse template chi tiết sản phẩm
	temp, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// Render template với dữ liệu sản phẩm
	err = temp.Execute(w, data)
	if err != nil {
		log.Println("Error rendering template:", err)
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
