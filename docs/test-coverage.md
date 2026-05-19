# Test Coverage — digital.vasic.config (round-253)

> Verbatim 2026-05-19 operator mandate: *"all existing tests and Challenges do work in anti-bluff manner - they MUST confirm that all tested codebase really works as expected! We had been in position that all tests do execute with success and all Challenges as well, but in reality the most of the features does not work and can't be used! This MUST NOT be the case and execution of tests and Challenges MUST guarantee the quality, the completition and full usability by end users of the product!"*

CONST-050(B) symbol-to-test ledger. Every exported symbol in `pkg/{config,env,validator,storageconfig}` is cross-referenced to the test name(s) that exercise it AND to the round-253 Challenge runner section that exercises it against real OS facilities (`os.TempDir` files + `os.Setenv` process env). No metadata-only PASS — every entry below names the production code path and the runtime evidence channel that proves it works.

## Anti-bluff posture (round-253)

- **Real on-disk transport.** `challenges/runner/main.go` Section 1 + 1b creates files in `os.MkdirTemp` and round-trips bilingual (en/sr/ja/ar/zh-CN) JSON through `config.SaveFile` and `config.LoadFile`. Byte-length is captured per locale in the PASS line.
- **Real process env transport.** Section 2 + 2b uses `os.Setenv` / `os.Unsetenv` and asserts both `env` tag-driven and `default` tag-driven paths survive the round-trip. No in-memory env stub.
- **Real reflection-driven binding.** `env.Load` walks the target struct via `reflect`. The runner uses public `EnvStruct` types with tagged fields; the bytes returned MUST equal the bytes set into the process env.
- **Every protocol exercised.** Section 3 runs `MarshalConfig` + every `Unmarshal*` of the 8 storage protocols (WebDAV / FTP / SFTP / SMB / GoogleDrive / Dropbox / OneDrive / Git) with a non-ASCII Name, asserting byte preservation AND `DefaultPort()` correctness.
- **Validator must surface violations.** Section 4b plants 3 deliberate violations (empty string for `Required`, out-of-range int for `Range`, unknown value for `OneOf`) and asserts ALL 3 errors are returned. A validator that silently swallows violations is a bluff under Article XI §11.9.
- **Paired mutation.** `config_describe_challenge.sh --anti-bluff-mutate` plants a `LoadFile -> LoadBogus_MUTATED` rename in a tmp ledger copy and asserts the gate exits 99. Proves the cross-reference gate catches ledger-vs-source drift instead of rubber-stamping it.

## pkg/config

| Exported symbol | Unit-test coverage | Runner section |
|-----------------|--------------------|----------------|
| `type Config` | `TestNew`, `TestNewWithOptions` (config_test.go) | n/a (constructor exercised indirectly via runner Section 1) |
| `type Loader` (interface) | compile-time only — interface | n/a |
| `type Option` | `TestNewWithOptions` (config_test.go) | n/a |
| `type Validator` (interface) | compile-time only — interface | n/a |
| `func New(opts ...Option) *Config` | `TestNew`, `TestNewWithOptions` | n/a |
| `func WithFile(path string) Option` | `TestNewWithOptions` | n/a |
| `func WithEnvPrefix(prefix string) Option` | `TestNewWithOptions` | n/a |
| `func LoadFile(path string, target interface{}) error` | `TestLoadFile`, `TestLoadFile_NotFound`, `TestLoadFile_InvalidJSON`, `TestLoadFile_EmptyFile`, `TestLoadFile_MalformedJSON`, `TestLoadFile_UnicodeValues`, `TestLoadFile_TypeMismatch`, `TestLoadFile_ExtremelyNestedConfig`, `TestLoadFile_MissingRequiredFields`, `TestRoundTrip` | Section 1 (round-trip per locale) |
| `func SaveFile(path string, config interface{}) error` | `TestSaveFile`, `TestSaveFile_CreatesDirs`, `TestSaveFile_UnmarshalableValue`, `TestRoundTrip` | Section 1 (writes to `os.TempDir`) |
| `func LoadOrCreate(path string, target interface{}, defaults interface{}) error` | `TestLoadOrCreate_CreatesNew`, `TestLoadOrCreate_LoadsExisting`, `TestLoadOrCreate_InvalidDirectory` | Section 1b (missing-file -> defaults applied) |

## pkg/env

