package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// displayDiff displays the differences in a formatted way
func displayDiff(diff *Diff) {
	if !diff.HasChanges() {
		color.Green("\n✓ No changes detected. Infrastructure is up to date.\n")
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 70))
	color.Cyan("  Infrastructure Change Plan")
	fmt.Println(strings.Repeat("=", 70))

	createCount := 0
	updateCount := 0
	deleteCount := 0

	for _, change := range diff.Changes {
		switch change.Type {
		case ChangeTypeCreate:
			createCount++
			displayCreate(change)

		case ChangeTypeUpdate:
			updateCount++
			displayUpdate(change)

		case ChangeTypeDelete:
			deleteCount++
			displayDelete(change)
		}
	}

	// Summary
	fmt.Println("\n" + strings.Repeat("-", 70))
	color.Cyan("\nPlan Summary:")
	fmt.Printf("  ")
	if createCount > 0 {
		color.Green("+ %d to create", createCount)
		fmt.Printf("  ")
	}
	if updateCount > 0 {
		color.Yellow("~ %d to update", updateCount)
		fmt.Printf("  ")
	}
	if deleteCount > 0 {
		color.Red("- %d to delete", deleteCount)
	}
	fmt.Println()

	fmt.Println("\n" + strings.Repeat("=", 70))
	color.Cyan("\nNote: This is a preview. No changes have been applied.")
	fmt.Println(strings.Repeat("=", 70) + "\n")
}

// displayCreate displays a resource that will be created
func displayCreate(change ResourceChange) {
	fmt.Println()
	color.Green("+ CREATE: %s.%s", change.Resource.Type, change.Resource.Name)
	fmt.Println(strings.Repeat("-", 50))

	for key, value := range change.Resource.Properties {
		color.Green("    + %s: %v", key, value)
	}
}

// displayUpdate displays a resource that will be updated
func displayUpdate(change ResourceChange) {
	fmt.Println()
	color.Yellow("~ UPDATE: %s.%s", change.Resource.Type, change.Resource.Name)
	fmt.Println(strings.Repeat("-", 50))

	for _, propChange := range change.PropertyChanges {
		if propChange.OldValue == nil && propChange.NewValue != nil {
			color.Green("    + %s: %v", propChange.Key, propChange.NewValue)
		} else if propChange.OldValue != nil && propChange.NewValue == nil {
			color.Red("    - %s: %v", propChange.Key, propChange.OldValue)
		} else {
			color.Yellow("    ~ %s:", propChange.Key)
			color.Red("        - %v", propChange.OldValue)
			color.Green("        + %v", propChange.NewValue)
		}
	}
}

// displayDelete displays a resource that will be deleted
func displayDelete(change ResourceChange) {
	fmt.Println()
	color.Red("- DELETE: %s.%s", change.Resource.Type, change.Resource.Name)
	fmt.Println(strings.Repeat("-", 50))

	for key, value := range change.Resource.Properties {
		color.Red("    - %s: %v", key, value)
	}
}
