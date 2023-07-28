package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

    router.HandlerFunc(http.MethodPost, "/v1/users", app.handleCreateUser)
    router.HandlerFunc(http.MethodGet, "/v1/users", app.handleGetAllUsers)

	return router

}