| Exported symbol | Unit-test coverage | Runner section |
|-----------------|--------------------|----------------|
| `func Load(target interface{}) error` | `TestLoad_Defaults`, `TestLoad_EnvOverrides`, `TestLoad_Duration`, `TestLoad_DurationFromEnv`, `TestLoad_Slice`, `TestLoad_SliceFromEnv`, `TestLoad_NestedStruct`, `TestLoad_NestedStructFromEnv`, `TestLoad_NoTags`, `TestLoad_NonPointerError`, `TestLoad_NonStructPointerError`, `TestLoad_InvalidIntValue`, `TestLoad_InvalidUintValue`, `TestLoad_InvalidFloatValue`, `TestLoad_InvalidBoolValue`, `TestLoad_InvalidDurationValue` | n/a (Section 2 exercises LoadWithPrefix which shares the implementation) |
| `func LoadWithPrefix(prefix string, target interface{}) error` | `TestLoadWithPrefix`, `TestEnvLoad_EmptyEnvVariable`, `TestEnvLoad_UnicodeEnvVariable`, `TestEnvLoad_WhitespaceEnvVariable`, `TestEnvLoad_UnexportedFields` | Section 2 (real `os.Setenv` round-trip per locale), Section 2b (defaults fallback) |

## pkg/validator

| Exported symbol | Unit-test coverage | Runner section |
|-----------------|--------------------|----------------|
| `type Rule` (function type) | every Test* uses Rule values | implicit in Section 4 |
| `type ValidationField` | every Test* | Section 4 |
| `func Required(fieldName string) Rule` | `TestRequired_NonZero`, `TestRequired_Zero`, `TestRequired_NonZeroInt`, `TestRequired_ZeroInt`, `TestValidator_Required_ZeroValues`, `TestValidator_Required_NilPanics` | Section 4 (non-ASCII Name) + Section 4b (empty-string violation) |
| `func MinLength(fieldName string, min int) Rule` | `TestMinLength_TooShort`, `TestMinLength_Valid`, `TestMinLength_ExactLength`, `TestMinLength_NotString`, `TestValidator_MinLength_EdgeCases` | Section 4 (rune-length aware) |
| `func Range(fieldName string, min, max int) Rule` | `TestRange_TooLow`, `TestRange_TooHigh`, `TestRange_Valid`, `TestRange_AtMin`, `TestRange_AtMax`, `TestRange_NotInt`, `TestValidator_Range_EdgeCases` | Section 4 + 4b (out-of-range violation) |
| `func OneOf(fieldName string, allowed ...string) Rule` | `TestOneOf_Valid`, `TestOneOf_Invalid`, `TestOneOf_NotString`, `TestValidator_OneOf_EdgeCases` | Section 4b (unknown-value violation) |
| `func Validate(fields ...ValidationField) []error` | `TestValidate_AllPass`, `TestValidate_AllFail`, `TestValidate_SomeFailures`, `TestValidate_Empty`, `TestValidate_MultipleRulesPerField`, `TestValidator_Validate_EmptyFields`, `TestValidator_Validate_NoErrors`, `TestValidator_Validate_MultipleErrors` | Section 4 + Section 4b (3 violations surfaced) |

## pkg/storageconfig

### Types

| Type | Unit-test coverage | Runner section |
|------|--------------------|----------------|
| `type StorageType` | `TestStorageTypeDefaultPorts`, `TestStorageTypeDisplayNames`, `TestStorageTypeSupportsFolders`, `TestStorageTypeSupportsEncryption`, `TestUnknownStorageType` | Section 3 + Section 3b |
| `type WebDavAuthType` | `TestWebDavAuthTypes` | indirect via Section 3 WebDAV |
| `type OneDriveDriveType` | `TestOneDriveDriveTypes` | indirect via Section 3 OneDrive |
| `type CommonConfig` | `TestCommonConfigDefaults`, `TestCommonConfigWithMetadata` | every Section 3 case carries a CommonConfig |
| `type WebDavConfig` | `TestNewWebDavConfig`, `TestWebDavConfigJSON` | Section 3 (Marshal+Unmarshal non-ASCII) |
| `type FtpConfig` | `TestNewFtpConfig`, `TestFtpConfigJSON` | Section 3 |
| `type SftpConfig` | `TestNewSftpConfig`, `TestSftpConfigJSON` | Section 3 |
| `type SmbConfig` | `TestNewSmbConfig`, `TestSmbConfigJSON` | Section 3 |
| `type GoogleDriveConfig` | `TestNewGoogleDriveConfig`, `TestGoogleDriveConfigJSON` | Section 3 |
| `type DropboxConfig` | `TestNewDropboxConfig`, `TestDropboxConfigJSON` | Section 3 |
| `type OneDriveConfig` | `TestNewOneDriveConfig`, `TestOneDriveConfigJSON` | Section 3 |
| `type GitConfig` | `TestNewGitConfig`, `TestGitConfigJSON` | Section 3 |
| `type StorageInfo` | `TestStorageInfo` | n/a (passive type) |
| `type QuotaInfo` | `TestQuotaInfo` | n/a (passive type) |
| `type FileInfo` | `TestFileInfo`, `TestFileInfoDirectory` | n/a (passive type) |

