# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Wrapper is a Go CLI application that manages configuration profiles for any binary (e.g., vault, aws, kubectl). It allows wrapping binaries with environment-specific configurations stored as profiles, enabling easy switching between different environments.

## Architecture

This project follows **Domain-Driven Design (DDD)** principles with clean architecture and Cobra CLI framework:

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

4. **Presentation Layer** (`cmd/`)
   - `root.go` - Root command, binary name management, dependency setup, unknown command fallback
   - `profile.go` - Profile parent command
   - `create.go`, `list.go`, `set.go`, `get.go`, `default.go` - Profile subcommands
   - Uses Cobra command structure with standard init() registration

### Command Hierarchy

```
wrapper <binary> [command]
├── profile                    # Profile management (cmd/profile.go)
│   ├── create <name>         # Create profile (cmd/create.go)
│   ├── list                  # List profiles (cmd/list.go)
│   ├── set <name>            # Set current (cmd/set.go)
│   ├── get                   # Get current (cmd/get.go)
│   └── default <name>        # Set default (cmd/default.go)
└── <any other command>       # Execute real binary via fallback
```

### Key Design Patterns

- **Repository Pattern**: Abstract persistence via ProfileRepository interface
- **Dependency Injection**: Dependencies setup in cmd/root.go helper functions
- **Command Pattern**: Cobra commands with init() registration
- **Strategy Pattern**: BinaryResolver for different resolution strategies

### How It Works

1. **Invocation Detection** (main.go:11-18)
   - If invoked as `wrapper` → Use cmd.Execute() (standard Cobra)
   - If invoked via symlink/alias → Use cmd.ExecuteWithBinary()

2. **Command Routing** (main.go:20-22)
   - All args passed to Cobra for routing
   - `profile` command → Cobra profile command tree
   - Unknown commands → Caught by error handler and executed as binary

3. **Binary Name Management** (cmd/root.go:56-60)
   - Extracted from first argument when called as `wrapper <binary>`
   - Passed via ExecuteWithBinary() when called via alias
   - Accessed via GetBinaryName() by all commands

4. **Unknown Command Fallback** (cmd/root.go:63-85)
   - Cobra executes and catches "unknown command" errors
   - Falls back to executing the real binary with all args
   - Preserves exit codes and stdio

5. **Profile Storage**
   - Location: `~/.config/wrapper/<binary>/`
   - Format: `.env` files (KEY=VALUE)
   - Current profile: `current.env` symlink
   - Default profile: `.default` file containing profile name

6. **Binary Resolution** (infrastructure/path_binary_resolver.go)
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
