package main

import (
	"encoding/json"
	"net/http"
	"time"

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

	// Logic to ensure default and max of 24h expiration time for token generation
	defaultExpiration := 60*60*24
	if params.ExpiresInSeconds > defaultExpiration || params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	}

	// token := jwt.NewWithClaims(
	// 	jwt.SigningMethodHS256,
	// 	jwt.RegisteredClaims{
	// 		Issuer: "chirpy",
	// 		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
	// 		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(params.ExpiresInSeconds))),
	// 		Subject: strconv.Itoa(user.ID),
	// 	},
	// 	)
	// tokenString, err := token.SignedString([]byte(cfg.jwtSecret))
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Could not generte JWT token")
	// 	return 
	// }

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not generte JWT token")
		return 
	}


	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
		},
		Token: token,
	})
}