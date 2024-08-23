package main

import (
	"encoding/json"
	"net/http"

	"github.com/lululululu5/chirpy/auth"
)

func(cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
	}
	type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Passwords or Email wrong")
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generte JWT token")
		return 
	}

	refreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate Refresh Token")
		return
	}

	err = cfg.db.StoreRefreshToken(refreshToken, user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not store Refresh Token")
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token: token,
		RefreshToken: refreshToken,
	})
}