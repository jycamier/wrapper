package application

import (
	"fmt"

	"github.com/jycamier/wrapper/internal/domain"
)

// ProfileService handles profile-related use cases
type ProfileService struct {
	repo domain.ProfileRepository
}

// NewProfileService creates a new ProfileService
func NewProfileService(repo domain.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

// CreateProfile creates a new profile
func (s *ProfileService) CreateProfile(name, binaryName string) error {
	// Check if profile already exists
	existing, err := s.repo.FindByName(name, binaryName)
	if err == nil && existing != nil {
		return fmt.Errorf("profile '%s' for binary '%s' already exists", name, binaryName)
	}

	// Create new profile
	profile, err := domain.NewProfile(name, binaryName, make(map[string]string))
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	// Save profile
	if err := s.repo.Save(profile); err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	return nil
}

// ListProfiles lists all profiles for a binary
func (s *ProfileService) ListProfiles(binaryName string) ([]*domain.Profile, error) {
	profiles, err := s.repo.List(binaryName)
	if err != nil {
		return nil, fmt.Errorf("failed to list profiles: %w", err)
	}
	return profiles, nil
}

// SetCurrentProfile sets the current profile for a binary
func (s *ProfileService) SetCurrentProfile(name, binaryName string) error {
	// Verify profile exists
	_, err := s.repo.FindByName(name, binaryName)
	if err != nil {
		if err == domain.ErrProfileNotFound {
			return fmt.Errorf("profile '%s' not found for binary '%s'", name, binaryName)
		}
		return fmt.Errorf("failed to find profile: %w", err)
	}

	// Set as current
	if err := s.repo.SetCurrent(name, binaryName); err != nil {
		return fmt.Errorf("failed to set current profile: %w", err)
	}

	return nil
}

// GetCurrentProfile gets the current profile name for a binary
func (s *ProfileService) GetCurrentProfile(binaryName string) (string, error) {
	profile, err := s.repo.GetCurrent(binaryName)
	if err != nil {
		if err == domain.ErrNoCurrentProfile {
			return "", fmt.Errorf("no current profile set for binary '%s'", binaryName)
		}
		return "", fmt.Errorf("failed to get current profile: %w", err)
	}

	return profile.Name(), nil
}

// SetDefaultProfile sets the default profile for a binary
func (s *ProfileService) SetDefaultProfile(name, binaryName string) error {
	// Verify profile exists
	_, err := s.repo.FindByName(name, binaryName)
	if err != nil {
		if err == domain.ErrProfileNotFound {
			return fmt.Errorf("profile '%s' not found for binary '%s'", name, binaryName)
		}
		return fmt.Errorf("failed to find profile: %w", err)
	}

	// Set as default
	if err := s.repo.SetDefault(name, binaryName); err != nil {
		return fmt.Errorf("failed to set default profile: %w", err)
	}

	return nil
}

// GetActiveProfile gets the active profile (current if set, otherwise default)
func (s *ProfileService) GetActiveProfile(binaryName string) (*domain.Profile, error) {
	profile, err := s.repo.GetActiveProfile(binaryName)
	if err != nil {
		if err == domain.ErrNoCurrentProfile || err == domain.ErrNoDefaultProfile {
			return nil, fmt.Errorf("no active profile for binary '%s': please create and set a profile first", binaryName)
		}
		return nil, fmt.Errorf("failed to get active profile: %w", err)
	}

	return profile, nil
}
