package presentation

import (
	"fmt"
	"os"

	"github.com/jycamier/wrapper/internal/application"
	"github.com/spf13/cobra"
)

// ProfileCommands handles profile-related CLI commands
type ProfileCommands struct {
	profileService *application.ProfileService
}

// NewProfileCommands creates a new ProfileCommands
func NewProfileCommands(profileService *application.ProfileService) *ProfileCommands {
	return &ProfileCommands{profileService: profileService}
}

// BuildProfileCommand builds the profile command tree
func (p *ProfileCommands) BuildProfileCommand(binaryName string) *cobra.Command {
	profileCmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage configuration profiles",
		Long:  fmt.Sprintf("Manage configuration profiles for %s", binaryName),
	}

	profileCmd.AddCommand(p.buildListCommand(binaryName))
	profileCmd.AddCommand(p.buildCreateCommand(binaryName))
	profileCmd.AddCommand(p.buildSetCommand(binaryName))
	profileCmd.AddCommand(p.buildGetCommand(binaryName))
	profileCmd.AddCommand(p.buildDefaultCommand(binaryName))

	return profileCmd
}

func (p *ProfileCommands) buildListCommand(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		Long:  fmt.Sprintf("List all configuration profiles for %s", binaryName),
		RunE: func(cmd *cobra.Command, args []string) error {
			profiles, err := p.profileService.ListProfiles(binaryName)
			if err != nil {
				return err
			}

			if len(profiles) == 0 {
				fmt.Printf("No profiles found for %s\n", binaryName)
				fmt.Printf("Create one with: %s profile create <name>\n", binaryName)
				return nil
			}

			fmt.Printf("Profiles for %s:\n", binaryName)
			for _, profile := range profiles {
				envCount := len(profile.Environment())
				fmt.Printf("  - %s (%d env vars)\n", profile.Name(), envCount)
			}

			return nil
		},
	}
}

func (p *ProfileCommands) buildCreateCommand(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new profile",
		Long:  fmt.Sprintf("Create a new configuration profile for %s", binaryName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			if err := p.profileService.CreateProfile(profileName, binaryName); err != nil {
				return err
			}

			fmt.Printf("✓ Profile '%s' created for %s\n", profileName, binaryName)
			fmt.Printf("  Edit at: ~/.config/wrapper/%s/%s.env\n", binaryName, profileName)
			fmt.Printf("  Set as current: %s profile set %s\n", binaryName, profileName)

			return nil
		},
	}
}

func (p *ProfileCommands) buildSetCommand(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name>",
		Short: "Set the current profile",
		Long:  fmt.Sprintf("Set the current configuration profile for %s", binaryName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			if err := p.profileService.SetCurrentProfile(profileName, binaryName); err != nil {
				return err
			}

			fmt.Printf("✓ Current profile set to '%s' for %s\n", profileName, binaryName)

			return nil
		},
	}
}

func (p *ProfileCommands) buildGetCommand(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get the current profile",
		Long:  fmt.Sprintf("Get the current configuration profile for %s", binaryName),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName, err := p.profileService.GetCurrentProfile(binaryName)
			if err != nil {
				return err
			}

			fmt.Printf("%s\n", profileName)

			return nil
		},
	}
}

func (p *ProfileCommands) buildDefaultCommand(binaryName string) *cobra.Command {
	return &cobra.Command{
		Use:   "default <name>",
		Short: "Set the default profile",
		Long:  fmt.Sprintf("Set the default configuration profile for %s (used when no current profile is set)", binaryName),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			if err := p.profileService.SetDefaultProfile(profileName, binaryName); err != nil {
				return err
			}

			fmt.Printf("✓ Default profile set to '%s' for %s\n", profileName, binaryName)

			return nil
		},
	}
}

// ExecuteProfileCommands executes the profile command tree
func ExecuteProfileCommands(binaryName string, args []string) {
	// Setup dependencies
	repo, err := setupRepository()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	profileService := application.NewProfileService(repo)
	commands := NewProfileCommands(profileService)

	rootCmd := commands.BuildProfileCommand(binaryName)
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
