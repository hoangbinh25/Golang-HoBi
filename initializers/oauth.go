package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var GoogleOAuthConfig *oauth2.Config

func InitOAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Không tìm thấy file .env hoặc lỗi khi đọc file")
	}

	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://golang-hobi.onrender.com/auth/callback",
		//http://golang-hobi.onrender.com/auth/callbackhttp://localhost:8080/auth/callback
		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
}
