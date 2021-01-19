package storage

import "fmt"

// Storage is the interface common to all storage
type Storage interface {
	// Check if a file exists
	FileExists(path string) bool
	// Upload a file to the storage
	Upload(string) error
}

// register functions creating new Storage instances.
var registry = make(map[string]Storage)

func Register(name string, storage Storage) error {
	if name == "" {
		return fmt.Errorf("Error registering storage: name cannot be empty")
	}
	if storage == nil {
		return fmt.Errorf("Error registering storage '%v': factory cannot be empty", name)
	}
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering storage '%v': already registered", name)
	}

	registry[name] = storage
	return nil
}

func GetStorage(name string) (Storage, error) {
	if _, exists := registry[name]; !exists {
		return nil, fmt.Errorf("Error creating storage. No such storage type exists: '%v'", name)
	}
	return registry[name], nil
}
