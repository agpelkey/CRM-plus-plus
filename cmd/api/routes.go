package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

    router.HandlerFunc(http.MethodPost, "/v1/users", app.handleCreateUser)
    router.HandlerFunc(http.MethodGet, "/v1/users", app.handleGetAllUsers)
    router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.handleGetUserByID)
    router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.handleDeleteUser)
    router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.handleUpdateUser)

	return router

}
