package usercontroller

import (
	"net/http"
	"time"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models/usermodel"
	"golang.org/x/crypto/bcrypt"
)

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
