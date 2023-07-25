package domain

import (
	"errors"
	"time"
)

var (
    ErrNoUsersFound = errors.New("no users to list")
)

type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Active      bool      `json:"-"`
        CreatedAt   time.Time `json:"-"`
}

// WrapUser wraps users for user representation
type WrapUser struct {
    User User `json:"user"` 
}

// WrapUserList wraps list of users for user representation
type WrapUserList struct {
    Users []User `json:"users"`
}

// User model for user creation
type UserCreate struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type UserService interface{
    
}
