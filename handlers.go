package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// create variable to read request values into
	var input struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}

	// use our readJSON function to decode the request body into the input struct.
	// If this returns an error we send the client the error message along with
	// a 400 Bad request status code.
	err := readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// copy the values from the input struct to a new user struct
	user := &User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		PhoneNumber: input.PhoneNumber,
		Email:       input.Email,
	}

	// call db insert method
	err = app.DB.CreateUser(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// send back http response with location header to let client know which URL they can find the newly
	// created resource at.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/clients/%d", user.ID))

	// write a JSON response with a 201 created status code, the user data, and the location header
	err = writeJSON(w, http.StatusCreated, envelope{"client info:": user}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) handleGetUser(id int) {

}