### Constants

| Constant | Unit-test coverage |
|----------|--------------------|
| `DefaultCommitAuthorName` | `TestNewGitConfig` |
| `DefaultCommitAuthorEmail` | `TestNewGitConfig` |

### Functions / methods

| Symbol | Unit-test coverage | Runner section |
|--------|--------------------|----------------|
| `func AllStorageTypes() []StorageType` | `TestAllStorageTypes` | Section 3b (iterates AllStorageTypes) |
| `(StorageType).DisplayName() string` | `TestStorageTypeDisplayNames`, `TestUnknownStorageType` | Section 3b |
| `(StorageType).DefaultPort() int` | `TestStorageTypeDefaultPorts` | Section 3 (asserts expectPort matches) |
| `(StorageType).SupportsFolders() bool` | `TestStorageTypeSupportsFolders` | Section 3 (invariant) |
| `(StorageType).SupportsEncryption() bool` | `TestStorageTypeSupportsEncryption` | Section 3 (invariant) |
| `func NewCommonConfig(name string, storageType StorageType) CommonConfig` | `TestCommonConfigDefaults` | indirect via every NewXConfig |
| `func NewWebDavConfig(name, url, username, password string) *WebDavConfig` | `TestNewWebDavConfig` | Section 3 |
| `func NewFtpConfig(name, host, username, password string) *FtpConfig` | `TestNewFtpConfig` | Section 3 |
| `func NewSftpConfig(name, host string) *SftpConfig` | `TestNewSftpConfig` | Section 3 |
| `func NewSmbConfig(name, host, share, username, password string) *SmbConfig` | `TestNewSmbConfig` | Section 3 |
| `func NewGoogleDriveConfig(name, clientID, clientSecret string) *GoogleDriveConfig` | `TestNewGoogleDriveConfig` | Section 3 |
| `func NewDropboxConfig(name, accessToken, appKey, appSecret string) *DropboxConfig` | `TestNewDropboxConfig` | Section 3 |
| `func NewOneDriveConfig(name, clientID, clientSecret string) *OneDriveConfig` | `TestNewOneDriveConfig` | Section 3 |
| `func NewGitConfig(name, repositoryURL, localCachePath string) *GitConfig` | `TestNewGitConfig` | Section 3 |
| `func MarshalConfig(v interface{}) ([]byte, error)` | every `Test*ConfigJSON` | Section 3 |
| `func UnmarshalWebDav(data []byte) (*WebDavConfig, error)` | `TestWebDavConfigJSON` | Section 3 |
| `func UnmarshalFtp(data []byte) (*FtpConfig, error)` | `TestFtpConfigJSON` | Section 3 |
| `func UnmarshalSftp(data []byte) (*SftpConfig, error)` | `TestSftpConfigJSON` | Section 3 |
| `func UnmarshalSmb(data []byte) (*SmbConfig, error)` | `TestSmbConfigJSON` | Section 3 |
| `func UnmarshalGoogleDrive(data []byte) (*GoogleDriveConfig, error)` | `TestGoogleDriveConfigJSON` | Section 3 |
| `func UnmarshalDropbox(data []byte) (*DropboxConfig, error)` | `TestDropboxConfigJSON` | Section 3 |
| `func UnmarshalOneDrive(data []byte) (*OneDriveConfig, error)` | `TestOneDriveConfigJSON` | Section 3 |
| `func UnmarshalGit(data []byte) (*GitConfig, error)` | `TestGitConfigJSON` | Section 3 |

## Verification

```bash
# Unit-test floor (testify) — all packages
go test -race -count=1 ./...

# Round-253 Challenge runner + paired-mutation
bash challenges/config_describe_challenge.sh                  # exit 0
bash challenges/config_describe_challenge.sh --anti-bluff-mutate  # exit 99
```

The paired-mutation invocation is the load-bearing proof that the cross-reference gate (Section 2 of `config_describe_challenge.sh`) catches ledger-vs-source drift — a ledger that silently lists nothing would PASS Section 2 vacuously without the mutation check. Exit 99 means the gate FAILED on the planted mutation, which is the desired anti-bluff behaviour.
