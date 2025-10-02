package domain

import (
	"errors"
	"fmt"
)

// Profile represents a configuration profile for a binary
type Profile struct {
	name        string
	binaryName  string
	environment map[string]string
}

// NewProfile creates a new profile
func NewProfile(name, binaryName string, env map[string]string) (*Profile, error) {
	if name == "" {
		return nil, errors.New("profile name cannot be empty")
	}
	if binaryName == "" {
		return nil, errors.New("binary name cannot be empty")
	}

	if env == nil {
		env = make(map[string]string)
	}

	return &Profile{
		name:        name,
		binaryName:  binaryName,
		environment: env,
	}, nil
}

// Name returns the profile name
func (p *Profile) Name() string {
	return p.name
}

// BinaryName returns the binary name
func (p *Profile) BinaryName() string {
	return p.binaryName
}

// Environment returns the environment variables
func (p *Profile) Environment() map[string]string {
	// Return a copy to prevent external modification
	env := make(map[string]string, len(p.environment))
	for k, v := range p.environment {
		env[k] = v
	}
	return env
}

// SetEnvironment updates the environment variables
func (p *Profile) SetEnvironment(env map[string]string) {
	p.environment = make(map[string]string, len(env))
	for k, v := range env {
		p.environment[k] = v
	}
}

// AddEnvironmentVariable adds or updates a single environment variable
func (p *Profile) AddEnvironmentVariable(key, value string) {
	p.environment[key] = value
}

// String returns a string representation of the profile
func (p *Profile) String() string {
	return fmt.Sprintf("Profile{name: %s, binary: %s, envCount: %d}",
		p.name, p.binaryName, len(p.environment))
}
