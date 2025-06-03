package oauth

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
	"github.com/Golang-Shoppe/models/usermodel"
)

type GoogleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := initializers.GoogleOAuthConfig.AuthCodeURL("random-state")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	token, err := initializers.GoogleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Lỗi token", http.StatusInternalServerError)
		return
	}

	client := initializers.GoogleOAuthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Lỗi lấy user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var gUser GoogleUser
	if err := json.Unmarshal(body, &gUser); err != nil {
		http.Error(w, "Lỗi giải mã JSON", http.StatusInternalServerError)
		return
	}

	// Kiểm tra xem user đã có trong DB chưa
	user, err := usermodel.FindUserByEmail(gUser.Email)
	if err != nil {
		http.Error(w, "Lỗi DB", http.StatusInternalServerError)
		return
	}

	if user == nil {
		newUser := &models.User{
			Username: strings.Split(gUser.Email, "@")[0],
			Name:     gUser.Name,
			Email:    gUser.Email,
			Role:     "user",
			Password: "",
			Address:  "",
			Phone:    "",
			Gender:   "",
		}
		usermodel.CreateUserFromGoogle(newUser)
		user = newUser
	}

	// Gán session
	session, _ := initializers.Store.Get(r, "session-name")
	session.Values["username"] = user.Username
	session.Values["idUser"] = user.IdUser
	session.Values["email"] = user.Email
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
