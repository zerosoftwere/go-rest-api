package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "app.db")

	if err != nil {
		panic("could not connect to databaes.")
	}

	DB.SetMaxOpenConns(5)
	DB.SetMaxIdleConns(2)

	createTables()
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		panic("Could not close connection to database")
	}
}

func createTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(128) UNIQUE NOT NULL,
		password VARCHAR(64) NOT NULL
	);
	`

	_, err := DB.Exec((createUserTable))
	if err != nil {
		panic("could not create users table")
	}

	createEventTable := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		dateTime DATETIME NOT NULL,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	`
	_, err = DB.Exec(createEventTable)
	if err != nil {
		panic("could not create events tables")
	}

	createRegistrationTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		event_id INTEGER NOT NULL REFERENCES events(id) ON DELETE CASCADE,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		UNIQUE(user_id, event_id)
	);
	`
	_, err = DB.Exec(createRegistrationTable)
	if err != nil {
		fmt.Println(err)
		panic("could not create registrations table")
	}
}
