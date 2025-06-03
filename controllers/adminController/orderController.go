package admincontroller

import (
	"html/template"
	"net/http"
)

func OrderHome(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/admin/order.html")
	if err != nil {
		http.Error(w, "Error when parsing file", http.StatusInternalServerError)
		return
	}
	temp.Execute(w, nil)
}
