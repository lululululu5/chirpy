package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lululululu5/chirpy/auth"
)

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	subject, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.GenerateHashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userIDInt, err := strconv.Atoi(subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
	}

	user, err := cfg.db.UpdateUser(userIDInt, hashedPassword, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user data")
		return
	}

	// Response if code was executed successfully
	respondWithJSON(w, http.StatusOK, response{
		User{
			ID: user.ID,
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})

}