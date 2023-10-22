package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Password  string    `json:"Password"`
	Email     string    `json:"email"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
			case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
					//invalid password
					return false, nil
			default:
					return false, err
		}
	}

	return true, nil
}
