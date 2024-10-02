package models

import (
	"time"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Email, FullName string
	password, passwordHash string
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func GetUserByEmail(email string) (*User, error) {
	sql := `SELECT id, Email, PasswordHash, FullName
		FROM User
		WHERE Email=?`

	stmt, err := db.Preparex(sql)
	if err != nil {
		return nil, err
	}

	var user *User
	err = stmt.Get(user, email)
	if err != nil {
		return nil, err
	}

}

func CreateUser(user User) int {
	sql:= "INSERT INTO user (email, password_hash, full_bane)"
	return 0
}

func AttemptLogin(email string, password string) (*User, error) {
	db := GetDBContext()

	
	err = bcrypt.CompareHashAndPassword([]byte(user.password), []byte(password))
	if err != nil {
		return nil, errors.New("Invalid Credentials")
	}

	return user, nil
}
