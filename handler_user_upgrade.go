package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserId int `json:"user_id"`
		} `json:"data"`
	}


	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	decoder.Decode(&params)

	if params.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "User event does not exist")
		return
	}

	// call db method to upgrade user status
	err := cfg.db.UpgradeUser(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not upgrade user")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")


}