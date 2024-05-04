package models

import (
	"log"
	"time"

	"alhassan.link/rest-api/db"
)

var logger = log.Default()

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `binding:"required" json:"name"`
	Description string    `binding:"required" json:"description"`
	Location    string    `binding:"required" json:"location"`
	DateTime    time.Time `binding:"required" json:"datetime"`
	UserID      int64     `json:"user_id"`
}

func New() Event {
	return Event{}
}

func (e *Event) Save() error {
	query := `
	INSERT INTO events (name, description, location, dateTime, user_id) VALUES(?, ?, ?, ?, ?)
	`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)

	if err != nil {
		return err
	}

	e.ID, err = result.LastInsertId()
	return err
}

func GetAllEvents() ([]Event, error) {
	query := `
	SELECT * FROM events
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	return events, nil
}

func GetEvent(eventId int64) (*Event, error) {
	var event Event

	query := `
	SELECT * FROM events WHERE id=?
	`
	row := db.DB.QueryRow(query, eventId)

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func UpdateEvent(eventId int64, event *Event) (*Event, error) {
	query := `
	UPDATE events SET name=?, description=?, location=?, datetime=? WHERE id=?
	`

	oldEvent, err := GetEvent(eventId)
	if err != nil {
		logger.Print(err)
		return nil, err
	}

	if oldEvent == nil {
		return nil, nil
	}

	_, err = db.DB.Exec(query, event.Name, event.Description, event.Location, event.DateTime, eventId)
	if err != nil {
		logger.Print(err)
		return nil, err
	}

	return &Event{
		ID:          eventId,
		Name:        event.Name,
		Description: event.Description,
		Location:    event.Location,
		DateTime:    event.DateTime,
		UserID:      oldEvent.UserID,
	}, nil
}

func DeleteEvent(eventId int64) (bool, error) {
	query := `
	DELETE FROM events WHERE id=?1
	`
	result, err := db.DB.Exec(query, eventId)
	if err != nil {
		return false, err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return affectedRows > 0, nil
}
