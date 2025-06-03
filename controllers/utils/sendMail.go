package utils

import "net/smtp"

func SendMail(to, subject, body string) error {
	from := "binh20xx@gmail.com"
	password := "aztq gdqi zglq gfgb" // App Password nếu dùng Gmail

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, password, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
}
