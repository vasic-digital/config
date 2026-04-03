package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"digital.vasic.config/pkg/config"
	"digital.vasic.config/pkg/env"
	"digital.vasic.config/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- LoadFile edge cases ---

func TestLoadFile_MalformedJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
	}{
		{"trailing comma", `{"host": "localhost",}`},
		{"single quotes", `{'host': 'localhost'}`},
		{"unclosed brace", `{"host": "localhost"`},
		{"unclosed string", `{"host": "localhost`},
		{"bare words", `{host: localhost}`},
		{"binary garbage", "\x00\x01\x02\x03\xff\xfe"},
		{"empty object with trailing data", `{}garbage`},
		{"array instead of object", `[1, 2, 3]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			path := filepath.Join(dir, "config.json")
			err := os.WriteFile(path, []byte(tt.content), 0644)
			require.NoError(t, err)

			type cfg struct {
				Host string `json:"host"`
			}
			var c cfg
			err = config.LoadFile(path, &c)
			assert.Error(t, err)
		})
	}
}

func TestLoadFile_EmptyFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	err := os.WriteFile(path, []byte(""), 0644)
	require.NoError(t, err)

	type cfg struct {
		Host string `json:"host"`
	}
	var c cfg
	err = config.LoadFile(path, &c)
	assert.Error(t, err, "empty file should fail to parse")
}

func TestLoadFile_TypeMismatch(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Port is a string in JSON but int in struct
	content := `{"host": "localhost", "port": "not_a_number"}`
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	type cfg struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	var c cfg
	err = config.LoadFile(path, &c)
	assert.Error(t, err, "string where int expected should fail")
}

func TestLoadFile_UnicodeValues(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	content := `{"host": "\u4e2d\u6587\u670d\u52a1\u5668", "name": "emoji server"}`
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	type cfg struct {
		Host string `json:"host"`
		Name string `json:"name"`
	}
	var c cfg
	err = config.LoadFile(path, &c)
	require.NoError(t, err)
	assert.Contains(t, c.Host, "\u4e2d\u6587")
}

func TestLoadFile_ExtremelyNestedConfig(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Build a deeply nested JSON object
	content := `{"a":{"b":{"c":{"d":{"e":{"f":{"g":{"h":"deep"}}}}}}}}`
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	var c map[string]interface{}
	err = config.LoadFile(path, &c)
	require.NoError(t, err)
	assert.NotNil(t, c["a"])
}

func TestLoadFile_MissingRequiredFields(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// JSON has no fields matching the struct
	content := `{"unrelated_field": "value"}`
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)

	type cfg struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	}
	var c cfg
	// JSON unmarshal does not error on missing fields -- they stay zero-valued
	err = config.LoadFile(path, &c)
	require.NoError(t, err)
	assert.Empty(t, c.Host)
	assert.Zero(t, c.Port)
}

func TestSaveFile_UnmarshalableValue(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// Channels cannot be marshaled
	type badCfg struct {
		Ch chan int `json:"ch"`
	}
	err := config.SaveFile(path, badCfg{Ch: make(chan int)})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal config")
}

func TestLoadOrCreate_InvalidDirectory(t *testing.T) {
	t.Parallel()
	// Use a path under /proc which is read-only
	path := "/proc/nonexistent/impossible/config.json"

	type cfg struct {
		Host string `json:"host"`
	}
	defaults := cfg{Host: "default"}
	var c cfg
	err := config.LoadOrCreate(path, &c, defaults)
	assert.Error(t, err, "should fail to create file in read-only filesystem")
}

// --- env.Load edge cases ---

func TestEnvLoad_EmptyEnvVariable(t *testing.T) {
	// An empty env var should override the default
	os.Setenv("EDGE_HOST", "")
	defer os.Unsetenv("EDGE_HOST")

	type cfg struct {
		Host string `env:"EDGE_HOST" default:"fallback"`
	}
	var c cfg
	err := env.Load(&c)
	require.NoError(t, err)
	// Empty string env var means use default (env is "" which is treated as unset)
	assert.Equal(t, "fallback", c.Host)
}

func TestEnvLoad_WhitespaceEnvVariable(t *testing.T) {
	os.Setenv("EDGE_WS_HOST", "   ")
	defer os.Unsetenv("EDGE_WS_HOST")

	type cfg struct {
		Host string `env:"EDGE_WS_HOST" default:"fallback"`
	}
	var c cfg
	err := env.Load(&c)
	require.NoError(t, err)
	assert.Equal(t, "   ", c.Host, "whitespace should be preserved")
}

func TestEnvLoad_UnicodeEnvVariable(t *testing.T) {
	os.Setenv("EDGE_UNI_NAME", "\u4e2d\u6587\u540d\u79f0")
	defer os.Unsetenv("EDGE_UNI_NAME")

	type cfg struct {
		Name string `env:"EDGE_UNI_NAME" default:"default"`
	}
	var c cfg
	err := env.Load(&c)
	require.NoError(t, err)
	assert.Equal(t, "\u4e2d\u6587\u540d\u79f0", c.Name)
}

func TestEnvLoad_UnexportedFields(t *testing.T) {
	t.Parallel()
	type cfg struct {
		Public  string `env:"EDGE_PUB" default:"visible"`
		private string `env:"EDGE_PRIV" default:"hidden"` //nolint:unused
	}
	var c cfg
	err := env.Load(&c)
	require.NoError(t, err)
	assert.Equal(t, "visible", c.Public)
}

// --- validator edge cases ---

func TestValidator_Required_ZeroValues(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{"empty string", "", true},
		{"zero int", 0, true},
		{"false bool", false, true},
		// Note: nil panics in reflect.ValueOf().IsZero() - tested separately
		{"non-empty string", "hello", false},
		{"non-zero int", 42, false},
		{"true bool", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := validator.Required("field")
			err := rule(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Required_NilPanics(t *testing.T) {
	t.Parallel()
	rule := validator.Required("field")
	// reflect.ValueOf(nil).IsZero() panics -- this documents the behavior
	assert.Panics(t, func() {
		_ = rule(nil)
	})
}

func TestValidator_MinLength_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   interface{}
		min     int
		wantErr bool
	}{
		{"exact minimum", "abc", 3, false},
		{"one below minimum", "ab", 3, true},
		{"empty string min 0", "", 0, false},
		{"non-string type", 12345, 3, true},
		{"unicode string length", "\u4e2d\u6587\u5b57", 3, false},
		{"min 0 always passes", "anything", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := validator.MinLength("field", tt.min)
			err := rule(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Range_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   interface{}
		min     int
		max     int
		wantErr bool
	}{
		{"at minimum", 1, 1, 10, false},
		{"at maximum", 10, 1, 10, false},
		{"below minimum", 0, 1, 10, true},
		{"above maximum", 11, 1, 10, true},
		{"negative range", -5, -10, -1, false},
		{"zero range", 5, 5, 5, false},
		{"string type", "hello", 1, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := validator.Range("field", tt.min, tt.max)
			err := rule(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_OneOf_EdgeCases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		value   interface{}
		allowed []string
		wantErr bool
	}{
		{"exact match", "dev", []string{"dev", "prod"}, false},
		{"no match", "staging", []string{"dev", "prod"}, true},
		{"empty allowed list", "anything", []string{}, true},
		{"empty string in allowed", "", []string{""}, false},
		{"non-string type", 42, []string{"42"}, true},
		{"case sensitive", "Dev", []string{"dev"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rule := validator.OneOf("field", tt.allowed...)
			err := rule(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidator_Validate_MultipleErrors(t *testing.T) {
	t.Parallel()
	errs := validator.Validate(
		validator.ValidationField{
			Value: "",
			Rules: []validator.Rule{
				validator.Required("name"),
				validator.MinLength("name", 3),
			},
		},
		validator.ValidationField{
			Value: 100,
			Rules: []validator.Rule{
				validator.Range("port", 1, 65535),
			},
		},
	)

	// First field should have 2 errors, second should have 0
	assert.Len(t, errs, 2)
}

func TestValidator_Validate_NoErrors(t *testing.T) {
	t.Parallel()
	errs := validator.Validate(
		validator.ValidationField{
			Value: "valid",
			Rules: []validator.Rule{
				validator.Required("name"),
				validator.MinLength("name", 1),
			},
		},
	)
	assert.Empty(t, errs)
}

func TestValidator_Validate_EmptyFields(t *testing.T) {
	t.Parallel()
	errs := validator.Validate()
	assert.Empty(t, errs)
}
