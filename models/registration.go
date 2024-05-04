package models

import (
	"fmt"

	"alhassan.link/rest-api/db"
)

type Registration struct {
	ID      int64
	EventID int64
	UserID  int64
}

func (r *Registration) Save() error {
	query := `
	INSERT INTO registrations (user_id, event_id) VALUES(?, ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(r.UserID, r.EventID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = id
	return nil
}

func GetEventRegistration(userId, eventId int64) (*Registration, error) {
	query := `
	SELECT id, user_id, event_id FROM registrations WHERE user_id=? AND event_id=?
	`
	row := db.DB.QueryRow(query, userId, eventId)
	var registration Registration
	err := row.Scan(&registration.ID, &registration.UserID, &registration.EventID)
	return &registration, err
}

func HasEventRegistration(userId, eventId int64) bool {
	query := `
	SELECT COUNT(*) FROM registrations WHERE user_id=? AND event_id=?
	`
	row := db.DB.QueryRow(query, userId, eventId)

	var count int64
	err := row.Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func GetEventRegistrations(userId int64) ([]Registration, error) {
	query := `
	SELECT id, user_id, event_id FROM registrations WHERE user_id=?
	`
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	registrations := []Registration{}
	for rows.Next() {
		var registration Registration
		err = rows.Scan(&registration.ID, &registration.UserID, &registration.EventID)

		if err != nil {
			return nil, err
		}

		registrations = append(registrations, registration)
	}
	return registrations, nil
}

func DeleteEventRegistration(userId, eventId int64) error {
	query := `
	DELETE FROM registrations WHERE user_id=? AND event_id=?
	`
	_, err := db.DB.Exec(query, userId, eventId)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
