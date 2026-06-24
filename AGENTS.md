# AGENTS.md

## INHERITED FROM constitution/AGENTS.md

All rules in `constitution/AGENTS.md` (and the `constitution/Constitution.md` it references) apply unconditionally. This file's rules below extend them — they MUST NOT weaken any inherited rule. See parent root `CLAUDE.md` §6.AD for the Lava-specific incorporation context (29th §6.L cycle, 2026-05-14) and §6.AD-debt for the implementation-gap inventory. Use `constitution/find_constitution.sh` from the parent project root to resolve the absolute path of the submodule from any nested location.

## INHERITED FROM the Helix Constitution

This module is governed by the Helix Constitution. All rules in the
constitution's `AGENTS.md` and the `Constitution.md` it references apply
unconditionally. Locate the constitution from any nested depth via its
`find_constitution.sh` helper — do NOT hardcode a path (this module stays
fully decoupled and project-agnostic per §11.4.28).

Canonical reference: https://github.com/HelixDevelopment/HelixConstitution

Guidelines for AI agents working on this codebase.

## Project Context

This is the `digital.vasic.config` Go module -- a configuration management library. It is a standalone module with no application binary; it provides packages for import by other projects.

## Development Guidelines

1. **Do not add a `main.go`** -- this is a library module, not an application.
2. **All public functions must have doc comments** following Go conventions.
3. **Tests use `testify`** (`assert` and `require`) -- do not introduce other test frameworks.
4. **Error messages** should use lowercase and wrap with `%w` where applicable.
5. **No external dependencies** beyond `testify` for testing -- keep the dependency footprint minimal.

## Testing

Run all tests before submitting changes:

```bash
go test ./... -count=1
```

Every new public function must have corresponding test coverage.

## Package Boundaries

- `pkg/config` -- file I/O only, no env var logic
- `pkg/env` -- env var loading only, no file I/O
- `pkg/validator` -- pure validation logic, no I/O of any kind

Keep these boundaries clean. Cross-package imports within this module should be avoided.
