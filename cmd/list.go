package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  `List all configuration profiles`,
	RunE: func(cmd *cobra.Command, args []string) error {
		binary := GetBinaryName()

		if binary == "" {
			return fmt.Errorf("binary name not specified")
		}

		service, err := getProfileService()
		if err != nil {
			return err
		}

		profiles, err := service.ListProfiles(binary)
		if err != nil {
			return err
		}

		if len(profiles) == 0 {
			fmt.Printf("No profiles found for %s\n", binary)
			fmt.Printf("Create one with: %s profile create <name>\n", binary)
			return nil
		}

		fmt.Printf("Profiles for %s:\n", binary)
		for _, profile := range profiles {
			envCount := len(profile.Environment())
			fmt.Printf("  - %s (%d env vars)\n", profile.Name(), envCount)
		}

		return nil
	},
}

func init() {
	profileCmd.AddCommand(listCmd)
}
