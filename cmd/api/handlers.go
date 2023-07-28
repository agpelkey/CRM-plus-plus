package main

import (
	"fmt"
	"net/http"

	"github.com/agpelkey/clients"
)

func (app *application) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &clients.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
	}

	err = app.UsersStore.Create(user)
	if err != nil {
        app.serverErrorResponse(w, r, err)
        return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/user/%d", user.ID))

	err = writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)

}

func (app *application) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := app.UsersStore.GetAll()
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    _ = writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
}

func (app *application) handleGetUserByLastName(w http.ResponseWriter, r *http.Request) {
    filter := clients.UserFilter{}    

    users, err := app.UsersStore.List(r.Context(), filter)
}



