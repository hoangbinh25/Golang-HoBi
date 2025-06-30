package usercontroller

import (
	"html/template"
	"net/http"

	"github.com/Golang-Shoppe/initializers"
)

func ConfirmReceivedHandler(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	// Cập nhật trạng thái đơn hàng thành "completed"
	_, err := initializers.DB.Exec("UPDATE orders SET status = 'completed', updated_at = NOW() WHERE id = ?", orderID)
	if err != nil {
		http.Error(w, "Error updating order status", http.StatusInternalServerError)
		return
	}

	// Hiển thị trang xác nhận thành công
	tmpl, err := template.ParseFiles("views/user/confirmReceived.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"orderID": orderID,
		"message": "Cảm ơn bạn đã xác nhận nhận hàng!",
	}

	tmpl.Execute(w, data)
}
