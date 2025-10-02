package infrastructure

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jycamier/wrapper/internal/domain"
)

// PathBinaryResolver resolves binaries from PATH, excluding the wrapper itself
type PathBinaryResolver struct {
	wrapperPath string
}

// NewPathBinaryResolver creates a new PathBinaryResolver
func NewPathBinaryResolver() (*PathBinaryResolver, error) {
	// Get wrapper's own path
	wrapperPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get wrapper path: %w", err)
	}

	// Resolve symlinks
	wrapperPath, err = filepath.EvalSymlinks(wrapperPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve wrapper path: %w", err)
	}

	return &PathBinaryResolver{wrapperPath: wrapperPath}, nil
}

// Resolve finds the real binary path, excluding the wrapper itself
func (r *PathBinaryResolver) Resolve(binaryName string) (string, error) {
	// Get PATH
	pathEnv := os.Getenv("PATH")
	if pathEnv == "" {
		return "", domain.ErrBinaryNotFound
	}

	// Split PATH
	paths := strings.Split(pathEnv, string(os.PathListSeparator))

	// Search for binary
	for _, dir := range paths {
		candidate := filepath.Join(dir, binaryName)

		// Check if file exists and is executable
		info, err := os.Stat(candidate)
		if err != nil {
			continue
		}

		// Check if it's a regular file
		if !info.Mode().IsRegular() {
			continue
		}

		// Check if it's executable
		if info.Mode().Perm()&0111 == 0 {
			continue
		}

		// Resolve symlinks
		resolvedPath, err := filepath.EvalSymlinks(candidate)
		if err != nil {
			continue
		}

		// Skip if it's the wrapper itself
		if resolvedPath == r.wrapperPath {
			continue
		}

		return candidate, nil
	}

	// Fallback: try exec.LookPath but verify it's not the wrapper
	path, err := exec.LookPath(binaryName)
	if err == nil {
		resolvedPath, err := filepath.EvalSymlinks(path)
		if err == nil && resolvedPath != r.wrapperPath {
			return path, nil
		}
	}

	return "", domain.ErrBinaryNotFound
}
