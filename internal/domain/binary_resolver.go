package domain

import "errors"

var (
	// ErrBinaryNotFound is returned when a binary is not found in PATH
	ErrBinaryNotFound = errors.New("binary not found in PATH")
)

// BinaryResolver defines the interface for resolving binary paths
type BinaryResolver interface {
	// Resolve finds the real binary path, excluding the wrapper itself
	Resolve(binaryName string) (string, error)
}
