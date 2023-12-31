package clients

import (
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
    ErrDuplicateEmail = errors.New("email already exists")

    ErrNoUsersFound = errors.New("no users were found")
    ErrRecordNotFound = errors.New("record not found")
)

// Custom user type
type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
    FollowUp    bool      `json:"follow_up"`
    CheckInDate time.Time `json:"check_in_date"`
}

// Custome data type to represent POST requests
type UserCreate struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
    //FollowUp    bool   `json:"follow_up"`
}

// Custom data type to represent PATCH requests
type UserUpdate struct {
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
}


// UserService is an interface for managing clients
type UserService interface {
    Create(user *User) error
    GetAll() ([]*User, error)
    GetUserByID(id int64) (*User, error)
    //List(ctx context.Context, filter UserFilter) ([]User, error)
    DeleteUser(id int64) error
    UpdateUser(id int64, user UserUpdate) error
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

// used to set input values to a new struct and return a new user.
// Will mostly be used for POST requests.
func (u UserCreate) CreateModel() User {
    return User{
        FirstName: u.FirstName,
        LastName: u.LastName,
        PhoneNumber: u.PhoneNumber,
        Email: u.Email,
        //FollowUp: u.Fol,
    }
}
