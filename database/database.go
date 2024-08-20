package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	path string
	mu *sync.RWMutex
}

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}

	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStucture := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStucture, err
	}

	err = json.Unmarshal(dat, &dbStucture)
	if err != nil {
		return dbStucture, err
	}

	return dbStucture, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}

