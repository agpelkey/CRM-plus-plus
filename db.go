package main

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Conn
}

type UserStorage interface {
	CreateUser(user *User) error
}

func (m *UserModel) CreateUser(user *User) error {
	// sql query for inserting new users into the db
	query := `INSERT INTO client_info (first_name, last_name, phone_number, email) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`

	args := []interface{}{user.FirstName, user.LastName, user.PhoneNumber, user.Email}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return m.DB.QueryRow(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}
