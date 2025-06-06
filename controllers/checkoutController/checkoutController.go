package checkoutcontroller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"

	"github.com/Golang-Shoppe/controllers/utils"
	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models/usermodel"
)

func ShowCheckoutPage(w http.ResponseWriter, r *http.Request) {
	// render template checkout.html
	tmpl := template.Must(template.ParseFiles("views/product/checkout.html"))

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := initializers.Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)
	rawId, ok := session.Values["idUser"]
	if !ok {
		fmt.Println("User ID not found in session")
		http.Error(w, "User ID not found", http.StatusNotFound)
		return
	}

	var idUser int
	switch v := rawId.(type) {
	case int:
		idUser = v
	case int64:
		idUser = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := usermodel.GetById(idUser)
	if err != nil {
		http.Error(w, "Id not found", http.StatusNotFound)
		return
	}

	rows, err := initializers.DB.Query(`
		SELECT 
		p.product_id,
		p.name,
		p.image,
		p.price, 
		p.quantity AS stock,
		p.sold_quantity AS sold,
		c.quantity
		FROM cart_items c
		JOIN products p ON p.product_id = c.product_id
		WHERE c.idUser = ?
	`, idUser)

	if err != nil {
		http.Error(w, "Lỗi khi truy vấn giỏ hàng", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cart []map[string]any
	var total int64

	for rows.Next() {
		var name, image string
		var price int64
		var product_id, stock, sold, cartQuantity int

		err := rows.Scan(
			&product_id,
			&name,
			&image,
			&price,
			&stock,
			&sold,
			&cartQuantity)

		if err != nil {
			_, file, line, _ := runtime.Caller(1)
			log.Printf("[ERROR] %s:%d %v\n", file, line, err)
		}

		itemTotal := price * int64(cartQuantity)

		total += itemTotal

		cart = append(cart, map[string]any{
			"product_id": product_id,
			"name":       name,
			"image":      image,
			"price":      price,
			"stock":      stock,
			"sold":       sold,
			"quantity":   cartQuantity,
			"total":      total,
		})
	}

	data := map[string]any{
		"username": username,
		"user":     user,
		"cart":     cart,
		"total":    total,
	}

	tmpl.Execute(w, data)
}

func CheckOutCartHandle(w http.ResponseWriter, r *http.Request) {
	session, _ := initializers.Store.Get(r, "session-name")
	rawId, ok := session.Values["idUser"]
	if !ok {
		http.Error(w, "User ID not found", http.StatusNotFound)
		return
	}

	fmt.Println("Session at checkout:", session.Values)

	var idUser int
	switch v := rawId.(type) {
	case int:
		idUser = v
	case int64:
		idUser = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Lấy giỏ hàng
	rows, err := initializers.DB.Query(`
		SELECT 
			product_id,
			quantity,
			price 
		FROM cart_items 
		WHERE idUser = ?`, idUser)
	if err != nil {
		http.Error(w, "Error selecting cart", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type CartItem struct {
		ProductID int
		Quantity  int
		Price     int64
	}

	var cartItems []CartItem
	var totalAmount int64

	for rows.Next() {
		var item CartItem
		if err := rows.Scan(&item.ProductID, &item.Quantity, &item.Price); err != nil {
			http.Error(w, "Error scanning cart", http.StatusInternalServerError)
			return
		}
		cartItems = append(cartItems, item)
		totalAmount += item.Price * int64(item.Quantity)
	}

	if len(cartItems) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	}

	// Bắt đầu transaction
	tx, err := initializers.DB.Begin()
	if err != nil {
		http.Error(w, "DB transaction error", http.StatusInternalServerError)
		return
	}

	// Tạo đơn hàng
	res, err := tx.Exec("INSERT INTO orders (idUser, total_amount, order_date) VALUES (?, ?, NOW())", idUser, totalAmount)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	orderID, _ := res.LastInsertId()

	// Thêm vào order_items
	for _, item := range cartItems {
		_, err := tx.Exec(`
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES (?, ?, ?, ?)`,
			orderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create order items", http.StatusInternalServerError)
			return
		}
	}

	// Xoá giỏ hàng
	if _, err := tx.Exec("DELETE FROM cart_items WHERE idUser = ?", idUser); err != nil {
		tx.Rollback()
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	// Commit
	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Gửi email xác nhận
	var email, username string
	err = initializers.DB.QueryRow("SELECT email, username FROM users WHERE idUser = ?", idUser).Scan(&email, &username)
	if err == nil && email != "" {
		subject := fmt.Sprintf("Xác nhận đơn hàng #%d", orderID)
		body := fmt.Sprintf("Chào %s,\n\nCảm ơn bạn đã đặt hàng tại cửa hàng của chúng tôi.\nMã đơn hàng của bạn là #%d.\nChúng tôi sẽ xử lý đơn hàng trong thời gian sớm nhất.\n\nTrân trọng!", username, orderID)
		if err := utils.SendMail(email, subject, body); err != nil {
			fmt.Println("Gửi mail thất bại:", err)
		} else {
			fmt.Println("Đã gửi email xác nhận đơn hàng tới", email)
		}
	} else {
		fmt.Println("Không thể lấy email người dùng:", err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
