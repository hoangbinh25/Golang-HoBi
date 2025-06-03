package cartcontroller

import (
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/Golang-Shoppe/initializers"
)

func ViewCartHandler(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("views/cart/viewCart.html")
	if err != nil {
		log.Println("Template err: ", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	session, _ := initializers.Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)
	rawId, ok := session.Values["idUser"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var userID int
	switch v := rawId.(type) {
	case int:
		userID = v
	case int64:
		userID = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rows, err := initializers.DB.Query(`
		SELECT p.product_id,
		p.name, p.image, p.price, 
		p.quantity AS stock,
		p.sold_quantity AS sold,
		c.quantity
		FROM cart_items c
		JOIN products p ON p.product_id = c.product_id
		WHERE c.idUser = ?
	`, userID)

	if err != nil {
		http.Error(w, "Lỗi khi truy vấn giỏ hàng", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cart []map[string]any
	for rows.Next() {
		var name, image string
		var price int64
		var product_id, stock, sold, cartQuantity int
		rows.Scan(&product_id, &name, &image, &price, &stock, &sold, &cartQuantity)
		cart = append(cart, map[string]any{
			"product_id": product_id,
			"name":       name,
			"image":      image,
			"price":      price,
			"stock":      stock,
			"sold":       sold,
			"quantity":   cartQuantity,
		})
	}

	// Truyền vào template
	data := map[string]any{
		"username": username,
		"cart":     cart,
	}

	tpl.Execute(w, data)
}

func UpdateCartQuantityHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, _ := initializers.Store.Get(r, "session-name")
	rawId, ok := session.Values["idUser"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

	productId, _ := strconv.Atoi(r.FormValue("product_id"))
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	_, err := initializers.DB.Exec(`
		UPDATE cart_items SET quantity = ? WHERE idUser = ? AND product_id = ? 	
	`, quantity, idUser, productId)
	if err != nil {
		http.Error(w, "Failed update quantity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, _ := initializers.Store.Get(r, "session-name")
	rawId, ok := session.Values["idUser"]
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	productId, _ := strconv.Atoi(r.FormValue("product_id"))

	_, err := initializers.DB.Exec(`
		DELETE FROM cart_items WHERE idUser = ? AND product_id = ?  
	`, idUser, productId)
	if err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
