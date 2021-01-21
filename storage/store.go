package storage

import (
	"fmt"

	ucfg "github.com/elastic/go-ucfg"
)

// Storage is the interface common to all storage
type Storage interface {
	// Check if a file exists
	FileExists(path string) bool
	// Upload a file to the storage
	Upload(string) error
}

//
type Factory = func(config *ucfg.Config) (Storage, error)

var registry = make(map[string]Factory)

func Register(name string, factory Factory) error {
	if name == "" {
		return fmt.Errorf("Error registering storage: name cannot be empty")
	}
	if factory == nil {
		return fmt.Errorf("Error registering storage '%v': factory cannot be empty", name)
	}
	if _, exists := registry[name]; exists {
		return fmt.Errorf("Error registering storage '%v': already registered", name)
	}

	registry[name] = factory
	return nil
}

func GetFactory(name string) (Factory, error) {
	if _, exists := registry[name]; !exists {
		return nil, fmt.Errorf("Error creating storage. No such storage type exists: '%v'", name)
	}
	return registry[name], nil
}
