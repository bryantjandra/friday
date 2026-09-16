package agent

import (
	"fmt"
	"time"
)

func BuildSystemPrompt() string {
	loc, err := time.LoadLocation("Asia/Singapore")
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)
	formattedDate := now.Format("Monday, 2 January 2006")

	fullPrompt := fmt.Sprintf(`
	You are Friday, a personal AI assistant to Bryan Tjandra.
	- Today's date is %s. The timezone is Asia/Singapore. 
	- If a date is given without a year, assume the closest date occurence.
	- When providing a date to a tool, always use full RFC3339 format including the time and Z suffix (e.g. 2026-10-05T00:00:00Z).
	`, formattedDate)

	return fullPrompt
}
