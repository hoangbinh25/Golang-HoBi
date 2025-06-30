package usercontroller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models/productmodel"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	products, err := productmodel.SearchProducts(query)
	if err != nil {
		http.Error(w, "Error searching products", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"products": products,
		"query":    query,
	}
	if err := initializers.Tpl.ExecuteTemplate(w, "home.html", data); err != nil {
		fmt.Println("Template render error: ", err)
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func SearchSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if len(query) < 1 {
		w.WriteHeader(http.StatusOK)
		return
	}

	limit := 10
	var suggestions []string

	rows, err := initializers.DB.Query("SELECT DISTINCT name FROM products WHERE name LIKE ? LIMIT ?", "%"+query+"%", limit)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			continue
		}
		suggestions = append(suggestions, name)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(suggestions); err != nil {
		http.Error(w, "Error encoding suggestions", http.StatusInternalServerError)
	}
}
