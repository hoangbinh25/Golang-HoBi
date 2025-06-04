package categorycontroller

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models/categorymodel"
	models "github.com/Golang-Shoppe/models/productmodel"
)

func ShowProductsCategory(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("views/category/detail.html")
	if err != nil {
		fmt.Println("Error when parsing file: ", err)
		http.Error(w, "Error when parsing file", http.StatusInternalServerError)
		return
	}

	session, _ := initializers.Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)
	idUser, ok := session.Values["idUser"].(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idCateStr := r.URL.Query().Get("id")
	idCate, err := strconv.Atoi(idCateStr)

	if err != nil {
		fmt.Println("Id not found: ", err)
		http.Error(w, "Id not found: ", http.StatusInternalServerError)
		return
	}

	productsCate, err := models.GetProductsByCategory(idCate)
	if err != nil {
		fmt.Println("Error when select product: ", err)
		http.Error(w, "Error when select product", http.StatusInternalServerError)
		return
	}

	categories := categorymodel.GetAll()

	data := map[string]any{
		"username":     username,
		"idUser":       idUser,
		"productsCate": productsCate,
		"categories":   categories,
	}

	tpl.Execute(w, data)

}
