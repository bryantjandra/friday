package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bryantjandra/friday/internal/agent"
	"github.com/bryantjandra/friday/internal/config"
	"github.com/bryantjandra/friday/internal/llm"
	"github.com/bryantjandra/friday/internal/store"
	"github.com/bryantjandra/friday/internal/tools"
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
	case "tooltest":
		return runToolTest()
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

	registry := tools.NewRegistry()
	registry.RegisterTool(tools.NewCreateReminderTool(db))
	client := llm.NewClient(cfg)

	ag := agent.NewAgent(client, registry)
	result, err := ag.Run(context.Background(), prompt)

	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

func runToolTest() error {
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

	/* Step 4: Build the registry and register the tool */
	registry := tools.NewRegistry()
	registry.RegisterTool(tools.NewCreateReminderTool(db))

	/* Step 5: look up the tool by name */
	tool, ok := registry.GetTool("create_reminder")
	if !ok {
		return fmt.Errorf("tool not found in registry")
	}

	args := json.RawMessage(`{"title":"mikis birthday","due_at":"2026-10-05T00:00:00Z"}`)
	result, err := tool.Execute(context.Background(), args)

	if err != nil {
		return err
	}

	fmt.Printf("Result: IsError=%v, Content=%q\n", result.IsError, result.Content)
	return nil
}
