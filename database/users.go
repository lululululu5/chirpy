package database

import (
	"errors"
)

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}

func (db *DB) CreateUser(email, hashedPassword string) (User, error){
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	exists, _ := checkForExistingUser(email, dbStructure.Users) 
	if exists {
		return User{}, errors.New("User with email already exists")
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID: id,
		Email: email,
		HashedPassword: hashedPassword,
		IsChirpyRed: false,
	}

	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil

}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[id]; if !ok {
		return User{}, errors.New("User does not exist")
	}

	return user, nil
} 

func (db *DB) GetUserByMail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	exists, id := checkForExistingUser(email, dbStructure.Users)
	if !exists {
		return User{}, errors.New("User does not exist")
	}
	
	user, ok := dbStructure.Users[id]; if !ok {
		return User{}, errors.New("User does not exist")
	}
	
	return user, nil
}

func (db *DB) UpdateUser(id int, hashedPassword, email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	oldUser, ok := dbStructure.Users[id]
	if !ok {
		return User{}, errors.New("User does not exist")
	}

	isChirpyRed := oldUser.IsChirpyRed

	user := User{
		ID: id,
		Email: email,
		HashedPassword: hashedPassword,
		IsChirpyRed: isChirpyRed,
	}

	dbStructure.Users[id] = user
	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func(db *DB) UpgradeUser(userID int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	
	user, ok := dbStructure.Users[userID]; if !ok {
		return errors.New("User does not exist")
	} else {
		user.IsChirpyRed = true
		dbStructure.Users[userID] = user
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}


func checkForExistingUser(lookupEmail string, lookupDB map[int]User) (bool, int) {
	for _,user:= range lookupDB{
	  if(user.Email == lookupEmail){
		return true, user.ID
	  }
	}
	return false, 0
  }

