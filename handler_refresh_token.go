package main

import (
	"net/http"

	"github.com/lululululu5/chirpy/auth"
)

func(cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Decode req Header => Auth is send in header
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Refresh Token")
	}

	// validate refreshToken
	userID, valid, err := cfg.db.ValidateRefreshToken(refreshToken)
	if err != nil || !valid {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate refresh token")
		return
	}

	// Generate new token based on refresh token
	expiresOneHour := 60 * 60 
	accessToken, err := auth.MakeJWT(userID, cfg.jwtSecret, expiresOneHour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generate new access token")
		return 
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func(cfg *apiConfig) handlerRevokeRefresh(w http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get Refresh Token")
	}

	// validate refreshToken
	_, valid, err := cfg.db.ValidateRefreshToken(refreshToken)
	if err != nil || !valid {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate refresh token")
		return
	}

	// delete refreshToken
	isRevoked, err:=  cfg.db.RevokeRefreshToken(refreshToken)
	if err != nil || !isRevoked{
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}