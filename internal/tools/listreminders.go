package tools

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bryantjandra/friday/internal/store"
)

type ListReminderTool struct {
	db       *sql.DB
	location *time.Location
}

type listReminderArgs struct {
	From             string `json:"from"`
	To               string `json:"to"`
	IncludeCompleted bool   `json:"include_completed"`
}

func NewListRemindersTool(db *sql.DB, location *time.Location) *ListReminderTool {
	return &ListReminderTool{
		db:       db,
		location: location,
	}
}

func (l *ListReminderTool) Name() string {
	return "list_reminders"
}

func (l *ListReminderTool) Description() string {
	return "List the user's reminders, sorted by due date. Use this to answer questions like 'what's coming up?', and always call it first to find a reminder's id before using complete_reminder, update_reminder, or delete_reminder. When looking for a specific reminder by name, omit from and to so a date filter doesn't hide it. Returns a JSON array of reminders with their ids; an empty result means nothing matched. The result times are in the user's local time (e.g. 2026-10-02T18:00:00+08:00)"
}

func (l *ListReminderTool) InputSchema() json.RawMessage {
	return json.RawMessage(`
		{
			"type": "object",
			"properties": {
				"from": {
					"type": "string",
					"description": "Optional, RFC3339 format, it is an inclusive lower bound, omit for no lower bound. Always convert dates to UTC with a Z suffix when passing them to tools (never use a timezone offset like +08:00). "
				},
				"to": {
					"type": "string",
					"description": "Optional, RFC3339 format, inclusive upper bound, omit for no upper bound. Always convert dates to UTC with a Z suffix when passing them to tools (never use a timezone offset like +08:00)."
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

	for i := range reminderList {
		dueAt, err := time.Parse(time.RFC3339, reminderList[i].DueAt)
		if err != nil {
			return Result{}, fmt.Errorf("reminder %d has invalid due_at, %q: %w", reminderList[i].ID, dueAt, err)
		}
		reminderList[i].DueAt = dueAt.In(l.location).Format(time.RFC3339)
	}

	/* data is json text, but in array of bytes format */
	data, err := json.Marshal(reminderList)
	if err != nil {
		return Result{}, err
	}

	/* content is json text, but in string format */
	return Result{IsError: false, Content: string(data)}, nil

}
