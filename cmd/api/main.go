package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/agpelkey/clients"
	"github.com/agpelkey/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type config struct {
	port int
	env  string
}

type application struct {
	config     config
	logger     *Logger
	UsersStore clients.UserService
}

func main() {

	// create config
	var cfg config

	// Read in arguments for port and env.
	flag.IntVar(&cfg.port, "port", 8080, "API Server Port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// create logger
	logger := New(os.Stdout, LevelInfo)

	// db connection
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("CRM_DB_DSN"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool")
		os.Exit(1)
	}
	defer dbpool.Close()

    logger.PrintInfo("db connection established", nil)

	// create application
	app := &application{
		config: cfg,
		logger: logger,
        UsersStore: postgres.NewUserStore(dbpool),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.PrintInfo("starting server", map[string]string{
		"Addr": srv.Addr,
		"env":  app.config.env,
	})

	err = srv.ListenAndServe()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}
