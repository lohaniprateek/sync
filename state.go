package main

import (
	"encoding/json"
)

// getCurrentState retrieves the current infrastructure state
// In a real implementation, this would query your actual infrastructure
func getCurrentState() (*SyncConfig, error) {
	// For demonstration purposes, this returns a mock current state
	// In a real scenario, you would:
	// - Query cloud provider APIs
	// - Read from a state file
	// - Connect to a database
	// - etc.

	return &SyncConfig{
		Version: "1.0",
		Resources: []Resource{
			{
				Type: "server",
				Name: "web-server-1",
				Properties: map[string]interface{}{
					"instance_type": "t2.micro",
					"region":        "us-east-1",
					"status":        "running",
				},
			},
			{
				Type: "database",
				Name: "main-db",
				Properties: map[string]interface{}{
					"engine":  "postgres",
					"version": "14.0",
					"size":    "small",
				},
			},
		},
	}, nil
}

// ResourceKey creates a unique identifier for a resource
func (r Resource) Key() string {
	return r.Type + "." + r.Name
}

// DeepEqual compares two property maps safely
func deepEqual(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare each key-value pair
	for key, valueA := range a {
		valueB, exists := b[key]
		if !exists {
			return false
		}

		if !deepEqualValue(valueA, valueB) {
			return false
		}
	}

	return true
}

// deepEqualValue compares individual values, handling different types
func deepEqualValue(a, b interface{}) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// For complex types, use JSON comparison as fallback
	// but handle marshaling errors properly
	aJSON, errA := json.Marshal(a)
	bJSON, errB := json.Marshal(b)

	if errA != nil || errB != nil {
		// If marshaling fails, fall back to direct comparison
		return a == b
	}

	return string(aJSON) == string(bJSON)
}
