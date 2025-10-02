# Wrapper - Profile Manager for Any Binary

Wrapper is a transparent configuration profile manager for any CLI binary. It allows you to manage multiple environment configurations (profiles) and seamlessly switch between them.

## Features

- 🔧 **Universal**: Works with any binary (vault, aws, kubectl, etc.)
- 📁 **Profile Management**: Create, list, and switch between profiles
- 🔄 **Transparent Execution**: Acts as a pass-through wrapper
- 💾 **Simple Storage**: Profiles stored as `.env` files
- 🎯 **Smart Resolution**: Automatically finds the real binary in PATH

## Installation

```bash
go install github.com/jycamier/wrapper@latest
```

Or build from source:

```bash
git clone https://github.com/jycamier/wrapper.git
cd wrapper
go build -o wrapper .
```

## Quick Start

### 1. Create a shell function

Add to your `.bashrc` or `.zshrc`:

```bash
vault() { wrapper vault "$@"; }
aws() { wrapper aws "$@"; }
kubectl() { wrapper kubectl "$@"; }
```

### 2. Create a profile

```bash
vault profile create prod
```

### 3. Edit the profile

Edit `~/.config/wrapper/vault/prod.env`:

```env
VAULT_ADDR=https://vault.example.com
VAULT_NAMESPACE=admin
VAULT_TOKEN=your-token
```

### 4. Set the profile

```bash
vault profile set prod
```

### 5. Use it!

```bash
vault status
# Executes /usr/bin/vault status with prod environment
```

## Usage

### Profile Commands

```bash
# List profiles
<binary> profile list

# Create a new profile
<binary> profile create <name>

# Set current profile
<binary> profile set <name>

# Get current profile
<binary> profile get

# Set default profile (used when no current profile)
<binary> profile default <name>
```

### Execute with Profile

Any command that is NOT `profile` will execute the real binary with the active profile's environment:

```bash
vault login
aws s3 ls
kubectl get pods
```

## How It Works

1. **Wrapper Detection**: Determines the binary name from how it was invoked
2. **Command Routing**:
   - If command is `profile` → manage profiles
   - Otherwise → execute real binary
3. **Profile Loading**: Loads environment from `current.env` (or default profile)
4. **Binary Execution**: Finds real binary in PATH and executes with loaded environment

## Architecture

Wrapper follows **Domain-Driven Design (DDD)** with clean architecture:

- **Domain Layer**: Core business logic (Profile, Repository interfaces)
- **Application Layer**: Use cases (ProfileService, ExecutorService)
- **Infrastructure Layer**: Filesystem storage, binary resolution
- **Presentation Layer**: CLI commands

See [CLAUDE.md](./CLAUDE.md) for detailed architecture documentation.

## Configuration

Profiles are stored in `~/.config/wrapper/<binary>/`:

```
~/.config/wrapper/
├── vault/
│   ├── prod.env      # Profile configurations
│   ├── dev.env
│   ├── current.env   # Symlink to active profile
│   └── .default      # Default profile name
├── aws/
│   └── ...
└── kubectl/
    └── ...
```

## Examples

### HashiCorp Vault

```bash
# Setup
vault() { wrapper vault "$@"; }

# Create profiles
vault profile create prod
vault profile create dev

# Configure prod
cat > ~/.config/wrapper/vault/prod.env <<EOF
VAULT_ADDR=https://vault.prod.example.com
VAULT_NAMESPACE=admin
EOF

# Configure dev
cat > ~/.config/wrapper/vault/dev.env <<EOF
VAULT_ADDR=https://vault.dev.example.com
VAULT_NAMESPACE=dev
EOF

# Switch between environments
vault profile set prod
vault status

vault profile set dev
vault status
```

### AWS CLI

```bash
# Setup
aws() { wrapper aws "$@"; }

# Create profiles
aws profile create work
aws profile create personal

# Configure
echo "AWS_PROFILE=work" > ~/.config/wrapper/aws/work.env
echo "AWS_PROFILE=personal" > ~/.config/wrapper/aws/personal.env

# Use
aws profile set work
aws s3 ls
```

## Development

```bash
# Build
go build -o wrapper .

# Test
go test ./...

# Run without building
go run . vault profile list
```

## License

See [LICENSE](./LICENSE) file.
