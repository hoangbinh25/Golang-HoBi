package usercontroller

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/usermodel"
	"golang.org/x/crypto/bcrypt"
)

func ProfileUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := initializers.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	username, _ := session.Values["username"].(string)

	// Đọc và convert idUser
	rawID, ok := session.Values["idUser"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var id int
	switch v := rawID.(type) {
	case int:
		id = v
	case int64:
		id = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Lấy user
	user, err := usermodel.GetById(id)
	if err != nil {
		log.Println("Error getting user by ID:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"username":   username,
		"user":       user,
		"gender":     user.Gender,
		"isLoggedIn": true,
	}

	// Render template
	initializers.Tpl.ExecuteTemplate(w, "profileUser.html", data)
}

func ProfileUserUpdateHandler(w http.ResponseWriter, r *http.Request) {
	session, err := initializers.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	rawID, ok := session.Values["idUser"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var id int
	switch v := rawID.(type) {
	case int:
		id = v
	case int64:
		id = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	req := models.UpdateProfileRequest{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Phone:   r.FormValue("phone"),
		Address: r.FormValue("address"),
		Gender:  r.FormValue("gender"),
	}

	oldUser, err := usermodel.GetById(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if req.Name == "" {
		req.Name = oldUser.Name
	}

	if req.Email == "" {
		req.Email = oldUser.Email
	}

	if req.Address == "" {
		req.Address = oldUser.Address
	}

	if req.Phone == "" {
		req.Phone = oldUser.Phone
	}

	if req.Gender == "" {
		req.Gender = oldUser.Gender
	}

	// Update database
	updateUser := models.User{
		Email:   req.Email,
		Address: req.Address,
		Phone:   req.Phone,
		Gender:  req.Gender,
		Name:    req.Name,
	}

	err = usermodel.UpdateUser(id, updateUser)
	if err != nil {
		fmt.Println("Error updating user: ", err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

}

func ProfileUserUpdatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	session, err := initializers.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	username, _ := session.Values["username"].(string)

	rawID, ok := session.Values["idUser"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var id int
	switch v := rawID.(type) {
	case int:
		id = v
	case int64:
		id = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		user, err := usermodel.GetById(id)
		if err != nil {
			log.Println("Error getting user by ID: ", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		data := map[string]interface{}{
			"username":   username,
			"user":       user,
			"isLoggedIn": true,
		}

		initializers.Tpl.ExecuteTemplate(w, "changePassword.html", data)
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Println("Lỗi parse multipart form:", err)
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		req := models.UpdatePasswordRequest{
			CurrentPassword: strings.TrimSpace(r.FormValue("currentPassword")),
			NewPassword:     strings.TrimSpace(r.FormValue("newPassword")),
			ConfirmPassword: strings.TrimSpace(r.FormValue("confirmPassword")),
		}

		if req.NewPassword != req.ConfirmPassword {
			http.Error(w, "New passwords do not match", http.StatusBadRequest)
			return
		}

		user, err := usermodel.GetById(id)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(strings.TrimSpace(user.Password)), []byte(req.CurrentPassword)); err != nil {
			fmt.Println("Error comparing password: ", err)
			http.Error(w, "Current password is incorrect", http.StatusBadRequest)
			return
		}

		hashedNew, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing new password", http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedNew)

		if err := usermodel.UpdatePassword(id, user); err != nil {
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"username":   username,
			"user":       user,
			"isLoggedIn": true,
			"success":    "Password updated successfully",
		}

		initializers.Tpl.ExecuteTemplate(w, "changePassword.html", data)
	}

}
