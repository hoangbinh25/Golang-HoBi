package admincontroller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/Golang-Shoppe/controllers/utils"
	"github.com/Golang-Shoppe/initializers"
	orderadminmodel "github.com/Golang-Shoppe/models/orderAdminModel"
)

func OrderHome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/admin/order.html")
	if err != nil {
		http.Error(w, "Error when parsing file", http.StatusInternalServerError)
		return
	}

	orders, err := orderadminmodel.GetAllOrders()
	if err != nil {
		http.Error(w, "Error when select orders", http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"orders": orders,
	}

	temp.Execute(w, data)
}

func ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID := r.FormValue("order_id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Bắt đầu transaction
	tx, err := initializers.DB.Begin()
	if err != nil {
		http.Error(w, "Database transaction error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Cập nhật trạng thái đơn hàng thành "confirmed"
	_, err = tx.Exec("UPDATE orders SET status = 'confirmed', updated_at = NOW() WHERE id = ?", orderID)
	if err != nil {
		http.Error(w, "Error updating order status", http.StatusInternalServerError)
		return
	}

	// Cập nhật trạng thái thanh toán thành "paid"
	_, err = tx.Exec("UPDATE payments SET status = 'paid', paid_at = NOW() WHERE order_id = ?", orderID)
	if err != nil {
		http.Error(w, "Error updating payment status", http.StatusInternalServerError)
		return
	}

	// Lấy thông tin đơn hàng và user để gửi email
	var userEmail, username string
	var totalAmount int64
	err = tx.QueryRow(`
		SELECT u.email, u.username, o.total_amount 
		FROM orders o 
		JOIN users u ON o.idUser = u.idUser 
		WHERE o.id = ?
	`, orderID).Scan(&userEmail, &username, &totalAmount)
	if err != nil {
		http.Error(w, "Error getting order details", http.StatusInternalServerError)
		return
	}

	// Cập nhật số lượng đã bán cho các sản phẩm trong đơn hàng
	_, err = tx.Exec(`
		UPDATE products p 
		JOIN order_items oi ON p.product_id = oi.product_id 
		SET p.sold_quantity = p.sold_quantity + oi.quantity 
		WHERE oi.order_id = ?
	`, orderID)
	if err != nil {
		http.Error(w, "Error updating sold quantity", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		http.Error(w, "Error committing transaction", http.StatusInternalServerError)
		return
	}

	// Gửi email xác nhận cho khách hàng
	if userEmail != "" {
		subject := fmt.Sprintf("Xác nhận đơn hàng #%s đã được giao thành công", orderID)
		body := fmt.Sprintf(`
			Chào %s,

			Đơn hàng #%s của bạn đã được giao thành công!

			Chi tiết đơn hàng:
			- Mã đơn hàng: #%s
			- Tổng tiền: %d VNĐ
			- Trạng thái: Đã giao hàng

			Vui lòng xác nhận rằng bạn đã nhận được hàng bằng cách click vào link bên dưới:
			http://localhost:8080/user/confirm-received?order_id=%s

			Cảm ơn bạn đã mua hàng tại cửa hàng của chúng tôi!

			Trân trọng,
			Đội ngũ Shopee
		`, username, orderID, orderID, totalAmount, orderID)

		if err := utils.SendMail(userEmail, subject, body); err != nil {
			fmt.Printf("Error sending email: %v\n", err)
		} else {
			fmt.Printf("Confirmation email sent to %s\n", userEmail)
		}
	}

	// Trả về response thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "message": "Order confirmed successfully"}`))
}
