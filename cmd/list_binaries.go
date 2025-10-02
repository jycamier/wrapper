package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// listBinariesCmd represents the list command for binaries
var listBinariesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all binaries with profiles",
	Long:  `List all binaries that have configuration profiles managed by wrapper`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := setupRepository()
		if err != nil {
			return err
		}

		// Get the wrapper config directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(homeDir, ".config", "wrapper")

		// Check if config directory exists
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			fmt.Println("No binaries configured yet")
			fmt.Println("Create a profile with: wrapper <binary> profile create <name>")
			return nil
		}

		// Read all directories in config dir
		entries, err := os.ReadDir(configDir)
		if err != nil {
			return fmt.Errorf("failed to read config directory: %w", err)
		}

		var binaries []string
		for _, entry := range entries {
			if entry.IsDir() {
				binaries = append(binaries, entry.Name())
			}
		}

		if len(binaries) == 0 {
			fmt.Println("No binaries configured yet")
			fmt.Println("Create a profile with: wrapper <binary> profile create <name>")
			return nil
		}

		fmt.Println("Configured binaries:")
		for _, binary := range binaries {
			// Get profiles for this binary
			profiles, err := repo.List(binary)
			if err != nil {
				continue
			}

			profileCount := len(profiles)
			if profileCount == 0 {
				fmt.Printf("  - %s (no profiles)\n", binary)
				continue
			}

			// Check current profile
			currentProfile, err := repo.GetCurrent(binary)
			currentName := ""
			if err == nil {
				currentName = currentProfile.Name()
			}

			if currentName != "" {
				fmt.Printf("  - %s (%d profiles, current: %s)\n", binary, profileCount, currentName)
			} else {
				fmt.Printf("  - %s (%d profiles)\n", binary, profileCount)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listBinariesCmd)
}
