package database

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
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

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mu: &sync.RWMutex{},
	}

	err := db.ensureDB()
	return db, err 
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dBStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dBStructure.Chirps) + 1
	chirp := Chirp{
		ID: id, 
		Body: body,
	}
	dBStructure.Chirps[id] = chirp

	err = db.writeDB(dBStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	Id, err := strconv.Atoi(id)
	if err != nil {
		return Chirp{}, err
	}

	val, ok := dbStructure.Chirps[Id]; if !ok {
		return Chirp{}, errors.New("Chirp Id does not exist")
	}

	return val, nil

}

func (db *DB) CreateUser(email string) (User, error){
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
	}

	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil

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