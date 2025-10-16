package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: sync <path/to/.sync>")
		os.Exit(1)
	}

	syncFilePath := os.Args[1]

	if err := run(syncFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(syncFilePath string) error {
	// Load desired state from .sync file
	desiredState, err := loadSyncFile(syncFilePath)
	if err != nil {
		return fmt.Errorf("failed to load sync file: %w", err)
	}

	// Get current state from infrastructure
	currentState, err := getCurrentState()
	if err != nil {
		return fmt.Errorf("failed to get current state: %w", err)
	}
	// Compare states and generate diff
	diff := compareStates(currentState, desiredState)

	// Display the changes
	displayDiff(diff)

	return nil
}
