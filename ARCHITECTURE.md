# Architecture -- Config

## Purpose

Go configuration management library with support for JSON file-based configuration, environment variable binding via struct tags, and composable validation rules with multi-error collection.

## Structure

```
pkg/
  config/      Core config management: JSON file loading, saving, load-or-create with defaults
  env/         Environment variable loading via struct tags (env, default, env_prefix)
  validator/   Validation rules: Required, MinLength, Range, OneOf with multi-error collection
```

## Key Components

- **`config.LoadFile`** / **`config.SaveFile`** -- JSON file serialization and deserialization
- **`config.LoadOrCreate`** -- Load existing config or create with defaults
- **`env.Load`** / **`env.LoadWithPrefix`** -- Populate structs from environment variables with type conversion (string, int, uint, float, bool, duration, string slices)
- **`env` struct tags** -- `env:"VAR_NAME"`, `default:"value"`, `env_prefix:"PREFIX"` for nested structs
- **`validator.Validate`** -- Run composable validation rules and collect all errors
- **`validator.Rule`** functions -- Required, MinLength, Range, OneOf

## Data Flow

```
LoadOrCreate("config.json", &cfg, defaults)
    |
    file exists? -> LoadFile -> unmarshal JSON
    not exists?  -> SaveFile(defaults) -> return defaults

env.Load(&cfg)
    |
    reflect over struct fields -> read env var -> type conversion -> apply default if empty

validator.Validate(fields...)
    |
    for each field: run Rules -> collect errors -> return []error
```

## Dependencies

- `github.com/stretchr/testify` -- Test assertions (only dependency)

## Testing Strategy

Table-driven tests with `testify`. Tests cover JSON loading/saving with defaults, environment variable binding with all supported types, nested struct prefix resolution, and validation rule composition.
