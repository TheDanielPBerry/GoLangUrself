package validate

import (
	"errors"
	"net/mail"
	"regexp"
)


func Email(email string) (error) {
	_, err := mail.ParseAddress(email)
	return err
}

func Password(password string) (error) {
	if len(password) < 10 {
		return errors.New("Password must be at least 10 characters");
	}
	if len(password) > 64 {
		return errors.New("Password cannot be longer than 64 characters")
	}
	match, err := regexp.MatchString(`^[^\n]{10,64}$/`, password)
	if err != nil {
		return err
	}
	if match {
		return nil
	}
	return errors.New("Password is Invalid")
}
