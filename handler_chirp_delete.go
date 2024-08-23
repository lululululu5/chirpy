package main

import (
	"net/http"
	"strconv"

	"github.com/lululululu5/chirpy/auth"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) {
	// Extract token from header
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not extract bearer token")
		return 
	}
	// Extract user_id from token
	currentUserIDString, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User authentication failed")
		return 
	}

	currentUserID, err := strconv.Atoi(currentUserIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not convert string to int")
		return
	}

	// Extract chirp id from header
	chirpID := req.PathValue("chirpID")
	chirp, err := cfg.db.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp does not exist")
		return 
	}
	// Check whether userId(author_id) from token and userid from chirp are the same
	if chirp.AuthorID != currentUserID {
		respondWithError(w, http.StatusForbidden, "User not allowed to delete chirp")
		return 
	}
	// If the same delete chirp from DB. => Create DB method for deletion
	err = cfg.db.DeleteChirp(chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete chirp")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")

}