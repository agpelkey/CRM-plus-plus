package postgres

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres connection pool
type Postgres struct {
	DB *pgxpool.Pool
}

// New returns a new instance of Postgres
func New(dsn string) (Postgres, error) {

	db, err := Connect(dsn)
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{DB: db}, nil
}

// Connect connects to Postgres Driver with given DSN
func Connect(dsn string) (*pgxpool.Pool, error) {
	var (
		once sync.Once
		err  error
		db   *pgxpool.Pool
	)

	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		db, err = pgxpool.New(ctx, dsn)
		err = db.Ping(ctx)
	})

	if err != nil || db == nil {
		return nil, err
	}

	return db, nil
}
