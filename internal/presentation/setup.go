package presentation

import (
	"github.com/jycamier/wrapper/internal/domain"
	"github.com/jycamier/wrapper/internal/infrastructure"
)

// setupRepository creates and configures the profile repository
func setupRepository() (domain.ProfileRepository, error) {
	return infrastructure.NewFilesystemRepository()
}

// setupBinaryResolver creates and configures the binary resolver
func setupBinaryResolver() (domain.BinaryResolver, error) {
	return infrastructure.NewPathBinaryResolver()
}
