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

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}
	
	user, err := cfg.db.GetUserByMail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User does not exist")
		return
	}

	// validatePassword
	err = auth.CheckHashPassword(user.HashedPassword, params.Password)
	fmt.Print(err)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Passwords or Email wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: user.ID,
		Email: user.Email,
	})
}