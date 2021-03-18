package utils

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	log, err := NewLogger("logger")
	if err != nil {
		t.Fatal(err)
	}
	log.Info("Simple logging message")
}
