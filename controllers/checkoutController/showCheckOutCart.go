package checkoutcontroller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"

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
		var productId, stock, sold, cartQuantity int

		err := rows.Scan(
			&productId,
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
			"productId": productId,
			"name":      name,
			"image":     image,
			"price":     price,
			"stock":     stock,
			"sold":      sold,
			"quantity":  cartQuantity,
			"total":     total,
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
