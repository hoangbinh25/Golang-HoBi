package cartcontroller

import (
	"encoding/json"
	"net/http"

	"github.com/Golang-Shoppe/initializers"
)

func GetCartHandle(w http.ResponseWriter, r *http.Request) {
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

	rows, err := initializers.DB.Query(`
		SELECT p.name, p.image, c.quantity, c.price
		FROM cart_items c
		JOIN products p ON p.product_id = c.product_id
		WHERE c.idUser = ?
	`, idUser)

	if err != nil {
		http.Error(w, "Error when query", http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var cart []map[string]any
	for rows.Next() {
		var name, image string
		var quantity int
		var price int64

		rows.Scan(&name, &image, &quantity, &price)

		cart = append(cart, map[string]any{
			"name":     name,
			"image":    image,
			"quantity": quantity,
			"price":    price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
