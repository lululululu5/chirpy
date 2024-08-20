package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

func(cfg  *apiConfig) handlerCreateUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.db.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new user")
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID: user.ID,
		Email: user.Email,
	})

} 