package presentation

import "fmt"

// ShowUsage displays usage information
func ShowUsage() {
	fmt.Println("wrapper - Configuration profile manager for any binary")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wrapper <binary> [args...]           Execute binary with active profile")
	fmt.Println("  wrapper <binary> profile list        List all profiles")
	fmt.Println("  wrapper <binary> profile create <name>")
	fmt.Println("                                       Create a new profile")
	fmt.Println("  wrapper <binary> profile set <name>  Set current profile")
	fmt.Println("  wrapper <binary> profile get         Get current profile")
	fmt.Println("  wrapper <binary> profile default <name>")
	fmt.Println("                                       Set default profile")
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  # Create a shell function for vault")
	fmt.Println("  vault() { wrapper vault \"$@\"; }")
	fmt.Println()
	fmt.Println("  # Create a profile")
	fmt.Println("  vault profile create prod")
	fmt.Println()
	fmt.Println("  # Edit ~/.config/wrapper/vault/prod.env with your variables")
	fmt.Println("  # VAULT_ADDR=https://vault.example.com")
	fmt.Println("  # VAULT_NAMESPACE=admin")
	fmt.Println()
	fmt.Println("  # Set as current profile")
	fmt.Println("  vault profile set prod")
	fmt.Println()
	fmt.Println("  # Execute vault with profile environment")
	fmt.Println("  vault status")
}
