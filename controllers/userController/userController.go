package usercontroller

import (
	"fmt"
	"net/http"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models/categorymodel"
	"github.com/Golang-Shoppe/models/productmodel"
)

var (
	err error
)

// Home page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := initializers.Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)
	products := productmodel.GetAll()
	categories := categorymodel.GetAll()

	data := map[string]any{
		"username":   username,
		"categories": categories,
		"products":   products,
	}

	if err := initializers.Tpl.ExecuteTemplate(w, "home.html", data); err != nil {
		fmt.Println("Template render error: ", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}
