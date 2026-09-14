package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bryantjandra/friday/internal/store"
)

type CreateReminderTool struct {
	db *sql.DB
}

func NewCreateReminderTool(db *sql.DB) *CreateReminderTool {
	return &CreateReminderTool{
		db: db,
	}
}

type createReminderArgs struct {
	Title string `json:"title"`
	DueAt string `json:"due_at"`
}

func (t *CreateReminderTool) Name() string {
	return "create_reminder"
}

func (t *CreateReminderTool) Description() string {
	return "Create a reminder for the user at a specific date. Use it when the user wants to create a reminder"
}

func (t *CreateReminderTool) InputSchema() json.RawMessage {
	return json.RawMessage(`
		{
			"type": "object",
			"properties": {
				"title": {
					"type": "string",
					"description": "What to be reminded about."
				},
				"due_at": {
					"type": "string",
					"description": "The date of the reminder in ISO 8601 format, e.g. 2026-10-05T00:00:00Z."
				}
			},
			"required": ["title", "due_at"]
		}
	`)
}

func (t *CreateReminderTool) Execute(ctx context.Context, args json.RawMessage) (result Result, err error) {
	var a createReminderArgs

	/*
		There are two ways how we are returning errors here:
			1. Inside the Result struct --> an issue with the model's response and the model can then try again.
			2. Inside the standard error --> a program failure from our code (model cannot fix by trying again).
	*/

	if err := json.Unmarshal(args, &a); err != nil {
		return Result{IsError: true, Content: "invalid arguments"}, nil
	}

	if a.Title == "" {
		return Result{IsError: true, Content: "title is required"}, nil
	}

	dueAtParsed, err := time.Parse(time.RFC3339, a.DueAt)
	if err != nil {
		return Result{IsError: true, Content: "due_at must be RFC3339, e.g. 2026-10-05T00:00:00Z"}, nil
	}

	id, err := store.CreateReminder(t.db, a.Title, dueAtParsed, "")
	if err != nil {
		return Result{}, err
	}

	return Result{IsError: false, Content: fmt.Sprintf("Reminder created with id %d", id)}, nil
}
