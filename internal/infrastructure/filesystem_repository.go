package infrastructure

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jycamier/wrapper/internal/domain"
)

// FilesystemRepository implements ProfileRepository using the filesystem
type FilesystemRepository struct {
	baseDir string
}

// NewFilesystemRepository creates a new FilesystemRepository
func NewFilesystemRepository() (*FilesystemRepository, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	baseDir := filepath.Join(homeDir, ".config", "wrapper")
	return &FilesystemRepository{baseDir: baseDir}, nil
}

// getBinaryDir returns the directory for a specific binary
func (r *FilesystemRepository) getBinaryDir(binaryName string) string {
	return filepath.Join(r.baseDir, binaryName)
}

// getProfilePath returns the file path for a profile
func (r *FilesystemRepository) getProfilePath(profileName, binaryName string) string {
	return filepath.Join(r.getBinaryDir(binaryName), profileName+".env")
}

// getCurrentSymlink returns the path to the current profile symlink
func (r *FilesystemRepository) getCurrentSymlink(binaryName string) string {
	return filepath.Join(r.getBinaryDir(binaryName), "current.env")
}

// getDefaultPath returns the path to the default profile marker
func (r *FilesystemRepository) getDefaultPath(binaryName string) string {
	return filepath.Join(r.getBinaryDir(binaryName), ".default")
}

// Save saves a profile to the filesystem
func (r *FilesystemRepository) Save(profile *domain.Profile) error {
	binaryDir := r.getBinaryDir(profile.BinaryName())

	// Create directory if it doesn't exist
	if err := os.MkdirAll(binaryDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	profilePath := r.getProfilePath(profile.Name(), profile.BinaryName())

	// Write environment variables
	file, err := os.Create(profilePath)
	if err != nil {
		return fmt.Errorf("failed to create profile file: %w", err)
	}
	defer file.Close()

	for key, value := range profile.Environment() {
		if _, err := fmt.Fprintf(file, "%s=%s\n", key, value); err != nil {
			return fmt.Errorf("failed to write environment variable: %w", err)
		}
	}

	return nil
}

// FindByName finds a profile by name
func (r *FilesystemRepository) FindByName(profileName, binaryName string) (*domain.Profile, error) {
	profilePath := r.getProfilePath(profileName, binaryName)

	// Check if file exists
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return nil, domain.ErrProfileNotFound
	}

	// Read environment variables
	env, err := r.readEnvFile(profilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile: %w", err)
	}

	return domain.NewProfile(profileName, binaryName, env)
}

// List lists all profiles for a binary
func (r *FilesystemRepository) List(binaryName string) ([]*domain.Profile, error) {
	binaryDir := r.getBinaryDir(binaryName)

	// Check if directory exists
	if _, err := os.Stat(binaryDir); os.IsNotExist(err) {
		return []*domain.Profile{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(binaryDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var profiles []*domain.Profile
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".env") {
			continue
		}

		// Skip current.env symlink
		if entry.Name() == "current.env" {
			continue
		}

		// Extract profile name
		profileName := strings.TrimSuffix(entry.Name(), ".env")

		// Load profile
		profile, err := r.FindByName(profileName, binaryName)
		if err != nil {
			continue // Skip invalid profiles
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// Delete deletes a profile
func (r *FilesystemRepository) Delete(profileName, binaryName string) error {
	profilePath := r.getProfilePath(profileName, binaryName)

	// Check if file exists
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return domain.ErrProfileNotFound
	}

	// Delete file
	if err := os.Remove(profilePath); err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	return nil
}

// SetCurrent sets the current profile
func (r *FilesystemRepository) SetCurrent(profileName, binaryName string) error {
	profilePath := r.getProfilePath(profileName, binaryName)

	// Check if profile exists
	if _, err := os.Stat(profilePath); os.IsNotExist(err) {
		return domain.ErrProfileNotFound
	}

	symlinkPath := r.getCurrentSymlink(binaryName)

	// Remove existing symlink if it exists
	_ = os.Remove(symlinkPath)

	// Create new symlink
	if err := os.Symlink(profilePath, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	return nil
}

// GetCurrent gets the current profile
func (r *FilesystemRepository) GetCurrent(binaryName string) (*domain.Profile, error) {
	symlinkPath := r.getCurrentSymlink(binaryName)

	// Check if symlink exists
	if _, err := os.Lstat(symlinkPath); os.IsNotExist(err) {
		return nil, domain.ErrNoCurrentProfile
	}

	// Read target
	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read symlink: %w", err)
	}

	// Extract profile name from target
	profileName := strings.TrimSuffix(filepath.Base(target), ".env")

	return r.FindByName(profileName, binaryName)
}

// SetDefault sets the default profile
func (r *FilesystemRepository) SetDefault(profileName, binaryName string) error {
	// Check if profile exists
	if _, err := r.FindByName(profileName, binaryName); err != nil {
		return err
	}

	defaultPath := r.getDefaultPath(binaryName)

	// Write default profile name
	if err := os.WriteFile(defaultPath, []byte(profileName), 0644); err != nil {
		return fmt.Errorf("failed to write default profile: %w", err)
	}

	return nil
}

// GetDefault gets the default profile
func (r *FilesystemRepository) GetDefault(binaryName string) (*domain.Profile, error) {
	defaultPath := r.getDefaultPath(binaryName)

	// Check if file exists
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		return nil, domain.ErrNoDefaultProfile
	}

	// Read default profile name
	data, err := os.ReadFile(defaultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read default profile: %w", err)
	}

	profileName := strings.TrimSpace(string(data))
	return r.FindByName(profileName, binaryName)
}

// GetActiveProfile gets the active profile (current if set, otherwise default)
func (r *FilesystemRepository) GetActiveProfile(binaryName string) (*domain.Profile, error) {
	// Try current first
	profile, err := r.GetCurrent(binaryName)
	if err == nil {
		return profile, nil
	}

	// Fall back to default
	profile, err = r.GetDefault(binaryName)
	if err == nil {
		return profile, nil
	}

	return nil, domain.ErrNoCurrentProfile
}

// readEnvFile reads environment variables from a .env file
func (r *FilesystemRepository) readEnvFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return env, nil
}
