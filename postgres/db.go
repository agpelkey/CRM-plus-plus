package main

import (

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
    DB *pgxpool.Pool
}

func New(dsn string) (Postgres, error) {
    db, err := Connect(dsn)
    
}
