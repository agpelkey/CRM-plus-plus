package postgres

import (
	"context"

	"github.com/agpelkey/clients"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userStore represents users database
type userStore struct {
	db *pgxpool.Pool
}

// NewUserStore returns an instance of UsersStore
func NewUserStore(db *pgxpool.Pool) userStore {
	return userStore{db: db}
}

/*
type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Actived     bool      `json:"-"`
	CreatedAt   time.Time `json:"-"`
}
*/

// Create creates a new user
func (u userStore) Create(ctx context.Context, user *clients.User) error {
	query := `
		INSERT INTO users(first_name, last_name, phone_number, email, activated)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	//args := []interface{}{user.FirstName, user.LastName, user.PhoneNumber, user.Email, user.Actived}
	args := pgx.NamedArgs{
		"first_name":   &user.FirstName,
		"last_name":    &user.LastName,
		"phone_number": &user.PhoneNumber,
		"email":        &user.Email,
		"activated":    &user.Actived,
	}

	err := u.db.QueryRow(ctx, query, args).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil

}
