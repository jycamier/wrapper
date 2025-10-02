package domain

import "errors"

// BinaryName is a value object representing a binary name
type BinaryName struct {
	value string
}

// NewBinaryName creates a new BinaryName
func NewBinaryName(name string) (BinaryName, error) {
	if name == "" {
		return BinaryName{}, errors.New("binary name cannot be empty")
	}
	return BinaryName{value: name}, nil
}

// String returns the string value
func (b BinaryName) String() string {
	return b.value
}

// Equals checks if two BinaryNames are equal
func (b BinaryName) Equals(other BinaryName) bool {
	return b.value == other.value
}
