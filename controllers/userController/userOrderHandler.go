package usercontroller

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
)

func UserOrdersHandler(w http.ResponseWriter, r *http.Request) {
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

	// Lấy thông tin user
	var user models.User
	err := initializers.DB.QueryRow("SELECT idUser, username, name, email, phone, address, gender FROM users WHERE idUser = ?", idUser).
		Scan(&user.IdUser, &user.Username, &user.Name, &user.Email, &user.Phone, &user.Address, &user.Gender)
	if err != nil {
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	// Lấy danh sách đơn hàng của user
	rows, err := initializers.DB.Query(`
		SELECT 
			o.id,
			o.total_amount,
			o.created_at,
			o.status
		FROM orders o
		WHERE o.idUser = ?
		ORDER BY o.created_at DESC
	`, idUser)
	if err != nil {
		http.Error(w, "Error getting orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrderWithItems struct {
		OrderID     int64
		TotalAmount int64
		OrderDate   string
		Status      string
		Items       []struct {
			ProductName string
			Quantity    int
			Price       int64
			Image       string
		}
	}

	var orders []OrderWithItems
	for rows.Next() {
		var order OrderWithItems
		var orderDate time.Time
		if err := rows.Scan(&order.OrderID, &order.TotalAmount, &orderDate, &order.Status); err != nil {
			fmt.Printf("Error scanning order: %v\n", err)
			http.Error(w, "Error scanning order", http.StatusInternalServerError)
			return
		}
		order.OrderDate = orderDate.Format("02/01/2006 15:04")
		fmt.Printf("Order ID: %d, Total: %d, Date: %s, Status: %s\n", order.OrderID, order.TotalAmount, order.OrderDate, order.Status)

		// Lấy chi tiết sản phẩm trong đơn hàng
		itemRows, err := initializers.DB.Query(`
			SELECT 
				p.name,
				oi.quantity,
				oi.price,
				p.image
			FROM order_items oi
			JOIN products p ON oi.product_id = p.product_id
			WHERE oi.order_id = ?
		`, order.OrderID)
		if err != nil {
			fmt.Printf("Error getting order items for order %d: %v\n", order.OrderID, err)
			http.Error(w, "Error getting order items", http.StatusInternalServerError)
			return
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item struct {
				ProductName string
				Quantity    int
				Price       int64
				Image       string
			}
			if err := itemRows.Scan(&item.ProductName, &item.Quantity, &item.Price, &item.Image); err != nil {
				fmt.Printf("Error scanning order item for order %d: %v\n", order.OrderID, err)
				http.Error(w, "Error scanning order item", http.StatusInternalServerError)
				return
			}
			fmt.Printf("Item: %s, Qty: %d, Price: %d, Image: %s\n", item.ProductName, item.Quantity, item.Price, item.Image)
			order.Items = append(order.Items, item)
		}
		orders = append(orders, order)
	}

	data := map[string]interface{}{
		"username": user.Username,
		"user":     user,
		"orders":   orders,
	}

	tmpl, err := template.ParseFiles("views/user/userOrders.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}
