package usermodel

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Golang-Shoppe/initializers"
	"github.com/Golang-Shoppe/models"
)

func GetById(id int) (models.User, error) {
	var user models.User
	row := initializers.DB.QueryRow(`
	SELECT
		idUser,
		username,
		password,
		name,
		email,
		address,  
		phone,
		gender,
		email_verified
	FROM users
	WHERE idUser = ?`, id)
	err := row.Scan(&user.IdUser, &user.Username, &user.Password, &user.Name, &user.Email, &user.Address, &user.Phone, &user.Gender, &user.EmailVerified)
	if err != nil {
		log.Println("Error when fetching user: ", err)
		return user, err
	}

	return user, nil
}

func UpdateUser(id int, user models.User) error {
	result, err := initializers.DB.Exec(`
	UPDATE users SET 
	email = ?,
	address = ?,
	phone = ?,
	gender = ?,
	name = ?
	WHERE idUser = ?`,
		user.Email,
		user.Address,
		user.Phone,
		user.Gender,
		user.Name,
		id)

	if err != nil {
		log.Println("Error when updating user: ", err)
		return nil
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		log.Println("Error when getting rows affected: ", err)
		return nil
	}

	if rowAff == 0 {
		log.Println("No data changed, but update executed successfully.")
		return nil
	}

	return nil
}

func UpdatePassword(id int, user models.User) error {
	result, err := initializers.DB.Exec(`
	UPDATE users SET
	password = ? 
	WHERE idUser = ?`,
		user.Password,
		id)
	if err != nil {
		fmt.Println("Error when updating password: ", err)
		return err
	}

	rowsAff, _ := result.RowsAffected()
	fmt.Println("Rows affected: ", rowsAff)
	return nil
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	row := initializers.DB.QueryRow(
		`SELECT 
		idUser,
		username,
		name,
		password,
		email, 
		address,
		phone,
		gender,
		createAt,
		updateAt,
		role,
		email_verified
	FROM users WHERE email = ?`, email)

	err := row.Scan(
		&user.IdUser,
		&user.Username,
		&user.Name,
		&user.Password,
		&user.Email,
		&user.Address,
		&user.Phone,
		&user.Gender,
		&user.CreateAt,
		&user.UpdateAt,
		&user.Role,
		&user.EmailVerified)
	if err == sql.ErrNoRows {
		return nil, nil // không có user
	} else if err != nil {
		return nil, err // lỗi truy vấn
	}
	return &user, nil
}

func CreateUserFromGoogle(u *models.User) error {
	now := time.Now()
	u.CreateAt = &now
	u.UpdateAt = &now

	_, err := initializers.DB.Exec(`INSERT INTO users (username, name, password, email, address, phone, gender, createAt, updateAt, role, email_verified) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Username, u.Name, u.Password, u.Email, u.Address, u.Phone, u.Gender, u.CreateAt, u.UpdateAt, u.Role, true)

	return err
}

func CreateUser(u *models.User) error {
	now := time.Now()
	u.CreateAt = &now
	u.UpdateAt = &now

	_, err := initializers.DB.Exec(`INSERT INTO users (username, name, password, email, address, phone, gender, createAt, updateAt, role, email_verified) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		u.Username, u.Name, u.Password, u.Email, u.Address, u.Phone, u.Gender, u.CreateAt, u.UpdateAt, u.Role, false)

	return err
}

func UpdateEmailVerified(id int) error {
	_, err := initializers.DB.Exec(`UPDATE users SET email_verified = TRUE WHERE idUser = ?`, id)
	return err
}

func UpdatePasswordByEmail(email, newPassword string) error {
	_, err := initializers.DB.Exec(`UPDATE users SET password = ? WHERE email = ?`, newPassword, email)
	return err
}

func CalculateCartTotal(idUser int) int64 {
	var total int64

	rows, err := initializers.DB.Query(`
		SELECT p.price, c.quantity
		FROM cart_items c
		JOIN products p ON c.product_id = p.id
		WHERE c.idUser = ?
	`, idUser)
	if err != nil {
		return 0
	}
	defer rows.Close()

	for rows.Next() {
		var price int64
		var quantity int
		err := rows.Scan(&price, &quantity)
		if err != nil {
			continue
		}
		total += price * int64(quantity)
	}

	return total
}

// Email Verification functions
func CreateEmailVerification(email, tokenType string) (*models.EmailVerification, error) {
	token := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour) // Token hết hạn sau 24 giờ

	verification := &models.EmailVerification{
		Email:     email,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		Used:      false,
	}

	_, err := initializers.DB.Exec(`
		INSERT INTO email_verifications (email, token, type, expires_at, created_at, used)
		VALUES (?, ?, ?, ?, ?, ?)
	`, verification.Email, verification.Token, verification.Type, verification.ExpiresAt, verification.CreatedAt, verification.Used)

	if err != nil {
		return nil, err
	}

	return verification, nil
}

func GetEmailVerificationByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	row := initializers.DB.QueryRow(`
		SELECT id, email, token, type, expires_at, created_at, used
		FROM email_verifications
		WHERE token = ?
	`, token)

	err := row.Scan(&verification.ID, &verification.Email, &verification.Token, &verification.Type, &verification.ExpiresAt, &verification.CreatedAt, &verification.Used)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &verification, nil
}

func MarkEmailVerificationAsUsed(token string) error {
	_, err := initializers.DB.Exec(`
		UPDATE email_verifications SET used = TRUE WHERE token = ?
	`, token)
	return err
}

func DeleteExpiredEmailVerifications() error {
	_, err := initializers.DB.Exec(`
		DELETE FROM email_verifications WHERE expires_at < NOW()
	`)
	return err
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
