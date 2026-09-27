package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bryantjandra/friday/internal/store"
)

type DeleteReminderTool struct {
	db *sql.DB
}

type deleteReminderArgs struct {
	ID int64 `json:"id"`
}

func NewDeleteReminderTool(db *sql.DB) *DeleteReminderTool {
	return &DeleteReminderTool{
		db: db,
	}
}

func (t *DeleteReminderTool) Name() string {
	return "delete_reminder"
}

func (t *DeleteReminderTool) Description() string {
	return "Permanently delete a reminder. This cannot be undone. Use this only when the user explicitly asks to delete, remove, or cancel a reminder; if they have done the task, use complete_reminder instead. Get the id from list_reminders first; never guess an id."
}

func (t *DeleteReminderTool) InputSchema() json.RawMessage {
	return json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"id": {
						"type": "integer",
						"description": "The id of the reminder to be deleted. You can retrieve this id through usage of list_reminders tool."
					}
				},
				"required": ["id"]
			}
`)
}

func (t *DeleteReminderTool) Execute(ctx context.Context, args json.RawMessage) (result Result, err error) {
	var a deleteReminderArgs

	if err := json.Unmarshal(args, &a); err != nil {
		return Result{Content: "invalid arguments", IsError: true}, nil
	}

	res, err := store.DeleteReminder(t.db, a.ID)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return Result{Content: fmt.Sprintf("no reminder found with id: %d", a.ID), IsError: true}, nil
		}
		return Result{}, err
	}

	return Result{Content: res, IsError: false}, nil
}
