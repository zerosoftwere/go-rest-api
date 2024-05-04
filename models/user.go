package models

import (
	"errors"

	"alhassan.link/rest-api/db"
	"alhassan.link/rest-api/utils"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `binding:"required" json:"email"`
	Password string `binding:"required" json:"password"`
}

func (u *User) Save() error {
	query := `
	INSERT INTO users (email, password) VALUES (?, ?);
	`
	stmts, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmts.Close()

	password, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmts.Exec(u.Email, password)
	if err != nil {
		return err
	}
	userId, err := result.LastInsertId()
	u.ID = userId
	return err
}

func (u *User) ValidateCredentials() error {
	query := `
	SELECT id, password FROM users WHERE email=?
	`
	row := db.DB.QueryRow(query, u.Email)

	var retrievedPassword string
	err := row.Scan(&u.ID, &retrievedPassword)
	if err != nil {
		return nil
	}

	isPasswordValid := utils.CheckPasswordHash(u.Password, retrievedPassword)
	if !isPasswordValid {
		return errors.New("invalid credentials")
	}
	return nil
}
