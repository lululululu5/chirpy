package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lululululu5/chirpy/auth"
)

func(cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}
	
	user, err := cfg.db.GetUserByMail(params.Email) // returns User Struct with hashed password
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User does not exist")
		return
	}


	fmt.Printf("This is the provided Password: %s\n", params.Password)
	fmt.Printf("This is the stored hash: %s\n", user.HashedPassword)
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	fmt.Println(err)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Passwords or Email wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
		},
	})
}