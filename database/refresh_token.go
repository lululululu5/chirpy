package database

import (
	"errors"
	"time"
)

type RefreshToken struct {
	ID string `json:"id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserId int `json:"user_id"`
}

func (db *DB) StoreRefreshToken(token string, userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}


	expirationDays := 60*24 // days * hours as hours is the only available const in the time package
	refreshToken := RefreshToken{
		ID: token,
		RefreshToken: token,
		ExpiresAt: time.Now().UTC().Add(time.Duration(expirationDays)*time.Hour), // add duration of now + 60 days
		UserId: userId, 
	}

	dbStructure.RefreshTokens[token] = refreshToken
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) ValidateRefreshToken(token string) (int, bool, error)  {
	dbStructure, err := db.loadDB()
	if err != nil {
		return 0, false, err
	}
	
	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return 0, false, errors.New("could not locate refresh token")
	} else if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		return 0, false, errors.New("refresh token expired")
	}

	return refreshToken.UserId, true, nil
}

func (db *DB) RevokeRefreshToken(token string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	delete(dbStructure.RefreshTokens, token)
	err = db.writeDB(dbStructure)
	if err != nil {
		return false, err
	}

	return true, nil
}