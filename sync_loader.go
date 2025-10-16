package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Resource represents a single infrastructure resource
type Resource struct {
	Type       string                 `yaml:"type"`
	Name       string                 `yaml:"name"`
	Properties map[string]interface{} `yaml:"properties"`
}

// SyncConfig represents the structure of the .sync file
type SyncConfig struct {
	Version   string     `yaml:"version"`
	Resources []Resource `yaml:"resources"`
}

// loadSyncFile reads and parses the .sync YAML file
func loadSyncFile(path string) (*SyncConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var config SyncConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig performs comprehensive validation on the loaded configuration
func validateConfig(config *SyncConfig) error {
	if config.Version == "" {
		return fmt.Errorf("version is required")
	}

	if len(config.Resources) == 0 {
		return fmt.Errorf("at least one resource is required")
	}

	// Track resource names to detect duplicates
	resourceKeys := make(map[string]bool)

	for i, resource := range config.Resources {
		// Validate resource type
		if resource.Type == "" {
			return fmt.Errorf("resource[%d]: type is required", i)
		}

		// Validate resource name
		if resource.Name == "" {
			return fmt.Errorf("resource[%d]: name is required", i)
		}

		// Check for duplicate resource keys
		key := resource.Key()
		if resourceKeys[key] {
			return fmt.Errorf("duplicate resource: %s", key)
		}
		resourceKeys[key] = true

		// Validate resource type is supported
		validTypes := map[string]bool{
			"server":       true,
			"database":     true,
			"loadbalancer": true,
		}
		if !validTypes[resource.Type] {
			return fmt.Errorf("resource[%d]: unsupported type '%s'", i, resource.Type)
		}

		// Validate properties exist
		if resource.Properties == nil || len(resource.Properties) == 0 {
			return fmt.Errorf("resource[%d]: properties are required", i)
		}
	}

	return nil
}
