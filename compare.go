package main

// ChangeType represents the type of change detected
type ChangeType string

const (
	ChangeTypeCreate ChangeType = "create"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
	ChangeTypeNone   ChangeType = "no-change"
)

// PropertyChange represents a change in a resource property
type PropertyChange struct {
	Key      string
	OldValue interface{}
	NewValue interface{}
}

// ResourceChange represents a change to a resource
type ResourceChange struct {
	Type            ChangeType
	Resource        Resource
	CurrentResource *Resource
	PropertyChanges []PropertyChange
}

// Diff represents all changes between current and desired state
type Diff struct {
	Changes []ResourceChange
}

// HasChanges returns true if there are any changes
func (d *Diff) HasChanges() bool {
	for _, change := range d.Changes {
		if change.Type != ChangeTypeNone {
			return true
		}
	}
	return false
}

// compareStates compares current and desired states and returns the differences
func compareStates(current, desired *SyncConfig) *Diff {
	diff := &Diff{
		Changes: []ResourceChange{},
	}

	// Create maps for quick lookup
	currentMap := make(map[string]*Resource)
	for i := range current.Resources {
		res := &current.Resources[i]
		currentMap[res.Key()] = res
	}

	desiredMap := make(map[string]*Resource)
	for i := range desired.Resources {
		res := &desired.Resources[i]
		desiredMap[res.Key()] = res
	}

	// Check for creates and updates
	for key, desiredRes := range desiredMap {
		currentRes, exists := currentMap[key]

		if !exists {
			// Resource needs to be created
			diff.Changes = append(diff.Changes, ResourceChange{
				Type:     ChangeTypeCreate,
				Resource: *desiredRes,
			})
		} else {
			// Resource exists, check for updates
			change := compareResources(currentRes, desiredRes)
			diff.Changes = append(diff.Changes, change)
		}
	}

	// Check for deletes
	for key, currentRes := range currentMap {
		if _, exists := desiredMap[key]; !exists {
			// Resource needs to be deleted
			diff.Changes = append(diff.Changes, ResourceChange{
				Type:            ChangeTypeDelete,
				Resource:        *currentRes,
				CurrentResource: currentRes,
			})
		}
	}

	return diff
}

// compareResources compares two resources and returns the changes
func compareResources(current, desired *Resource) ResourceChange {
	change := ResourceChange{
		Type:            ChangeTypeNone,
		Resource:        *desired,
		CurrentResource: current,
		PropertyChanges: []PropertyChange{},
	}

	// Compare properties
	if !deepEqual(current.Properties, desired.Properties) {
		change.Type = ChangeTypeUpdate

		// Find specific property changes
		allKeys := make(map[string]bool)
		for k := range current.Properties {
			allKeys[k] = true
		}
		for k := range desired.Properties {
			allKeys[k] = true
		}

		for key := range allKeys {
			currentVal := current.Properties[key]
			desiredVal := desired.Properties[key]

			// Simple comparison (in production, use deep comparison)
			if currentVal != desiredVal {
				change.PropertyChanges = append(change.PropertyChanges, PropertyChange{
					Key:      key,
					OldValue: currentVal,
					NewValue: desiredVal,
				})
			}
		}
	}

	return change
}
