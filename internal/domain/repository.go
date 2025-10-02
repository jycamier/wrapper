package domain

import "errors"

var (
	// ErrProfileNotFound is returned when a profile is not found
	ErrProfileNotFound = errors.New("profile not found")
	// ErrProfileAlreadyExists is returned when trying to create a profile that already exists
	ErrProfileAlreadyExists = errors.New("profile already exists")
	// ErrNoCurrentProfile is returned when no current profile is set
	ErrNoCurrentProfile = errors.New("no current profile set")
	// ErrNoDefaultProfile is returned when no default profile is set
	ErrNoDefaultProfile = errors.New("no default profile set")
)

// ProfileRepository defines the interface for profile persistence
type ProfileRepository interface {
	// Save saves a profile
	Save(profile *Profile) error

	// FindByName finds a profile by name and binary name
	FindByName(profileName, binaryName string) (*Profile, error)

	// List lists all profiles for a binary
	List(binaryName string) ([]*Profile, error)

	// Delete deletes a profile
	Delete(profileName, binaryName string) error

	// SetCurrent sets the current profile for a binary
	SetCurrent(profileName, binaryName string) error

	// GetCurrent gets the current profile for a binary
	GetCurrent(binaryName string) (*Profile, error)

	// SetDefault sets the default profile for a binary
	SetDefault(profileName, binaryName string) error

	// GetDefault gets the default profile for a binary
	GetDefault(binaryName string) (*Profile, error)

	// GetActiveProfile gets the active profile (current if set, otherwise default)
	GetActiveProfile(binaryName string) (*Profile, error)
}
