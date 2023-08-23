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


// create a new user entry in the database
func (u userStore) Create(user *clients.User) error {
    query := `INSERT INTO users(first_name, last_name, phone_number, email)
            VALUES ($1, $2, $3, $4)
            RETURNING id`


    args := []interface{}{user.FirstName, user.LastName, user.PhoneNumber, user.Email}

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    err := u.db.QueryRow(ctx, query, args...).Scan(&user.ID) 
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            fmt.Println(pgErr.Message)
            fmt.Println(pgErr.Code)
        }
    }

    return nil
}

// get a list of all users from the database
func (u userStore) GetAll() ([]*clients.User, error) {
    query := `
        SELECT id, first_name, last_name, phone_number, email, follow_up, check_in_date 
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
            &u.FollowUp,
            &u.CheckInDate,
        )
        if err != nil {
            return nil, err
        }

        users = append(users, &u)
    }

    return users, nil
}

// get user by ID from database
func (u userStore) GetUserByID(id int64) (*clients.User, error) {

    query := `SELECT first_name, last_name, phone_number, email FROM users WHERE id = $1`

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    var user clients.User

    err := u.db.QueryRow(ctx, query, id).Scan(
        &user.FirstName,
        &user.LastName,
        &user.PhoneNumber,
        &user.Email,
    )

    if err != nil {
        var pgErr *pgconn.PgError
        switch {
        case errors.Is(err, clients.ErrNoUsersFound):
            return nil, fmt.Errorf(pgErr.Message)
        default:
            return nil, err 
        }
    }

    return &user, nil
}

// update a user entry in the database
func (u userStore) UpdateUser(id int64, user clients.UserUpdate) error {

    query := `
        UPDATE users 
        SET first_name = COALESCE(@first_name, first_name),
            last_name = COALESCE(@last_name, last_name),
            phone_number = COALESCE(@phone_number, phone_number),
            email = COALESCE(@email, email)
        WHERE id = @id
    `

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    args := pgx.NamedArgs{
        "first_name":   &user.FirstName,
        "last_name":    &user.LastName,
        "phone_number": &user.PhoneNumber,
        "email":        &user.Email,
        "id":           &id,
    }


    _, err := u.db.Query(ctx, query, args)
    if err != nil {
        return fmt.Errorf("failed to query update client: %v", err)
    }

    return nil
}

// delete a user from the database
func (u userStore) DeleteUser(id int64) error {
    query := `DELETE FROM users WHERE id = $1`

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    payload, err := u.db.Exec(ctx, query, id)

    if err != nil {
        return fmt.Errorf("failed to delete from users: %v", err) 
    }

    if rows := payload.RowsAffected(); rows != 1 {
        return clients.ErrNoUsersFound
    }

    return nil
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




