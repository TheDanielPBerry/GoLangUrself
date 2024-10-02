package models

import (
	"errors"
	"log"
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)


type User struct {
	Id int `db:"Id"`
	Email string `db:"Email"`
	PasswordHash string `db:"PasswordHash"`
	FullName string `db:"FullName"`
	Password string
}

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}

func GetUserByEmail(email string) (*User, error) {
	db := GetDBContext()

	query := `SELECT Id, Email, PasswordHash, FullName
		FROM User
		WHERE Email=?`

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Panic(err)
		return nil, errors.New("Fatal Error")
	}

	var user User
	err = stmt.Get(&user, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		log.Panic(err)
		return nil, errors.New("Fatal Error")
	}

	return &user, nil
}


func GetUserById(id int) (*User, bool) {
	db := GetDBContext()

	query := `SELECT Id, Email, PasswordHash, FullName
		FROM User
		WHERE Id=?`

	stmt, err := db.Preparex(query)
	if err != nil {
		log.Panic(err)
		return nil, false
	}

	var user User
	err = stmt.Get(&user, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Panic(err)
		return nil, false
	}

	return &user, true
}


func CreateUser(user User) (int, error) {
	db := GetDBContext()
	
	fetchedUser, err := GetUserByEmail(user.Email)
	if err == nil || fetchedUser != nil {
		return 0, errors.New("User Already Exists")
	}
	
	sql:= `INSERT INTO User 
		(Email, PasswordHash, FullName)
		VALUES
		(?,?,?)`
	
	stmt, err := db.Preparex(sql)
	if err != nil {
		log.Panic(err)
		return 0, errors.New("Fatal Error")
	}
	user.PasswordHash, err = HashPassword(user.Password)
	if err != nil {
		log.Print(err)
		return 0, errors.New("Invalid Password")
	}

	result, err := stmt.Exec(user.Email, user.PasswordHash, user.FullName)
	if err != nil {
		log.Print(err)
		return 0, errors.New("Could not create user")
	}
	id, err := result.LastInsertId()
	return int(id), err
}


func AuthenticateLogin(email string, password string) (*User, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		log.Print(err)
		return nil, errors.New("Invalid Credentials") 
	}
	if user == nil {
		return nil, errors.New("Invalid Credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Print(err)
		return nil, errors.New("Invalid Credentials")
	}

	return user, nil
}

