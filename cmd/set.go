package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <name>",
	Short: "Set the current profile",
	Long:  `Set the current configuration profile`,
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

		if err := service.SetCurrentProfile(profileName, binary); err != nil {
			return err
		}

		fmt.Printf("✓ Current profile set to '%s' for %s\n", profileName, binary)

		return nil
	},
}

func init() {
	profileCmd.AddCommand(setCmd)
}
