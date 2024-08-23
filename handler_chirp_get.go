package main

import (
	"net/http"
	"sort"
)


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
	authId := req.URL.Query().Get("author_id")
	dbChirps, err := cfg.db.GetChirps(authId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not load chirps")
		return 
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID: dbChirp.ID,
			Body: dbChirp.Body,
			AuthorID: dbChirp.AuthorID,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func(cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("chirpID")
	dbChirp, err := cfg.db.GetChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not retrieve Chirp with ID")
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID: dbChirp.ID,
		Body: dbChirp.Body,
		AuthorID: dbChirp.AuthorID,
	})

}