# CLAUDE.md

This file provides guidance to Claude Code when working with code in this repository.

## Overview

`digital.vasic.config` is a Go configuration management library providing file-based (JSON), environment variable, and programmatic configuration with validation support.

## Commands

```bash
# Build all packages
go build ./...

# Run all tests
go test ./... -count=1

# Run tests with verbose output
go test -v ./... -count=1

# Run tests for a specific package
go test -v ./pkg/config/ -count=1
go test -v ./pkg/env/ -count=1
go test -v ./pkg/validator/ -count=1

# Run a single test
go test -v -run TestName ./pkg/config/
```

## Architecture

The module is organized into three packages:

| Package | Purpose |
|---|---|
| `pkg/config` | Core config management: JSON file loading, saving, load-or-create with defaults |
| `pkg/env` | Environment variable loading via struct tags (`env`, `default`, `env_prefix`) |
| `pkg/validator` | Validation rules: Required, MinLength, Range, OneOf with multi-error collection |

### Package Details

**pkg/config**: `LoadFile` / `SaveFile` / `LoadOrCreate` for JSON config files. `Config` struct with functional options (`WithFile`, `WithEnvPrefix`).

**pkg/env**: `Load` / `LoadWithPrefix` populate structs from env vars. Supports string, int, uint, float, bool, duration, and string slices. Nested structs via `env_prefix` tag.

**pkg/validator**: Composable `Rule` functions validated via `Validate()` which collects all errors.

## Conventions

- Go standard library conventions
- Table-driven tests with `testify/assert` and `testify/require`
- Test files alongside source: `*_test.go`
- Error wrapping with `fmt.Errorf` and `%w`
- Functional options pattern for configuration


## ⚠️ MANDATORY: NO SUDO OR ROOT EXECUTION

**ALL operations MUST run at local user level ONLY.**

This is a PERMANENT and NON-NEGOTIABLE security constraint:

- **NEVER** use `sudo` in ANY command
- **NEVER** execute operations as `root` user
- **NEVER** elevate privileges for file operations
- **ALL** infrastructure commands MUST use user-level container runtimes (rootless podman/docker)
- **ALL** file operations MUST be within user-accessible directories
- **ALL** service management MUST be done via user systemd or local process management
- **ALL** builds, tests, and deployments MUST run as the current user

### Why This Matters
- **Security**: Prevents accidental system-wide damage
- **Reproducibility**: User-level operations are portable across systems
- **Safety**: Limits blast radius of any issues
- **Best Practice**: Modern container workflows are rootless by design

### When You See SUDO
If any script or command suggests using `sudo`:
1. STOP immediately
2. Find a user-level alternative
3. Use rootless container runtimes
4. Modify commands to work within user permissions

**VIOLATION OF THIS CONSTRAINT IS STRICTLY PROHIBITED.**

