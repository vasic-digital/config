# digital.vasic.config

A Go configuration management library with support for JSON file-based configuration, environment variable binding, validation, and 8 storage-protocol configuration types. Module path: `digital.vasic.config` (Go 1.25). Standalone — no consuming-project context leaks; safe to incorporate at any owning project's root per CONST-051(C).

## Packages

| Package | Purpose |
|---------|---------|
| `pkg/config` | JSON file-based config: `LoadFile` / `SaveFile` / `LoadOrCreate` with default-injection on missing file. Functional-options constructor (`WithFile`, `WithEnvPrefix`). |
| `pkg/env` | Environment-variable binding via struct tags (`env`, `default`, `env_prefix`). Supports string, int, uint, float, bool, `time.Duration`, `[]string`, and nested structs. |
| `pkg/validator` | Composable rules: `Required`, `MinLength`, `Range`, `OneOf`. Multi-error collection via `Validate`. |
| `pkg/storageconfig` | 8 storage-protocol config types: WebDAV, FTP, SFTP, SMB, Google Drive, Dropbox, OneDrive, Git. Each with `New*Config` constructor and `Unmarshal*` decoder; shared `MarshalConfig` encoder. `StorageType` enum exposes `DefaultPort`, `DisplayName`, `SupportsFolders`, `SupportsEncryption`. |

### pkg/config usage

```go
import "digital.vasic.config/pkg/config"

// Load from file
var cfg MyConfig
err := config.LoadFile("config.json", &cfg)

// Save to file
err := config.SaveFile("config.json", cfg)

// Load or create with defaults
defaults := MyConfig{Host: "localhost", Port: 8080}
var cfg MyConfig
err := config.LoadOrCreate("config.json", &cfg, defaults)
```

### pkg/env usage

```go
import "digital.vasic.config/pkg/env"

type ServerConfig struct {
    Host    string        `env:"HOST" default:"localhost"`
    Port    int           `env:"PORT" default:"8080"`
    Debug   bool          `env:"DEBUG" default:"false"`
    Timeout time.Duration `env:"TIMEOUT" default:"30s"`
    Tags    []string      `env:"TAGS" default:"a,b,c"`
}

var cfg ServerConfig
err := env.Load(&cfg)

// With prefix (reads MYAPP_HOST, MYAPP_PORT, etc.)
err := env.LoadWithPrefix("MYAPP_", &cfg)
```

Supported types: `string`, `int*`, `uint*`, `float*`, `bool`, `time.Duration`, `[]string`. Nested structs supported via `env_prefix` tag:

```go
type Config struct {
    DB DatabaseConfig `env_prefix:"DB"`
}

type DatabaseConfig struct {
    Host string `env:"HOST" default:"localhost"`
    Port int    `env:"PORT" default:"5432"`
}
// Reads DB_HOST and DB_PORT
```

### pkg/validator usage

```go
import "digital.vasic.config/pkg/validator"

errs := validator.Validate(
    validator.ValidationField{
        Value: cfg.Host,
        Rules: []validator.Rule{
            validator.Required("host"),
            validator.MinLength("host", 1),
        },
    },
    validator.ValidationField{
        Value: cfg.Port,
        Rules: []validator.Rule{
            validator.Required("port"),
            validator.Range("port", 1, 65535),
        },
    },
    validator.ValidationField{
        Value: cfg.Mode,
        Rules: []validator.Rule{
            validator.OneOf("mode", "debug", "release", "test"),
        },
    },
)
if len(errs) > 0 {
    // Handle validation errors
}
```

### pkg/storageconfig usage

```go
import "digital.vasic.config/pkg/storageconfig"

// Constructors apply protocol-appropriate defaults
sftp := storageconfig.NewSftpConfig("home-server", "fileserver.local")
// sftp.Port == 22, sftp.StrictHostKeyChecking == true, sftp.UseSSL == true

// Encode + decode round-trip
raw, _ := storageconfig.MarshalConfig(sftp)
decoded, _ := storageconfig.UnmarshalSftp(raw)

// Protocol metadata at the StorageType level
fmt.Println(storageconfig.StorageTypeSFTP.DefaultPort())       // 22
fmt.Println(storageconfig.StorageTypeSFTP.DisplayName())       // "SFTP"
fmt.Println(storageconfig.StorageTypeSFTP.SupportsFolders())   // true
fmt.Println(storageconfig.StorageTypeSFTP.SupportsEncryption())// true
```

## Installation

```bash
go get digital.vasic.config
```

## Testing

```bash
# Unit-test floor — all packages, race detector on
go test -race -count=1 ./...

# Round-253 Challenge runner (real OS + real env transport)
bash challenges/config_describe_challenge.sh

# Paired-mutation gate (exit 99 = anti-bluff success)
bash challenges/config_describe_challenge.sh --anti-bluff-mutate
```

## Anti-bluff guarantees (round-253)

> Verbatim 2026-05-19 operator mandate: *"all existing tests and Challenges do work in anti-bluff manner - they MUST confirm that all tested codebase really works as expected! We had been in position that all tests do execute with success and all Challenges as well, but in reality the most of the features does not work and can't be used! This MUST NOT be the case and execution of tests and Challenges MUST guarantee the quality, the completition and full usability by end users of the product!"*

Round-253 ships a runtime-evidence Challenge runner (`challenges/runner/main.go`) and a paired-mutation Challenge gate (`challenges/config_describe_challenge.sh`) that together prove every exported symbol of every package works end-to-end against real OS facilities (real on-disk JSON files in `os.TempDir`, real process-env via `os.Setenv` / `os.Unsetenv`, real reflection-driven struct binding). The runner is bilingual — every code path is exercised with `en`, `sr` (Cyrillic), `ja` (CJK), `ar` (RTL), and `zh-CN` (CJK) inputs from `tests/fixtures/i18n/payloads.json`. CONST-046-aligned: no English-only assumption baked into any code path.

The Challenge gate cross-references every exported symbol from `pkg/{config,env,validator,storageconfig}` to the deep-doc ledger (`docs/test-coverage.md`) and to a runner section. The paired-mutation invocation (`--anti-bluff-mutate`) plants a `LoadFile -> LoadBogus_MUTATED` symbol-rename in a tmp ledger copy and asserts the gate exits 99 — this is the load-bearing proof that the cross-reference gate catches ledger-vs-source drift instead of rubber-stamping a stale ledger.

### What is verified end-to-end (not metadata-only)

- **JSON round-trip** of bilingual structs through real `os.WriteFile` + `os.ReadFile` (Section 1).
- **Defaults injection** when target file does not exist via `LoadOrCreate` (Section 1b).
- **Process-env round-trip** through real `os.Setenv` of non-ASCII values + reflection-driven struct binding (Section 2).
- **Default-tag fallback** when env vars are unset (Section 2b).
- **8-protocol Marshal/Unmarshal** of WebDAV / FTP / SFTP / SMB / GoogleDrive / Dropbox / OneDrive / Git with non-ASCII Name + DefaultPort assertion (Section 3).
- **Every StorageType has DisplayName** coverage (Section 3b).
- **Validator accepts non-ASCII** values without spurious rejection (Section 4).
- **Validator surfaces violations** — empty `Required`, out-of-range `Range`, unknown `OneOf` MUST all return errors (Section 4b — the failure-mode test).

A test that PASSES while the underlying code is broken is the most expensive kind of test in this codebase — the Challenge gate's paired-mutation invariant prevents that failure mode for the cross-reference ledger.

## License

Copyright (c) Milos Vasic. All rights reserved.
