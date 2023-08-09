package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/agpelkey/clients"
	"github.com/julienschmidt/httprouter"
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

func (app *application) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
    
    id, err := app.readIDParam(r)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    user, err := app.UsersStore.GetUserByID(id)
    if err != nil {
        switch {
        case errors.Is(err, clients.ErrRecordNotFound):
            app.notFoundResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
        }
        return
    }

    headers := make(http.Header)
    headers.Set("Location", fmt.Sprintf("v1/users/%d", user.ID))

    err = writeJSON(w, http.StatusOK, envelope{"users":user}, headers)
        
}

func (app *application) handleGetUserByLastName(w http.ResponseWriter, r *http.Request) {
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil {
        return 0, errors.New("invalid id parameter")
    }

    return id, nil
}



