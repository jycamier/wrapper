package main

import (
	"os"
	"path/filepath"

	"github.com/jycamier/wrapper/cmd"
)

func main() {
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
