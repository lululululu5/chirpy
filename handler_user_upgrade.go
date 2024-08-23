package main

import (
	"encoding/json"
	"net/http"

	"github.com/lululululu5/chirpy/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Authorization header wrong")
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Polka API key is invalid")
		return 
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	decoder.Decode(&params)

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNotFound, "User event does not exist")
		return
	}

	// call db method to upgrade user status
	err = cfg.db.UpgradeUser(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not upgrade user")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")


}