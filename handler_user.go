package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lululululu5/chirpy/auth"
)

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"-"`
}

// type CustomClaims struct {
// 	jwt.RegisteredClaims
// 	Subject string
// }

func(cfg  *apiConfig) handlerCreateUsers(w http.ResponseWriter, req *http.Request) {
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	hashedPassword, err := auth.GenerateHashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password")
	}

	user, err := cfg.db.CreateUser(params.Email, hashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create new user")
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
} 

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, req *http.Request) {
	rawToken := strings.TrimSpace(strings.ReplaceAll(req.Header.Get("Authorization"), "Bearer ", "")) // Extract value from different Headers
	token, err := jwt.ParseWithClaims(
		rawToken,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error in decoding token")
		return
	}


	if !token.Valid{
		// Checks for a range of different validations incl. expiration etc.
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	claims := token.Claims.(*jwt.RegisteredClaims)
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not convert id to string")
		return
	}

	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
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

	user, err := cfg.db.UpdateUser(id, hashedPassword, params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update user data")
		return
	}

	// Response if code was executed successfully
	respondWithJSON(w, http.StatusOK, response{
		User{
			ID: user.ID,
			Email: user.Email,
		},
	})

}