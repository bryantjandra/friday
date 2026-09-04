package main

import (
	"fmt"
	"os"
)

/* Go forces main to take no parameters, and no return value. Thus, we use a seperate run() function to know what error was produced*/
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
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}
