package cartcontroller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Golang-Shoppe/initializers"
)

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Lấy session
	session, _ := initializers.Store.Get(r, "session-name")
	rawUserID, ok := session.Values["idUser"]
	if !ok {
		http.Error(w, "You are not logged in", http.StatusUnauthorized)
		return
	}

	var userID int
	switch v := rawUserID.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Lấy product_id từ form (POST)
	productIDStr := r.FormValue("product_id")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Lấy quantity từ form (POST)
	quantityStr := r.FormValue("quantity")
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity < 1 {
		quantity = 1
	}

	// Kiểm tra sản phẩm có tồn tại không
	var price int64
	err = initializers.DB.QueryRow(`
	SELECT price FROM products WHERE product_id = ?`,
		productID).Scan(&price)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Kiểm tra sản phẩm đã có trong giỏ hàng chưa
	var oldQuantity int
	err = initializers.DB.QueryRow(`
		SELECT quantity FROM cart_items WHERE idUser = ? AND product_id = ?
	`, userID, productID).Scan(&oldQuantity)

	if err == sql.ErrNoRows {
		// Chưa có → thêm mới
		_, err = initializers.DB.Exec(`
			INSERT INTO cart_items (idUser, product_id, quantity, price)
			VALUES (?, ?, ?, ?)`,
			userID, productID, quantity, price)
		if err != nil {
			http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
			return
		}
		fmt.Println("Đã thêm mới vào giỏ hàng: ", productID)
	} else if err == nil {
		// Đã có → cập nhật lại số lượng
		_, err = initializers.DB.Exec(`
			UPDATE cart_items SET quantity = ?, price = ?
			WHERE idUser = ? AND product_id = ?`,
			quantity, price, userID, productID)

		fmt.Println("Đã cập nhật số lượng giỏ hàng: ", productID)
		if err != nil {
			http.Error(w, "Failed to update cart", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Error checking cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Thêm thành công",
	})

	fmt.Println("idUser: ", userID)
	fmt.Println("ProductID: ", productID)
}
