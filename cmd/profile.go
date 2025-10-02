package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long:  `Manage configuration profiles for the binary`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand, default to list
		binary := GetBinaryName()

		if binary == "" {
			return fmt.Errorf("binary name not specified")
		}

		repo, err := setupRepository()
		if err != nil {
			return err
		}

		profiles, err := repo.List(binary)
		if err != nil {
			return err
		}

		if len(profiles) == 0 {
			fmt.Printf("No profiles found for %s\n", binary)
			fmt.Printf("Create one with: %s profile create <name>\n", binary)
			return nil
		}

		// Get current profile
		currentProfile, _ := repo.GetCurrent(binary)
		var currentName string
		if currentProfile != nil {
			currentName = currentProfile.Name()
		}

		// Get default profile
		defaultProfile, _ := repo.GetDefault(binary)
		var defaultName string
		if defaultProfile != nil {
			defaultName = defaultProfile.Name()
		}

		fmt.Printf("%sProfiles for %s:%s\n", colorCyan, binary, colorReset)
		for _, profile := range profiles {
			envCount := len(profile.Environment())
			name := profile.Name()

			// Determine display
			if name == currentName && name == defaultName {
				// Both current and default
				fmt.Printf("  %s✓%s %s%s%s (%d env vars) %s[current, default]%s\n",
					colorGreen, colorReset,
					colorGreen, name, colorReset,
					envCount,
					colorYellow, colorReset)
			} else if name == currentName {
				// Current only
				fmt.Printf("  %s✓%s %s%s%s (%d env vars) %s[current]%s\n",
					colorGreen, colorReset,
					colorGreen, name, colorReset,
					envCount,
					colorYellow, colorReset)
			} else if name == defaultName {
				// Default only
				fmt.Printf("  %s●%s %s (%d env vars) %s[default]%s\n",
					colorYellow, colorReset,
					name,
					envCount,
					colorYellow, colorReset)
			} else {
				// Regular profile
				fmt.Printf("  - %s (%d env vars)\n", name, envCount)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
