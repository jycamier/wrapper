package domain

import "errors"

// ProfileName is a value object representing a profile name
type ProfileName struct {
	value string
}

// NewProfileName creates a new ProfileName
func NewProfileName(name string) (ProfileName, error) {
	if name == "" {
		return ProfileName{}, errors.New("profile name cannot be empty")
	}
	return ProfileName{value: name}, nil
}

// String returns the string value
func (p ProfileName) String() string {
	return p.value
}

// Equals checks if two ProfileNames are equal
func (p ProfileName) Equals(other ProfileName) bool {
	return p.value == other.value
}
