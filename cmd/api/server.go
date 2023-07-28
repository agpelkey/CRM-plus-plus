package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)


func (app *application) server() error {
    // HTTP server setting
    srv := &http.Server{
        Addr: fmt.Sprintf(":%d", app.config.port),
        Handler: app.routes(),
        IdleTimeout: time.Minute,
        ReadTimeout: 10*time.Second,
        WriteTimeout: 30*time.Second,
    }

    // get spicy with a lil shutdownErr channel
    shutdownErr := make(chan error)

    // start background go routine
    go func() {
        // create a seperate quit channel to carry os signals
        quit := make(chan os.Signal, 1)

        // signal.Notify listens for incoming SIGIN and SIGTERM signals
        // and then relays them to the quit channel
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

        // read the signals from the channel into a variable
        s := <-quit

        // log a message to sayy the signal has been caught. 
        app.logger.PrintInfo("shutting down server", map[string]string {
            "signal": s.String(),
        })

        // context bb
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        // This is where the magic happens. Call Shutdown() on our sever, and relay the err message
        // to the shutdownErr channel
        shutdownErr <-srv.Shutdown(ctx)
    }()

    app.logger.PrintInfo("starting server", map[string]string {
        "add": srv.Addr,
        "env": app.config.env,
    })

    // by calling Shutdown() on our server ListenAndServe() will return http.ErrServerClosed. So, if we see this error,
    // it actually indicates that the graceful shutdown has started.
    err := srv.ListenAndServe()
    if !errors.Is(err, http.ErrServerClosed){
        return err
    }

    // otherwise, we wait for the channel to return the value from Shutdown()
    err = <-shutdownErr
    if err != nil {
        return err
    }

    app.logger.PrintInfo("stopped server", map[string]string {
        "addr": srv.Addr,
    })

    return nil
}
