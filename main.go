package main

import (
	"os"
	"path/filepath"

	"github.com/jycamier/wrapper/internal/presentation"
)

func main() {
	// Determine binary name from how wrapper was invoked
	binaryName := filepath.Base(os.Args[0])

	// If called as "wrapper", the first arg is the binary name
	if binaryName == "wrapper" {
		if len(os.Args) < 2 {
			// Show usage
			presentation.ShowUsage()
			os.Exit(1)
		}

		binaryName = os.Args[1]
		os.Args = append(os.Args[:1], os.Args[2:]...) // Remove binary name from args
	}

	// Check if this is a profile command
	if len(os.Args) > 1 && os.Args[1] == "profile" {
		// Execute profile command
		presentation.ExecuteProfileCommands(binaryName, os.Args[2:])
		return
	}

	// Execute the real binary with active profile environment
	presentation.ExecuteBinary(binaryName, os.Args[1:])
}
