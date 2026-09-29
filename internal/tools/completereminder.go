package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bryantjandra/friday/internal/store"
)

type CompleteReminderTool struct {
	db *sql.DB
}

type completeReminderArgs struct {
	ID int64 `json:"id"`
}

func NewCompleteReminderTool(db *sql.DB) *CompleteReminderTool {
	return &CompleteReminderTool{
		db: db,
	}
}

func (t *CompleteReminderTool) Name() string {
	return "complete_reminder"
}

func (t *CompleteReminderTool) Description() string {
	return "Mark a reminder as done. Use this when the user says they have done, finished, or completed something, even if they don't mention a reminder (e.g. 'I paid rent', 'done with the dentist'). Prefer this over delete_reminder unless the user explicitly asks to remove it. Get the id from list_reminders first; never guess an id. If two or more reminders match the user's words (e.g. two titles containing 'dentist'), do NOT call this tool: list the matching reminders and ask the user which one they mean."
}

func (t *CompleteReminderTool) InputSchema() json.RawMessage {
	return json.RawMessage(`
			{
				"type": "object",
				"properties": {
					"id": {
						"type": "integer",
						"description": "The id of the reminder to be marked as completed. You can retrieve this id through usage of list_reminders tool."
					}
				},
				"required": ["id"]
			}
`)
}

func (t *CompleteReminderTool) Execute(ctx context.Context, args json.RawMessage) (result Result, err error) {
	var a completeReminderArgs

	if err := json.Unmarshal(args, &a); err != nil {
		return Result{Content: "invalid arguments", IsError: true}, nil
	}

	res, err := store.CompleteReminder(t.db, a.ID)

	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return Result{Content: fmt.Sprintf("no reminder found with id: %d", a.ID), IsError: true}, nil
		}
		return Result{}, err
	}

	return Result{Content: res, IsError: false}, nil
}
