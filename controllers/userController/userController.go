package usercontroller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/Golang-Shoppe/controllers/utils"
	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/categorymodel"
	"github.com/Golang-Shoppe/models/productmodel"
	"github.com/Golang-Shoppe/models/usermodel"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var (
	err error
)

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

	initializers.Tpl.ExecuteTemplate(w, "home.html", data)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := initializers.Tpl.ExecuteTemplate(w, "login.html", nil); err != nil {
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
		initializers.Tpl.ExecuteTemplate(w, "login.html", "Vui lòng xác nhận email trước khi đăng nhập. Kiểm tra hộp thư của bạn.")
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

		if role == "admin" {
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
		} else if role == "user" {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
		} else {
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

		http://localhost:8080/verify-email?token=%s

		Link này sẽ hết hạn sau 24 giờ.

		Nếu bạn không thực hiện đăng ký này, vui lòng bỏ qua email này.

		Trân trọng,
		Đội ngũ HoBi
	`, user.Username, verification.Token)

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

func UserOrdersHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := initializers.Store.Get(r, "session-name")
	rawId, ok := session.Values["idUser"]
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	var idUser int
	switch v := rawId.(type) {
	case int:
		idUser = v
	case int64:
		idUser = int(v)
	default:
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Lấy thông tin user
	var user models.User
	err := initializers.DB.QueryRow("SELECT idUser, username, name, email, phone, address, gender FROM users WHERE idUser = ?", idUser).
		Scan(&user.IdUser, &user.Username, &user.Name, &user.Email, &user.Phone, &user.Address, &user.Gender)
	if err != nil {
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	// Lấy danh sách đơn hàng của user
	rows, err := initializers.DB.Query(`
		SELECT 
			o.id,
			o.total_amount,
			o.created_at,
			o.status
		FROM orders o
		WHERE o.idUser = ?
		ORDER BY o.created_at DESC
	`, idUser)
	if err != nil {
		http.Error(w, "Error getting orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type OrderWithItems struct {
		OrderID     int64
		TotalAmount int64
		OrderDate   string
		Status      string
		Items       []struct {
			ProductName string
			Quantity    int
			Price       int64
			Image       string
		}
	}

	var orders []OrderWithItems
	for rows.Next() {
		var order OrderWithItems
		var orderDate time.Time
		if err := rows.Scan(&order.OrderID, &order.TotalAmount, &orderDate, &order.Status); err != nil {
			fmt.Printf("Error scanning order: %v\n", err)
			http.Error(w, "Error scanning order", http.StatusInternalServerError)
			return
		}
		order.OrderDate = orderDate.Format("02/01/2006 15:04")
		fmt.Printf("Order ID: %d, Total: %d, Date: %s, Status: %s\n", order.OrderID, order.TotalAmount, order.OrderDate, order.Status)

		// Lấy chi tiết sản phẩm trong đơn hàng
		itemRows, err := initializers.DB.Query(`
			SELECT 
				p.name,
				oi.quantity,
				oi.price,
				p.image
			FROM order_items oi
			JOIN products p ON oi.product_id = p.product_id
			WHERE oi.order_id = ?
		`, order.OrderID)
		if err != nil {
			fmt.Printf("Error getting order items for order %d: %v\n", order.OrderID, err)
			http.Error(w, "Error getting order items", http.StatusInternalServerError)
			return
		}
		defer itemRows.Close()

		for itemRows.Next() {
			var item struct {
				ProductName string
				Quantity    int
				Price       int64
				Image       string
			}
			if err := itemRows.Scan(&item.ProductName, &item.Quantity, &item.Price, &item.Image); err != nil {
				fmt.Printf("Error scanning order item for order %d: %v\n", order.OrderID, err)
				http.Error(w, "Error scanning order item", http.StatusInternalServerError)
				return
			}
			fmt.Printf("Item: %s, Qty: %d, Price: %d, Image: %s\n", item.ProductName, item.Quantity, item.Price, item.Image)
			order.Items = append(order.Items, item)
		}
		orders = append(orders, order)
	}

	data := map[string]interface{}{
		"username": user.Username,
		"user":     user,
		"orders":   orders,
	}

	tmpl, err := template.ParseFiles("views/user/userOrders.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

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

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		initializers.Tpl.ExecuteTemplate(w, "forgotPassword.html", nil)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		email := r.FormValue("email")
		if email == "" {
			initializers.Tpl.ExecuteTemplate(w, "forgotPassword.html", "Email không được để trống")
			return
		}

		// Kiểm tra email có tồn tại không
		user, err := usermodel.FindUserByEmail(email)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if user == nil {
			initializers.Tpl.ExecuteTemplate(w, "forgotPassword.html", "Email không tồn tại trong hệ thống")
			return
		}

		// Tạo token reset password
		verification, err := usermodel.CreateEmailVerification(email, "forgot_password")
		if err != nil {
			http.Error(w, "Error creating reset token", http.StatusInternalServerError)
			return
		}

		// Gửi email reset password
		subject := "Đặt lại mật khẩu - HoBi"
		body := fmt.Sprintf(`
			Chào %s,

			Bạn đã yêu cầu đặt lại mật khẩu cho tài khoản %s.

			Click vào link bên dưới để đặt lại mật khẩu:

			http://localhost:8080/reset-password?token=%s

			Link này sẽ hết hạn sau 24 giờ.

			Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này.

			Trân trọng,
			Đội ngũ HoBi
		`, user.Username, verification.Token)

		if err := utils.SendMail(email, subject, body); err != nil {
			http.Error(w, "Error sending reset email", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"email":   email,
			"message": "Email đặt lại mật khẩu đã được gửi. Vui lòng kiểm tra hộp thư của bạn.",
		}
		initializers.Tpl.ExecuteTemplate(w, "resetPasswordSent.html", data)
	}
}

func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token không hợp lệ", http.StatusBadRequest)
		return
	}

	// Lấy thông tin verification
	verification, err := usermodel.GetEmailVerificationByToken(token)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if verification == nil {
		http.Error(w, "Token không tồn tại", http.StatusBadRequest)
		return
	}

	if verification.Used {
		http.Error(w, "Token đã được sử dụng", http.StatusBadRequest)
		return
	}

	if verification.Type != "register" {
		http.Error(w, "Token không hợp lệ", http.StatusBadRequest)
		return
	}

	if time.Now().After(verification.ExpiresAt) {
		http.Error(w, "Token đã hết hạn", http.StatusBadRequest)
		return
	}

	// Cập nhật trạng thái email_verified cho user
	user, err := usermodel.FindUserByEmail(verification.Email)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User không tồn tại", http.StatusBadRequest)
		return
	}

	if err := usermodel.UpdateEmailVerified(user.IdUser); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	// Đánh dấu token đã được sử dụng
	if err := usermodel.MarkEmailVerificationAsUsed(token); err != nil {
		http.Error(w, "Error updating token", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"message": "Xác nhận email thành công! Bạn có thể đăng nhập ngay bây giờ.",
	}
	initializers.Tpl.ExecuteTemplate(w, "emailVerified.html", data)
}

func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Token không hợp lệ", http.StatusBadRequest)
			return
		}

		data := map[string]interface{}{
			"token": token,
		}
		initializers.Tpl.ExecuteTemplate(w, "resetPassword.html", data)
		return
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		token := r.FormValue("token")
		newPassword := r.FormValue("newPassword")
		confirmPassword := r.FormValue("confirmPassword")

		if newPassword != confirmPassword {
			initializers.Tpl.ExecuteTemplate(w, "resetPassword.html", map[string]interface{}{
				"token": token,
				"error": "Mật khẩu không khớp",
			})
			return
		}

		// Validate password
		if isValid, msg := validatePassword(newPassword); !isValid {
			initializers.Tpl.ExecuteTemplate(w, "resetPassword.html", map[string]interface{}{
				"token": token,
				"error": msg,
			})
			return
		}

		// Kiểm tra token
		verification, err := usermodel.GetEmailVerificationByToken(token)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if verification == nil {
			http.Error(w, "Token không tồn tại", http.StatusBadRequest)
			return
		}

		if verification.Used {
			http.Error(w, "Token đã được sử dụng", http.StatusBadRequest)
			return
		}

		if verification.Type != "forgot_password" {
			http.Error(w, "Token không hợp lệ", http.StatusBadRequest)
			return
		}

		if time.Now().After(verification.ExpiresAt) {
			http.Error(w, "Token đã hết hạn", http.StatusBadRequest)
			return
		}

		// Hash password mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		// Cập nhật password
		if err := usermodel.UpdatePasswordByEmail(verification.Email, string(hashedPassword)); err != nil {
			http.Error(w, "Error updating password", http.StatusInternalServerError)
			return
		}

		// Đánh dấu token đã được sử dụng
		if err := usermodel.MarkEmailVerificationAsUsed(token); err != nil {
			http.Error(w, "Error updating token", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"message": "Đặt lại mật khẩu thành công! Bạn có thể đăng nhập với mật khẩu mới.",
		}
		initializers.Tpl.ExecuteTemplate(w, "passwordResetSuccess.html", data)
	}
}
