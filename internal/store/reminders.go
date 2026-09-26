package store

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrNotFound = errors.New("reminder not found")

type Reminder struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	DueAt     string `json:"due_at"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

func CreateReminder(db *sql.DB, title string, dueAt time.Time) (id int64, err error) {

	/* Parameterised queries are used to avoid SQL injection attacks (the SQL command and the value are passed to the DB seperately, it plugs in the values as pure data). */
	query := "INSERT INTO reminders (title, due_at, created_at) VALUES (?, ?, ?);"

	/* sqlite stores eerything as strings (there is no date/time type). Thus, the RFC3339 format is needed to ensure the dates can be chronologically ordered, even though it's a string. */
	result, err := db.Exec(query, title, dueAt.UTC().Format(time.RFC3339), time.Now().UTC().Format(time.RFC3339))

	if err != nil {
		return 0, err
	}

	id64, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id64, nil

}

func CompleteReminder(db *sql.DB, id int64) (res string, err error) {
	query := "UPDATE reminders SET completed=1 WHERE id = ?"
	result, err := db.Exec(query, id)

	if err != nil {
		return "", err
	}

	/* we need to check how many rows are affected, because it actually doesn't return as an error if we pass in an id that doesn't exist in our DB */
	/* that is why we need to make sure the number of rows affected is actually > 0, to prove that our function actually works as intended */
	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows == 0 {
		return "", ErrNotFound
	}

	return fmt.Sprintf("Successfully completed reminder %d", id), nil
}

func DeleteReminder(db *sql.DB, id int64) (res string, err error) {
	query := "DELETE FROM reminders WHERE id=?"
	result, err := db.Exec(query, id)

	if err != nil {
		return "", err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows == 0 {
		return "", ErrNotFound
	}

	return fmt.Sprintf("Successfully deleted reminder %d", id), nil
}

func ListReminders(db *sql.DB, from time.Time, to time.Time, includeCompleted bool) (reminders []Reminder, err error) {
	var conditions []string
	var args []any
	var reminderList []Reminder

	if !from.IsZero() {
		conditions = append(conditions, "due_at >= ?")
		args = append(args, from.UTC().Format(time.RFC3339))
	}

	if !to.IsZero() {
		conditions = append(conditions, "due_at <= ?")
		args = append(args, to.UTC().Format(time.RFC3339))
	}

	if !includeCompleted {
		conditions = append(conditions, "completed=0")
	}

	query := "SELECT id, title, due_at, completed, created_at FROM reminders"

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY due_at ASC"
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r Reminder
		if err := rows.Scan(&r.ID, &r.Title, &r.DueAt, &r.Completed, &r.CreatedAt); err != nil {
			return nil, err
		}
		reminderList = append(reminderList, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminderList, nil
}

func UpdateReminder(db *sql.DB, id int64, title string, dueAt time.Time) (res string, err error) {
	var conditions []string
	var args []any

	if title != "" {
		conditions = append(conditions, "title=?")
		args = append(args, title)
	}

	if !dueAt.IsZero() {
		conditions = append(conditions, "due_at=?")
		args = append(args, dueAt.UTC().Format(time.RFC3339))
	}

	query := "UPDATE reminders"
	args = append(args, id)

	if len(conditions) > 0 {
		query += " SET " + strings.Join(conditions, ",")
	}

	if len(conditions) == 0 {
		return "", errors.New("nothing to update as no parameters were supplied")
	}

	query += " WHERE id = ?"

	result, err := db.Exec(query, args...)
	if err != nil {
		return "", err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}

	if rows == 0 {
		return "", ErrNotFound
	}

	return fmt.Sprintf("Successfully updated reminder %d", id), nil

}
