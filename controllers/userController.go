package controllers

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/categorymodel"
	"github.com/Golang-Shoppe/models/productmodel"
	"github.com/Golang-Shoppe/models/usermodel"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	tpl *template.Template
	err error
)

func InitParseFiles() {
	tpl, err = tpl.ParseGlob("views/**/*.html")
	tpl, err = tpl.ParseGlob("views/*.html")
	if err != nil {
		log.Fatal("Error parsing templates: ", err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Read session from browser
	session, _ := initializers.Store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)
	products := productmodel.GetAll()
	categories := categorymodel.GetAll()

	data := map[string]any{
		"username":   username,
		"categories": categories,
		"products":   products,
	}

	tpl.ExecuteTemplate(w, "home.html", data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
		fmt.Println("Error executing")
		return
	}
}

func LoginAuthHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var (
		id         int
		hash, role string
	)

	stmt := "SELECT idUser, password, role FROM `go_ecommerce`.users WHERE username = ?;"
	row := initializers.DB.QueryRow(stmt, user.Username)
	err = row.Scan(&id, &hash, &role)

	if err != nil {
		fmt.Println("Error selecting Hash in db by username")
		tpl.ExecuteTemplate(w, "login.html", "Check username or password")
		return
	}
	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
	if err == nil {
		user.IdUser = id
		user.Role = role

		session, err := initializers.Store.Get(r, "session-name")
		if err != nil {
			fmt.Println("Error getting session:", err)
			tpl.ExecuteTemplate(w, "login.html", "Session error. Please try again.")
			return
		}
		session.Values["username"] = user.Username
		session.Values["role"] = user.Role
		session.Values["idUser"] = user.IdUser

		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   300,  // Second
			HttpOnly: true, // Cookie only http
			SameSite: http.SameSiteLaxMode,
		}

		if err := session.Save(r, w); err != nil {
			fmt.Println("Error saving session")
			tpl.ExecuteTemplate(w, "login.html", "A error. Please try again")
			return
		}

		if role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else if role == "user" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
			tpl.ExecuteTemplate(w, "login.html", "Invalid role configuration")
		}
		return

	}

	tpl.ExecuteTemplate(w, "login.html", "Invalid username or password")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := initializers.Store.Get(r, "session-name")
	if err != nil {
		fmt.Println("Error getting session:", err)
	}

	// Xóa hết session values
	session.Values = make(map[interface{}]interface{})

	session.Options.MaxAge = -1

	// Lưu session (gửi cookie xóa về trình duyệt)
	if err := session.Save(r, w); err != nil {
		fmt.Println("Error saving session:", err)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
		fmt.Println("Error executing template")
		return
	}
}

func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error when parsing form")
		return
	}

	db := initializers.DB

	user := models.User{
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	// Check username for only alphanumeric characters
	var nameAlphanumeric = true
	for _, char := range user.Username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			nameAlphanumeric = false
		}
	}
	// Check usernameLength
	var nameLength bool
	if len(user.Username) >= 5 && len(user.Username) <= 50 {
		nameLength = true
	}

	fmt.Println("nameAlphanumeric: ", nameAlphanumeric)

	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		tpl.ExecuteTemplate(w, "register.html", "Invalid email address")
		return
	}

	fmt.Println("nameLength: ", nameLength)

	// Check password
	fmt.Println("Password: ", user.Password, "\npasswordLength: ", len(user.Password))

	if isValid, msg := validatePassword(user.Password); !isValid {
		fmt.Println("Error: ", msg)
		return
	}

	// check username already exists
	stmt := "SELECT idUser FROM `go_ecommerce`.users WHERE username = ?;"
	row := db.QueryRow(stmt, user.Username)
	var uID string
	if err := row.Scan(&uID); err != sql.ErrNoRows {
		if err != nil {
			fmt.Println("Error checking username: ", err)
		} else {
			fmt.Println("Username already exists, err: ", err)
		}
		tpl.ExecuteTemplate(w, "register.html", "username already exists")
		return
	}

	// Create hash from password
	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err: ", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering account")
		return
	}

	fmt.Println("hash: ", hash)
	fmt.Println("String(hash): ", string(hash))

	result, err := db.Exec("INSERT INTO `go_ecommerce`.users (username, password, email) VALUES (?, ?, ?);",
		user.Username, string(hash), user.Email)

	if err != nil {
		fmt.Println("Error executing insert: ", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem registering account")
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected: ", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem")
		return
	}

	lastIns, err := result.LastInsertId()
	if err != nil {
		fmt.Println("Error getting last inserted id: ", err)
		tpl.ExecuteTemplate(w, "register.html", "There was a problem")
		return
	}

	fmt.Println("Rows affected: ", rows)
	fmt.Println("Last INSERT ID: ", lastIns)

	session, _ := initializers.Store.Get(r, "session-name")
	session.Values["username"] = user.Username
	if err := session.Save(r, w); err != nil {
		fmt.Println("Error saving session: ", err)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func ProfileUserHandler(w http.ResponseWriter, r *http.Request) {
	session, err := initializers.Store.Get(r, "session-name")
	if err != nil {
		http.Error(w, "Session error", http.StatusInternalServerError)
		return
	}

	// fmt.Println("Tất cả giá trị trong session: ", session.Values)

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
	tpl.ExecuteTemplate(w, "profileUser.html", data)
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

		tpl.ExecuteTemplate(w, "changePassword.html", data)
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

		tpl.ExecuteTemplate(w, "changePassword.html", data)
	}

}

func validatePassword(password string) (bool, string) {
	if len(password) < 8 || len(password) > 60 {
		return false, "Password must be between 8 and 60 characters"
	}
	hasNumber := false
	hasSpecial := false
	hasNoSpace := true

	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		case unicode.IsSpace(char):
			hasNoSpace = false
		}
	}

	if !hasNumber || !hasSpecial || !hasNoSpace {
		return false, "Password must contain number, a special character and no space "
	}

	return true, ""
}
