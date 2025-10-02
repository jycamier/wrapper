package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// defaultCmd represents the default command
var defaultCmd = &cobra.Command{
	Use:   "default <name>",
	Short: "Set the default profile",
	Long:  `Set the default configuration profile (used when no current profile is set)`,
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

		if err := service.SetDefaultProfile(profileName, binary); err != nil {
			return err
		}

		fmt.Printf("✓ Default profile set to '%s' for %s\n", profileName, binary)

		return nil
	},
}

func init() {
	profileCmd.AddCommand(defaultCmd)
}
