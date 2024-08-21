package database

import (
	"errors"
	"strconv"
)



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