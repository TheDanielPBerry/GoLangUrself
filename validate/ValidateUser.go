package validate

import (
	"errors"
	"net/mail"
	"regexp"
)


func FullName(name string) error {
	if len(name) <= 0 {
		return errors.New("Name is required")
	}
	if len(name) > 128 {
		return errors.New("Name cannot be longer than 128 characters")
	}
	return nil
}

func Email(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func Password(password string) error {
	if len(password) < 10 {
		return errors.New("Password must be at least 10 characters");
	}
	if len(password) > 64 {
		return errors.New("Password cannot be longer than 64 characters")
	}
	_, err := regexp.MatchString(`^[^\n]{10,64}$/`, password)
	if err != nil {
		return err
	}
	return nil
}
