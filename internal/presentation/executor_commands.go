package presentation

import (
	"fmt"
	"os"

	"github.com/jycamier/wrapper/internal/application"
)

// ExecuteBinary executes a binary with the active profile environment
func ExecuteBinary(binaryName string, args []string) {
	// Setup dependencies
	repo, err := setupRepository()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	resolver, err := setupBinaryResolver()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	executorService := application.NewExecutorService(repo, resolver)

	// Execute binary
	if err := executorService.Execute(binaryName, args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
