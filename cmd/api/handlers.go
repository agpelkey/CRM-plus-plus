package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/agpelkey/clients"
	"github.com/julienschmidt/httprouter"
)

// Method to handle POST requests and create new users
func (app *application) handleCreateUser(w http.ResponseWriter, r *http.Request) {

    input := clients.UserCreate{}

	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

    err = input.Validate()
    if err != nil {
        app.badRequestResponse(w, r, err)    
        return
    }


    user := input.CreateModel()

	err = app.UsersStore.Create(&user)
	if err != nil {
        app.serverErrorResponse(w, r, err)
        return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/user/%d", user.ID))

	err = writeJSON(w, http.StatusCreated, envelope{"user": user}, headers)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

}

// Method to get all users from database
func (app *application) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
    users, err := app.UsersStore.GetAll()
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    _ = writeJSON(w, http.StatusOK, envelope{"users": users}, nil)
}

// Method to get user by ID
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

    err = writeJSON(w, http.StatusOK, envelope{"users":user}, nil)
        
}

// handleUpdateUser handles the "PATCH /v1/users/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (app *application) handleUpdateUser(w http.ResponseWriter, r *http.Request) {

    id, err := app.readIDParam(r)
    if err != nil {
        app.notFoundResponse(w, r)
        return
    }
    
    input := clients.UserUpdate{}
    err = readJSON(w, r, &input)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }


    err = app.UsersStore.UpdateUser(id, input)
    if err != nil {
       app.serverErrorResponse(w, r, err) 
       return
    }

    err = writeJSON(w, http.StatusOK, envelope{"user": err}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return 
    }
}

// Method to handle DELETE requests from client
func (app *application) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
  
    id, err := app.readIDParam(r)
    if err != nil {
        app.serverErrorResponse(w, r, err)
        return
    }

    err = app.UsersStore.DeleteUser(id)
    if err != nil {
       switch {
       case errors.Is(err, clients.ErrRecordNotFound):
           app.notFoundResponse(w, r)
        default:
            app.serverErrorResponse(w, r, err)
       } 
    }

    err = writeJSON(w, http.StatusOK, envelope{"message": "user successfully deleted"}, nil)
    if err != nil {
        app.serverErrorResponse(w, r, err)
    }

}

// still playing around with this
/*
func (app *application) handleGetUserByLastName(w http.ResponseWriter, r *http.Request) {
}
*/

// function to get ID from url 
func (app *application) readIDParam(r *http.Request) (int64, error) {
    params := httprouter.ParamsFromContext(r.Context())

    id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
    if err != nil {
        return 0, errors.New("invalid id parameter")
    }

    return id, nil
}



