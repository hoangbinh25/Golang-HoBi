package admincontroller

import (
	"html/template"
	"net/http"

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
