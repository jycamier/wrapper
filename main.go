package main

import (
	"os"
	"path/filepath"

	"github.com/jycamier/wrapper/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info in cmd package
	cmd.Version = version
	cmd.Commit = commit
	cmd.Date = date

	// Determine binary name from how wrapper was invoked
	binaryName := filepath.Base(os.Args[0])

	// If called as "wrapper", use standard cobra execution
	if binaryName == "wrapper" {
		cmd.Execute()
		return
	}

	// Called via symlink/alias (e.g., vault, aws, kubectl)
	// Pass all args to ExecuteWithBinary - Cobra will route to profile or execute
	cmd.ExecuteWithBinary(binaryName, os.Args[1:])
}
