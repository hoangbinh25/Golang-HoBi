package admincontroller

import (
	"html/template"
	"net/http"
)

func AdminHomeHandler(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/admin/admin.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}
