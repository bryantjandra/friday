package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bryantjandra/friday/internal/store"
)

type UpdateReminderTool struct {
	db *sql.DB
}

type updateReminderArgs struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	DueAt string `json:"due_at"`
}

func NewUpdateReminderTool(db *sql.DB) *UpdateReminderTool {
	return &UpdateReminderTool{
		db: db,
	}
}

func (t *UpdateReminderTool) Name() string {
	return "update_reminder"
}

func (t *UpdateReminderTool) Description() string {
	return "Change an existing reminder's title, due date, or both. Use this when the user wants to reschedule, move, push back, bring forward, or rename a reminder (e.g. 'move the dentist to the 20th'). Never use create_reminder to change an existing reminder. Only send the fields that change. Get the id from list_reminders first; never guess an id."
}

func (t *UpdateReminderTool) InputSchema() json.RawMessage {
	return json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"id": {
						"type": "integer",
						"description": "The id of the reminder to be updated. You can retrieve this id through usage of list_reminders tool."
					},
					"title": {
						"type": "string",
						"description": "Optional. The new title. Omit to keep the current title."
					},
					"due_at": {
						"type": "string",
						"description": "Optional. The new due date and time, in RFC3339 format converted to UTC with a Z suffix, e.g. 2026-10-05T00:00:00Z (never use an offset like +08:00). Omit to keep the current due date."
					}
				},
				"required": ["id"]
			}
`)
}

func (t *UpdateReminderTool) Execute(ctx context.Context, args json.RawMessage) (result Result, err error) {
	var a updateReminderArgs

	if err := json.Unmarshal(args, &a); err != nil {
		return Result{Content: "invalid arguments", IsError: true}, nil
	}

	var dueAtParsed time.Time
	if a.DueAt != "" {
		dueAtParsed, err = time.Parse(time.RFC3339, a.DueAt)
		if err != nil {
			return Result{Content: "due_at must be RFC3339 format, e.g. 2026-10-05T00:00:00Z", IsError: true}, nil
		}
	}

	res, err := store.UpdateReminder(t.db, a.ID, a.Title, dueAtParsed)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return Result{Content: fmt.Sprintf("no reminder found with id: %d", a.ID), IsError: true}, nil
		}
		if errors.Is(err, store.ErrNothingToUpdate) {
			return Result{Content: fmt.Sprintf("nothing to update as no parameters were given, provide at least one of title or due_at"), IsError: true}, nil
		}
		return Result{}, err
	}

	return Result{Content: res, IsError: false}, nil
}
