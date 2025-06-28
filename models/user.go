package models

import "time"

type User struct {
	IdUser        int        `json:"idUser"`
	Name          string     `json:"name"`
	Username      string     `json:"username"`
	Password      string     `json:"password"`
	Email         string     `json:"email"`
	Address       string     `json:"address"`
	Phone         string     `json:"phone"`
	Gender        string     `json:"gender"`
	EmailVerified bool       `json:"email_verified"`
	CreateAt      *time.Time `json:"createAt"`
	UpdateAt      *time.Time `json:"updateAt"`
	Role          string     `json:"role"`
}

type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Gender  string `json:"gender"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}

type EmailVerification struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	Type      string    `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	Used      bool      `json:"used"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"newPassword"`
	ConfirmPassword string `json:"confirmPassword"`
}
