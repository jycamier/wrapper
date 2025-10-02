package application

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jycamier/wrapper/internal/domain"
)

// ExecutorService handles binary execution with profile environment
type ExecutorService struct {
	profileRepo    domain.ProfileRepository
	binaryResolver domain.BinaryResolver
}

// NewExecutorService creates a new ExecutorService
func NewExecutorService(profileRepo domain.ProfileRepository, binaryResolver domain.BinaryResolver) *ExecutorService {
	return &ExecutorService{
		profileRepo:    profileRepo,
		binaryResolver: binaryResolver,
	}
}

// Execute executes a binary with the active profile environment
func (s *ExecutorService) Execute(binaryName string, args []string) error {
	// Get active profile
	profile, err := s.profileRepo.GetActiveProfile(binaryName)
	if err != nil {
		if err == domain.ErrNoCurrentProfile || err == domain.ErrNoDefaultProfile {
			return fmt.Errorf("no active profile for '%s': use '%s profile create <name>' to create one", binaryName, binaryName)
		}
		return fmt.Errorf("failed to get active profile: %w", err)
	}

	// Resolve real binary path
	binaryPath, err := s.binaryResolver.Resolve(binaryName)
	if err != nil {
		if err == domain.ErrBinaryNotFound {
			return fmt.Errorf("binary '%s' not found in PATH", binaryName)
		}
		return fmt.Errorf("failed to resolve binary: %w", err)
	}

	// Prepare command
	cmd := exec.Command(binaryPath, args...)

	// Set environment variables from profile
	cmd.Env = os.Environ()
	for key, value := range profile.Environment() {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Wire up stdio
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute
	if err := cmd.Run(); err != nil {
		// Preserve exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return fmt.Errorf("failed to execute binary: %w", err)
	}

	return nil
}
