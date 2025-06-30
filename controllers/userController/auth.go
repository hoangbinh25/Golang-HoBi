package usercontroller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"unicode"

	"github.com/Golang-Shoppe/controllers/utils"
	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/usermodel"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	user := models.User{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	var (
		id            int
		hash, role    string
		emailVerified bool
	)

	stmt := "SELECT idUser, password, role, email_verified FROM users WHERE username = ?;"
	row := initializers.DB.QueryRow(stmt, user.Username)
	err = row.Scan(&id, &hash, &role, &emailVerified)

	if err != nil {
		fmt.Println("Error selecting Hash in db by username")
		initializers.Tpl.ExecuteTemplate(w, "login.html", "Check username or password")
		return
	}

	// Kiểm tra email đã được xác nhận chưa
	if !emailVerified {
		initializers.Tpl.ExecuteTemplate(w, "login.html", "Please confirm email before login. Check your email.")
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
			initializers.Tpl.ExecuteTemplate(w, "login.html", "Session error. Please try again.")
			return
		}
		session.Values["username"] = user.Username
		session.Values["role"] = user.Role
		session.Values["idUser"] = user.IdUser

		session.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   3000, // Second
			HttpOnly: true, // Cookie only http
			SameSite: http.SameSiteLaxMode,
		}

		if err := session.Save(r, w); err != nil {
			fmt.Println("Error saving session")
			initializers.Tpl.ExecuteTemplate(w, "login.html", "A error. Please try again")
			return
		}

		switch role {
		case "admin":
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			log.Println("Login success, role:", role)
		case "user":
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		default:
			initializers.Tpl.ExecuteTemplate(w, "login.html", "Invalid role configuration")
		}
		return

	}

	initializers.Tpl.ExecuteTemplate(w, "login.html", "Invalid username or password")
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
	if err := initializers.Tpl.ExecuteTemplate(w, "register.html", nil); err != nil {
		fmt.Println("Error executing template")
		return
	}
}

func RegisterAuthHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error when parsing form")
		return
	}

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

	if !nameAlphanumeric {
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Username chỉ được chứa chữ cái và số")
		return
	}

	if !nameLength {
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Username phải có độ dài từ 5 đến 50 ký tự")
		return
	}

	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Invalid email address")
		return
	}

	fmt.Println("nameLength: ", nameLength)

	// Check password
	fmt.Println("Password: ", user.Password, "\npasswordLength: ", len(user.Password))

	if isValid, msg := validatePassword(user.Password); !isValid {
		fmt.Println("Error: ", msg)
		initializers.Tpl.ExecuteTemplate(w, "register.html", msg)
		return
	}

	// check username already exists
	var usernameExists bool
	err := initializers.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", user.Username).Scan(&usernameExists)
	if err != nil {
		fmt.Println("Error when checking username exists: ", err)
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Database error. Please try again.")
		return
	}

	if usernameExists {
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Username already exists")
		return
	}

	// check email already exists
	var emailExists bool
	err = initializers.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&emailExists)
	if err != nil {
		fmt.Println("Error when checking email exists: ", err)
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Database error. Please try again.")
		return
	}

	if emailExists {
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Email already exists")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error when hashing password: ", err)
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Error processing password. Please try again.")
		return
	}

	user.Password = string(hashedPassword)
	user.Role = "user"

	// Tạo user với email chưa được xác nhận
	if err := usermodel.CreateUser(&user); err != nil {
		fmt.Println("Error when creating user: ", err)
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Error creating user")
		return
	}

	// Tạo token xác nhận email
	verification, err := usermodel.CreateEmailVerification(user.Email, "register")
	if err != nil {
		fmt.Println("Error creating email verification: ", err)
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Error creating verification token")
		return
	}

	// Gửi email xác nhận
	subject := "Xác nhận đăng ký tài khoản - HoBi"
	body := fmt.Sprintf(`
		Chào %s,

		Cảm ơn bạn đã đăng ký tài khoản tại HoBi!

		Để hoàn tất quá trình đăng ký, vui lòng click vào link bên dưới để xác nhận email của bạn:

		%s/verify-email?token=%s

		Link này sẽ hết hạn sau 24 giờ.

		Nếu bạn không thực hiện đăng ký này, vui lòng bỏ qua email này.

		Trân trọng,
		Đội ngũ HoBi
	`, user.Username, os.Getenv("APP_URL"), verification.Token)

	if err := utils.SendMail(user.Email, subject, body); err != nil {
		fmt.Println("Error sending verification email: ", err)
		initializers.Tpl.ExecuteTemplate(w, "register.html", "Error sending verification email")
		return
	}

	// Hiển thị trang thông báo đã gửi email xác nhận
	data := map[string]interface{}{
		"email":   user.Email,
		"message": "Đăng ký thành công! Vui lòng kiểm tra email để xác nhận tài khoản.",
	}
	initializers.Tpl.ExecuteTemplate(w, "verifyEmailSent.html", data)
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
