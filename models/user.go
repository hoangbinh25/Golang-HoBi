package models

import "time"

type User struct {
	IdUser   int        `json:"idUser`
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Password string     `json:"password"`
	Email    string     `json:"email"`
	Address  string     `json:"address"`
	Phone    string     `json:"phone"`
	Gender   string     `json:"gender"`
	CreateAt *time.Time `json:"createAt"`
	UpdateAt *time.Time `json:"updateAt"`
	Role     string     `json:"role"`
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
