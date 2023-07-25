package clients

import (
	"context"
	"errors"
	"time"
)

var (
	errFirstNameRequired = errors.New("clients first name is required")
	errLastNameRequired  = errors.New("clients last name is required")

	errPhoneNumberRequired = errors.New("phone number is required")

	errEmailRequired = errors.New("email is required")
	errEmailTooLong  = errors.New("email length is too long")
	errEmailInvalid  = errors.New("email is not valid")
)

type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Actived     bool      `json:"-"`
	CreatedAt   time.Time `json:"-"`
}

// WrapUser wraps user for user representation
type WrapUser struct {
	User User `json:"user"`
}

// WrapUserList wraps list of users for representation
type WrapUserList struct {
	User []User `json:"users"`
}

type UserCreate struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

// UserService is an interface for managing clients
type UserService interface {
	Create(ctx context.Context, user *User) error
}

// Validate is called upon POST requests
func (u UserCreate) Validate() error {
	switch {
	case u.FirstName == "":
		return errFirstNameRequired
	case u.LastName == "":
		return errLastNameRequired
	case len(u.Email) >= 500:
		return errEmailTooLong
	case u.PhoneNumber == "":
		return errPhoneNumberRequired
	}
	return nil
}
