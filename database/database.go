package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	chirpsCount int
	mux *sync.RWMutex
}

type Chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(path, []byte{}, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	db := &DB{
		path: path,
		chirpsCount: 0,
		mux: new(sync.RWMutex),

	}
	return db, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	newChirp := Chirp{
		Id: db.chirpsCount +1,
		Body: body,
	}
	db.chirpsCount++
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}
	dbStructure.Chirps[db.chirpsCount] = newChirp
	db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
	}
	return newChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}

	Chirps := []Chirp{}
	for _, chirp := range dbStructure.Chirps {
		Chirps = append(Chirps, chirp)
	}
	// sorting of a slice based on the Id field
	sort.Slice(Chirps, func(i,j int) bool {
		return Chirps[i].Id < Chirps[j].Id
	})
	return Chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error { 
	if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile(db.path, []byte{}, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	// should fall back to ensureDB, to make sure it does not cause any problems
	err := db.ensureDB()
	if err != nil {
		return DBStructure{}, err
	}

	rawDat, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(rawDat, &dbStructure)
	if err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}


func (db *DB) writeDB(dbStructure DBStructure) error {
	err := db.ensureDB()
	if err != nil {
		return err
	}

	dat, err := json.Marshal(dbStructure) 
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0666)
	if err != nil {
		return err
	}

	return nil
}
