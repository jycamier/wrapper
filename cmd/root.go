package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/jycamier/wrapper/internal/application"
	"github.com/jycamier/wrapper/internal/infrastructure"
	"github.com/spf13/cobra"
)

var (
	binaryName string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wrapper <binary> [command]",
	Short: "Configuration profile manager for any binary",
	Long: `Wrapper is a transparent configuration profile manager for any CLI binary.
It allows you to manage multiple environment configurations (profiles) and
seamlessly switch between them.

Example:
  # Create a shell function
  vault() { wrapper vault "$@"; }

  # Create and use a profile
  vault profile create prod
  vault profile set prod
  vault status  # Runs with prod profile environment`,
	SilenceUsage:         true,
	SilenceErrors:        true,
	DisableFlagParsing:   false,
	DisableSuggestions:   true,
	DisableFlagsInUseLine: true,
	// If no known subcommand is found, execute the binary
	RunE: func(cmd *cobra.Command, args []string) error {
		if binaryName == "" {
			return fmt.Errorf("binary name not specified")
		}

		// Execute the real binary with all args
		service, err := getExecutorService()
		if err != nil {
			return err
		}

		return service.Execute(binaryName, args)
	},
}

// isWrapperCommand checks if the argument is a known wrapper command
func isWrapperCommand(arg string) bool {
	wrapperCommands := []string{"list", "alias", "version", "help", "--help", "-h", "completion"}
	for _, cmd := range wrapperCommands {
		if arg == cmd {
			return true
		}
	}
	return false
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	// Determine binary name from invocation
	// Skip if first arg is a wrapper command
	if len(os.Args) >= 2 && !isWrapperCommand(os.Args[1]) {
		binaryName = os.Args[1]
		// Remove binary name from args for cobra parsing
		os.Args = append(os.Args[:1], os.Args[2:]...)
	}

	// Handle unknown command errors by executing the binary
	if err := rootCmd.Execute(); err != nil {
		// Check if it's an unknown command error
		errStr := err.Error()
		if binaryName != "" && (strings.Contains(errStr, "unknown command") || strings.Contains(errStr, "unknown flag")) {
			// Try to execute as binary command
			service, execErr := getExecutorService()
			if execErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", execErr)
				os.Exit(1)
			}

			// Get remaining args after "wrapper <binary>"
			if execErr := service.Execute(binaryName, os.Args[1:]); execErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", execErr)
				os.Exit(1)
			}
			return
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// ExecuteWithBinary executes a command for a specific binary (used when invoked via symlink)
func ExecuteWithBinary(binary string, args []string) {
	binaryName = binary
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		// Check if it's an unknown command error
		errStr := err.Error()
		if strings.Contains(errStr, "unknown command") || strings.Contains(errStr, "unknown flag") {
			// Try to execute as binary command
			service, execErr := getExecutorService()
			if execErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", execErr)
				os.Exit(1)
			}

			if execErr := service.Execute(binaryName, args); execErr != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", execErr)
				os.Exit(1)
			}
			return
		}

		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// GetBinaryName returns the current binary name
func GetBinaryName() string {
	return binaryName
}

// setupRepository creates the profile repository
func setupRepository() (*infrastructure.FilesystemRepository, error) {
	return infrastructure.NewFilesystemRepository()
}

// setupBinaryResolver creates the binary resolver
func setupBinaryResolver() (*infrastructure.PathBinaryResolver, error) {
	return infrastructure.NewPathBinaryResolver()
}

// getProfileService returns an initialized ProfileService
func getProfileService() (*application.ProfileService, error) {
	repo, err := setupRepository()
	if err != nil {
		return nil, err
	}
	return application.NewProfileService(repo), nil
}

// getExecutorService returns an initialized ExecutorService
func getExecutorService() (*application.ExecutorService, error) {
	repo, err := setupRepository()
	if err != nil {
		return nil, err
	}

	resolver, err := setupBinaryResolver()
	if err != nil {
		return nil, err
	}

	return application.NewExecutorService(repo, resolver), nil
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	// Allow unknown flags and commands to pass through to the real binary
	rootCmd.FParseErrWhitelist = cobra.FParseErrWhitelist{UnknownFlags: true}
}
