package database

import (
	"errors"
	"strconv"
)

type Chirp struct {
	ID int `json:"id"`
	Body string `json:"body"`
	AuthorID int `json:"author_id"`
}


func (db *DB) CreateChirp(body string, userID int) (Chirp, error) {
	dBStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dBStructure.Chirps) + 1
	chirp := Chirp{
		ID: id, 
		Body: body,
		AuthorID: userID,
	}
	dBStructure.Chirps[id] = chirp

	err = db.writeDB(dBStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps(authorIDString string) ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	authorID, err := strconv.Atoi(authorIDString)
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		// if chirp belongs to authorId append
		if authorID != 0 {
			if authorID == chirp.AuthorID {
				chirps = append(chirps, chirp)		
			} else {
				continue
			}
		} else {
			chirps = append(chirps, chirp)
		}
		
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

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	delete(dbStructure.Chirps, id)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	
	return nil
}