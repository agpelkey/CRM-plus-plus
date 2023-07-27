package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/agpelkey/clients"
	"github.com/jackc/pgx/v5/pgconn"
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
// Create creates a new user
func (u userStore) Create(ctx context.Context, user *clients.User) error {
	query := `
		INSERT INTO users(first_name, last_name, phone_number, email)
		VALUES (@first_name, @last_name, @phone_number, @email)
		RETURNING id, created_at
	`

	//args := []interface{}{user.FirstName, user.LastName, user.PhoneNumber, user.Email, user.Actived}
	args := pgx.NamedArgs{
		"first_name":   &user.FirstName,
		"last_name":    &user.LastName,
		"phone_number": &user.PhoneNumber,
		"email":        &user.Email,
	}

	err := u.db.QueryRow(ctx, query, args).Scan(&user.ID)
	if err != nil {
        fmt.Println(err)
	}

	return nil

}
*/

func (u userStore) Create(user *clients.User) error {
    query := `INSERT INTO users(first_name, last_name, phone_number, email)
            VALUES ($1, $2, $3, $4)
            RETURNING id, created_at`


    args := []interface{}{user.FirstName, user.LastName, user.PhoneNumber, user.Email}

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := u.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            fmt.Println(pgErr.Message)
            fmt.Println(pgErr.Code)
        }
    }

    return nil
}







