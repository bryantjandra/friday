package store

import (
	"database/sql"
	"time"
)

func CreateReminder(db *sql.DB, title string, dueAt time.Time, notes string) (id int64, err error) {

	/* Parameterised queries are used to avoid SQL injection attacks (the SQL command and the value are passed to the DB seperately, it plugs in the values as pure data). */
	query := "INSERT INTO reminders (title, due_at, notes, created_at) VALUES (?, ?, ?, ?);"

	/* sqlite stores eerything as strings (there is no date/time type). Thus, the RFC3339 format is needed to ensure the dates can be chronologically ordered, even though it's a string. */
	result, err := db.Exec(query, title, dueAt.UTC().Format(time.RFC3339), notes, time.Now().UTC().Format(time.RFC3339))

	if err != nil {
		return 0, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id64, nil

}
