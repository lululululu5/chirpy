package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/lululululu5/chirpy/database"
)

type apiConfig struct {
	fileserverHits int
	db *database.DB
}



func main() {
	const filepathRoot = "."
	const port = "8080"

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatal(err)
	}


	apiCfg := apiConfig{
		fileserverHits: 0,
		db: db, 
	}

	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerFileserverHitsCount)
	
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerPostValidation)	
	mux.HandleFunc("/api/reset", apiCfg.handlerFileserverHitsReset)
	

	srv := &http.Server{
		Addr: ":" + port,
		Handler: mux, 
	}

	

	log.Printf("Serving files from %s on port %s\n", filepathRoot,  port)
	log.Fatal(srv.ListenAndServe())
}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}


func (cfg *apiConfig) handlerFileserverHitsReset(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits = 0
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) handlerPostValidation(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}


	// Read data from Post
	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	//Response in case of Wrong length
	const maxChirpLenght = 140
	if len(params.Body) > maxChirpLenght {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// clean Body 
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	bodySlice := strings.Split(params.Body, " ")
	for i, word := range bodySlice {
		// Check if lower cased word is in badWords slice
		if slices.Contains(badWords, strings.ToLower(word)) {
			bodySlice = slices.Replace(bodySlice, i, i+1, "****")
		}
		}
	params.Body = strings.Join(bodySlice, " ")
	
	// Add Chirp to Database
	newChirp, err := cfg.db.CreateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "")
	}

	// Pass cleaned body to responseWithJSON formula
	respondWithJSON(w, http.StatusCreated, newChirp)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct{
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}


