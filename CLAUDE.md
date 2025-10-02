# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Wrapper is a Go CLI application that manages configuration profiles for any binary (e.g., vault, aws, kubectl). It allows wrapping binaries with environment-specific configurations stored as profiles, enabling easy switching between different environments.

## Architecture

This project follows **Domain-Driven Design (DDD)** principles with clean architecture:

### DDD Layers

1. **Domain Layer** (`internal/domain/`)
   - `profile.go` - Profile entity with business logic
   - `binary_name.go`, `profile_name.go` - Value objects
   - `repository.go` - ProfileRepository interface
   - `binary_resolver.go` - BinaryResolver interface
   - Contains core business rules, no external dependencies

2. **Application Layer** (`internal/application/`)
   - `profile_service.go` - Profile management use cases
   - `executor_service.go` - Binary execution orchestration
   - Coordinates domain objects and repositories

3. **Infrastructure Layer** (`internal/infrastructure/`)
   - `filesystem_repository.go` - Profile persistence using filesystem (.env files)
   - `path_binary_resolver.go` - Binary resolution from PATH (excluding wrapper itself)
   - Concrete implementations of domain interfaces

4. **Presentation Layer** (`internal/presentation/`)
   - `profile_commands.go` - CLI commands for profile management
   - `executor_commands.go` - Binary execution entry point
   - `usage.go` - Help/usage display
   - `setup.go` - Dependency injection setup
   - CLI interface using Cobra

### Key Design Patterns

- **Repository Pattern**: Abstract persistence via ProfileRepository interface
- **Dependency Injection**: Dependencies configured in presentation/setup.go
- **Command Pattern**: Profile operations as distinct commands
- **Strategy Pattern**: BinaryResolver for different resolution strategies

### How It Works

1. **Invocation Detection** (main.go:11-24)
   - If invoked as `wrapper <binary>` → extract binary name
   - If invoked via symlink/alias → use executable name as binary name

2. **Command Routing** (main.go:26-34)
   - `profile` command → Profile management (presentation/profile_commands.go)
   - Any other command → Execute real binary (presentation/executor_commands.go)

3. **Profile Storage**
   - Location: `~/.config/wrapper/<binary>/`
   - Format: `.env` files (KEY=VALUE)
   - Current profile: `current.env` symlink
   - Default profile: `.default` file containing profile name

4. **Binary Resolution** (infrastructure/path_binary_resolver.go)
   - Searches PATH for real binary
   - Excludes wrapper itself by comparing resolved paths
   - Returns first matching executable

## Development Commands

### Building
```bash
go build -o wrapper .
```

### Running Directly
```bash
# Using wrapper command
./wrapper vault profile create prod
./wrapper vault status

# Via go run
go run . vault profile list
```

### Testing
```bash
go test ./...
go test ./internal/domain/...
go test ./internal/application/...
```

### Installing
```bash
go install
# Then create shell functions:
# vault() { wrapper vault "$@"; }
```

### Dependencies
```bash
go mod download
go mod tidy
```

## Usage Examples

```bash
# Create a profile
vault profile create prod

# Edit profile (manually)
# Edit ~/.config/wrapper/vault/prod.env
# Add: VAULT_ADDR=https://vault.example.com

# Set current profile
vault profile set prod

# Execute with profile
vault status  # Runs /usr/bin/vault with prod environment

# List profiles
vault profile list

# Set default profile
vault profile default prod
```

## Key Dependencies
- `github.com/spf13/cobra` v1.10.1 - CLI framework
