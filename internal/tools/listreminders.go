package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/bryantjandra/friday/internal/store"
)

type ListReminderTool struct {
	db *sql.DB
}

type listReminderArgs struct {
	From             string `json:"from"`
	To               string `json:"to"`
	IncludeCompleted bool   `json:"include_completed"`
}

func NewListRemindersTool(db *sql.DB) *ListReminderTool {
	return &ListReminderTool{
		db: db,
	}
}

func (l *ListReminderTool) Name() string {
	return "list_reminders"
}

func (l *ListReminderTool) Description() string {
	return "List reminders. Takes in optional arguments (`from`, `to`, `include_completed`)."
}

func (l *ListReminderTool) InputSchema() json.RawMessage {
	return json.RawMessage(`
		{
			"type": "object",
			"properties": {
				"from": {
					"type": "string",
					"description": "Optional, RFC3339 format, it is an inclusive lower bound, omit for no lower bound."
				},
				"to": {
					"type": "string",
					"description": "Optional, RFC3339 format, inclusive upper bound, omit for no upper bound."
				},
				"include_completed": {
					"type": "boolean",
					"description": "Optional, defaults to false, whether to include complete reminders or not."
				}
			}
		}
	`)
}

func (l *ListReminderTool) Execute(ctx context.Context, args json.RawMessage) (Result, error) {
	var a listReminderArgs
	if err := json.Unmarshal(args, &a); err != nil {
		return Result{IsError: true, Content: "invalid arguments"}, nil
	}

	var fromParsed time.Time
	var err error
	if a.From != "" {
		fromParsed, err = time.Parse(time.RFC3339, a.From)
		if err != nil {
			return Result{IsError: true, Content: "from must be RFC3339 format, e.g. e.g. 2026-10-05T00:00:00Z"}, nil
		}
	}

	var toParsed time.Time
	if a.To != "" {
		toParsed, err = time.Parse(time.RFC3339, a.To)
		if err != nil {
			return Result{IsError: true, Content: "to must be RFC3339 format, e.g. 2026-10-05T00:00:00Z"}, nil
		}
	}

	reminderList, err := store.ListReminders(l.db, fromParsed, toParsed, a.IncludeCompleted)

	if err != nil {
		return Result{}, err
	}

	/* data is json text, but in array of bytes format */
	data, err := json.Marshal(reminderList)
	if err != nil {
		return Result{}, err
	}

	/* content is json text, but in string format */
	return Result{IsError: false, Content: string(data)}, nil

}
