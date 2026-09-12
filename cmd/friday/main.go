package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bryantjandra/friday/internal/config"
	"github.com/bryantjandra/friday/internal/llm"
	"github.com/bryantjandra/friday/internal/store"
)

/* Go forces main to take no parameters, and no return value. Thus, we use a seperate run() function to know what error was produced */
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return nil
	}

	switch os.Args[1] {
	case "version":
		fmt.Println("friday version")
		return nil
	default:
		/* Not a reserved subcommand, thus treat the whole string as the prompt */
		return runPrompt(os.Args[1])
	}
}

func runPrompt(prompt string) error {

	/* Step 1: Load the config values */
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	/* Step 2: Open the DB connection */
	db, err := store.Open(cfg.DBPath)
	if err != nil {
		return err
	}

	/*
		When we defer something, we ensure it always runs whenever this function (runPrompt) exits, no matter what.
		We ensure we close our DB when this function returns to ensure we don't waste resources.
	*/
	defer db.Close()

	/* Step 3: Run the schema */
	if err := store.Migrate(db); err != nil {
		return err
	}

	client := llm.NewClient(cfg)

	req := llm.Request{
		MaxTokens: 1024,
		Messages: []llm.Message{
			{
				Role: "user",
				Content: []llm.ContentBlock{
					{Type: "text", Text: prompt},
				},
			},
		},
	}

	resp, err := client.CreateMessage(context.Background(), req)

	if err != nil {
		return err
	}

	for _, block := range resp.Content {
		if block.Type == "text" {
			fmt.Println(block.Text)
		}
	}

	return nil
}
