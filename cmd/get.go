package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current profile",
	Long:  `Get the current configuration profile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		binary := GetBinaryName()

		if binary == "" {
			return fmt.Errorf("binary name not specified")
		}

		service, err := getProfileService()
		if err != nil {
			return err
		}

		profileName, err := service.GetCurrentProfile(binary)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", profileName)

		return nil
	},
}

func init() {
	profileCmd.AddCommand(getCmd)
}
