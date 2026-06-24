# CLAUDE.md

## INHERITED FROM constitution/CLAUDE.md

All rules in `constitution/CLAUDE.md` (and the `constitution/Constitution.md` it references) apply unconditionally. This file's rules below extend them — they MUST NOT weaken any inherited rule. See parent root `CLAUDE.md` §6.AD for the Lava-specific incorporation context (29th §6.L cycle, 2026-05-14) and §6.AD-debt for the implementation-gap inventory. Use `constitution/find_constitution.sh` from the parent project root to resolve the absolute path of the submodule from any nested location.

## INHERITED FROM the Helix Constitution

This module is governed by the Helix Constitution. All rules in the
constitution's `CLAUDE.md` and the `Constitution.md` it references apply
unconditionally. Locate the constitution from any nested depth via its
`find_constitution.sh` helper — do NOT hardcode a path (this module stays
fully decoupled and project-agnostic per §11.4.28).

Canonical reference: https://github.com/HelixDevelopment/HelixConstitution

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

### Acceptance demo for this module

```bash
# TODO — replace with a 10-line real-system demo that loads a JSON config,
# overrides via an env var, validates it, and prints the resolved values.
```
