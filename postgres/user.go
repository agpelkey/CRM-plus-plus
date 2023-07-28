package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/agpelkey/clients"
	"github.com/jackc/pgx/v5"
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

func (u userStore) GetAll() ([]*clients.User, error) {
    query := `
        SELECT id, first_name, last_name, phone_number, email 
        FROM users
        ORDER BY first_name`

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := u.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*clients.User

    for rows.Next() {
        var u clients.User
        err := rows.Scan(
            &u.ID,
            &u.FirstName,
            &u.LastName,
            &u.PhoneNumber,
            &u.Email,
        )
        if err != nil {
            return nil, err
        }

        users = append(users, &u)
    }

    return users, nil
}

// for future authentication
func (u userStore) GetUserByLastName(ctx context.Context, lastName string) (clients.User, error) {
    user, err := u.List(ctx, clients.UserFilter{LastName: lastName})
    if err != nil {
        return clients.User{}, nil
    }

    return user[0], nil
}


func (u userStore) List(ctx context.Context, filter clients.UserFilter) ([]clients.User, error) {
    query := `
    SELECT id, first_name, last_name, phone_number
    FROM users
    where 1=1
    `

    rows, err := u.db.Query(ctx, query)
    if err != nil {
        return nil, fmt.Errorf("failed to query list users: %v", err)
    }

    users, err := pgx.CollectRows(rows, pgx.RowToStructByName[clients.User])
    if err != nil {
        return nil, fmt.Errorf("failed to scan rows of users: %v", err)
    }
    
    if len(users) == 0 {
        return nil, clients.ErrNoUsersFound
    }

    return users, nil
}




