package cmd

import (
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command
var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long:  `Manage configuration profiles for the binary`,
}

func init() {
	rootCmd.AddCommand(profileCmd)
}
