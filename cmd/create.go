package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new profile",
	Long:  `Create a new configuration profile`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		profileName := args[0]
		binary := GetBinaryName()

		if binary == "" {
			return fmt.Errorf("binary name not specified")
		}

		service, err := getProfileService()
		if err != nil {
			return err
		}

		if err := service.CreateProfile(profileName, binary); err != nil {
			return err
		}

		fmt.Printf("✓ Profile '%s' created for %s\n", profileName, binary)
		fmt.Printf("  Edit at: ~/.config/wrapper/%s/%s.env\n", binary, profileName)
		fmt.Printf("  Set as current: %s profile set %s\n", binary, profileName)

		return nil
	},
}

func init() {
	profileCmd.AddCommand(createCmd)
}
