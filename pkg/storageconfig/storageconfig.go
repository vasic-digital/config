// Package storageconfig provides configuration types for 8 storage protocols.
// This mirrors the Config-KMP Kotlin module's StorageConfig sealed class hierarchy.
package storageconfig

import "encoding/json"

// StorageType represents a storage protocol type.
type StorageType string

const (
	StorageTypeWebDAV      StorageType = "WEBDAV"
	StorageTypeFTP         StorageType = "FTP"
	StorageTypeSFTP        StorageType = "SFTP"
	StorageTypeSMB         StorageType = "SMB"
	StorageTypeGoogleDrive StorageType = "GOOGLE_DRIVE"
	StorageTypeDropbox     StorageType = "DROPBOX"
	StorageTypeOneDrive    StorageType = "ONEDRIVE"
	StorageTypeGit         StorageType = "GIT"
)

// AllStorageTypes returns all supported storage types.
func AllStorageTypes() []StorageType {
	return []StorageType{
		StorageTypeWebDAV, StorageTypeFTP, StorageTypeSFTP, StorageTypeSMB,
		StorageTypeGoogleDrive, StorageTypeDropbox, StorageTypeOneDrive, StorageTypeGit,
	}
}

// DisplayName returns a human-readable name for the storage type.
func (st StorageType) DisplayName() string {
	switch st {
	case StorageTypeWebDAV:
		return "WebDAV"
	case StorageTypeFTP:
		return "FTP"
	case StorageTypeSFTP:
		return "SFTP"
	case StorageTypeSMB:
		return "SMB/CIFS"
	case StorageTypeGoogleDrive:
		return "Google Drive"
	case StorageTypeDropbox:
		return "Dropbox"
	case StorageTypeOneDrive:
		return "OneDrive"
	case StorageTypeGit:
		return "Git"
	default:
		return string(st)
	}
}

// DefaultPort returns the default port for the storage type.
func (st StorageType) DefaultPort() int {
	switch st {
	case StorageTypeWebDAV:
		return 443
	case StorageTypeFTP:
		return 21
	case StorageTypeSFTP:
		return 22
	case StorageTypeSMB:
		return 445
	case StorageTypeGoogleDrive:
		return 443
	case StorageTypeDropbox:
		return 443
	case StorageTypeOneDrive:
		return 443
	case StorageTypeGit:
		return 22
	default:
		return 0
	}
}

// SupportsFolders returns whether the protocol supports folder operations.
func (st StorageType) SupportsFolders() bool {
	return st != StorageTypeFTP
}

// SupportsEncryption returns whether the protocol supports encryption.
func (st StorageType) SupportsEncryption() bool {
	return st != StorageTypeFTP
}

// WebDavAuthType represents WebDAV authentication types.
type WebDavAuthType string

const (
	WebDavAuthBasic  WebDavAuthType = "BASIC"
	WebDavAuthDigest WebDavAuthType = "DIGEST"
	WebDavAuthOAuth  WebDavAuthType = "OAUTH"
	WebDavAuthNone   WebDavAuthType = "NONE"
)

// OneDriveDriveType represents OneDrive drive types.
type OneDriveDriveType string

const (
	OneDriveDriveMe         OneDriveDriveType = "ME"
	OneDriveDriveBusiness   OneDriveDriveType = "BUSINESS"
	OneDriveDriveSharePoint OneDriveDriveType = "SHAREPOINT"
	OneDriveDriveGroup      OneDriveDriveType = "GROUP"
)

// CommonConfig holds fields shared by all storage configurations.
type CommonConfig struct {
	Name      string            `json:"name"`
	Type      StorageType       `json:"storageType"`
	IsEnabled bool              `json:"isEnabled"`
	Priority  int               `json:"priority"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// NewCommonConfig creates a CommonConfig with defaults.
func NewCommonConfig(name string, storageType StorageType) CommonConfig {
	return CommonConfig{
		Name:      name,
		Type:      storageType,
		IsEnabled: true,
		Priority:  100,
	}
}

// WebDavConfig holds WebDAV storage configuration.
type WebDavConfig struct {
	CommonConfig
	URL                string         `json:"url"`
	Username           string         `json:"username"`
	Password           string         `json:"password"`
	AuthenticationType WebDavAuthType `json:"authenticationType"`
	SSLEnabled         bool           `json:"sslEnabled"`
	VerifyCertificate  bool           `json:"verifyCertificate"`
	ConnectionTimeout  int            `json:"connectionTimeout"`
	ReadTimeout        int            `json:"readTimeout"`
}

// NewWebDavConfig creates a WebDavConfig with defaults.
func NewWebDavConfig(name, url, username, password string) *WebDavConfig {
	return &WebDavConfig{
		CommonConfig:       NewCommonConfig(name, StorageTypeWebDAV),
		URL:                url,
		Username:           username,
		Password:           password,
		AuthenticationType: WebDavAuthBasic,
		SSLEnabled:         true,
		VerifyCertificate:  true,
		ConnectionTimeout:  30000,
		ReadTimeout:        60000,
	}
}

// FtpConfig holds FTP storage configuration.
type FtpConfig struct {
	CommonConfig
	Host              string `json:"host"`
	Port              int    `json:"port"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	RootPath          string `json:"rootPath"`
	PassiveMode       bool   `json:"passiveMode"`
	SecureFtp         bool   `json:"secureFtp"`
	Encoding          string `json:"encoding"`
	ConnectionTimeout int    `json:"connectionTimeout"`
}

// NewFtpConfig creates an FtpConfig with defaults.
func NewFtpConfig(name, host, username, password string) *FtpConfig {
	return &FtpConfig{
		CommonConfig:      NewCommonConfig(name, StorageTypeFTP),
		Host:              host,
		Port:              21,
		Username:          username,
		Password:          password,
		RootPath:          "/",
		PassiveMode:       true,
		SecureFtp:         false,
		Encoding:          "UTF-8",
		ConnectionTimeout: 30000,
	}
}

// SftpConfig holds SFTP storage configuration.
type SftpConfig struct {
	CommonConfig
	Host                  string `json:"host"`
	Port                  int    `json:"port"`
	Username              string `json:"username,omitempty"`
	Password              string `json:"password,omitempty"`
	PrivateKeyPath        string `json:"privateKeyPath,omitempty"`
	PrivateKeyPassphrase  string `json:"privateKeyPassphrase,omitempty"`
	KnownHostsPath        string `json:"knownHostsPath,omitempty"`
	StrictHostKeyChecking bool   `json:"strictHostKeyChecking"`
	RootPath              string `json:"rootPath"`
	UseSSL                bool   `json:"useSsl"`
	ConnectionTimeout     int    `json:"connectionTimeout"`
}

// NewSftpConfig creates an SftpConfig with defaults.
func NewSftpConfig(name, host string) *SftpConfig {
	return &SftpConfig{
		CommonConfig:          NewCommonConfig(name, StorageTypeSFTP),
		Host:                  host,
		Port:                  22,
		StrictHostKeyChecking: true,
		RootPath:              "/",
		UseSSL:                true,
		ConnectionTimeout:     30000,
	}
}

// SmbConfig holds SMB/CIFS storage configuration.
type SmbConfig struct {
	CommonConfig
	Host              string `json:"host"`
	Share             string `json:"share"`
	Domain            string `json:"domain,omitempty"`
	Username          string `json:"username"`
	Password          string `json:"password"`
	Path              string `json:"path"`
	Port              int    `json:"port"`
	Encryption        bool   `json:"encryption"`
	Signing           bool   `json:"signing"`
	UseSSL            bool   `json:"useSsl"`
	ConnectionTimeout int    `json:"connectionTimeout"`
}

// NewSmbConfig creates an SmbConfig with defaults.
func NewSmbConfig(name, host, share, username, password string) *SmbConfig {
	return &SmbConfig{
		CommonConfig:      NewCommonConfig(name, StorageTypeSMB),
		Host:              host,
		Share:             share,
		Username:          username,
		Password:          password,
		Path:              "/",
		Port:              445,
		Encryption:        true,
		Signing:           true,
		UseSSL:            false,
		ConnectionTimeout: 30000,
	}
}

// GoogleDriveConfig holds Google Drive storage configuration.
type GoogleDriveConfig struct {
	CommonConfig
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	RefreshToken string `json:"refreshToken,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RootFolderID string `json:"rootFolderId,omitempty"`
	TeamDriveID  string `json:"teamDriveId,omitempty"`
}

// NewGoogleDriveConfig creates a GoogleDriveConfig with defaults.
func NewGoogleDriveConfig(name, clientID, clientSecret string) *GoogleDriveConfig {
	return &GoogleDriveConfig{
		CommonConfig: NewCommonConfig(name, StorageTypeGoogleDrive),
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

// DropboxConfig holds Dropbox storage configuration.
type DropboxConfig struct {
	CommonConfig
	AccessToken  string `json:"accessToken"`
	AppKey       string `json:"appKey"`
	AppSecret    string `json:"appSecret"`
	RefreshToken string `json:"refreshToken,omitempty"`
	RootPath     string `json:"rootPath"`
}

// NewDropboxConfig creates a DropboxConfig with defaults.
func NewDropboxConfig(name, accessToken, appKey, appSecret string) *DropboxConfig {
	return &DropboxConfig{
		CommonConfig: NewCommonConfig(name, StorageTypeDropbox),
		AccessToken:  accessToken,
		AppKey:       appKey,
		AppSecret:    appSecret,
		RootPath:     "",
	}
}

// OneDriveConfig holds OneDrive storage configuration.
type OneDriveConfig struct {
	CommonConfig
	ClientID     string            `json:"clientId"`
	ClientSecret string            `json:"clientSecret"`
	RefreshToken string            `json:"refreshToken,omitempty"`
	AccessToken  string            `json:"accessToken,omitempty"`
	DriveType    OneDriveDriveType `json:"driveType"`
	DriveID      string            `json:"driveId,omitempty"`
	RootFolderID string            `json:"rootFolderId,omitempty"`
}

// NewOneDriveConfig creates an OneDriveConfig with defaults.
func NewOneDriveConfig(name, clientID, clientSecret string) *OneDriveConfig {
	return &OneDriveConfig{
		CommonConfig: NewCommonConfig(name, StorageTypeOneDrive),
		ClientID:     clientID,
		ClientSecret: clientSecret,
		DriveType:    OneDriveDriveMe,
	}
}

// GitConfig holds Git storage configuration.
type GitConfig struct {
	CommonConfig
	RepositoryURL        string `json:"repositoryUrl"`
	Branch               string `json:"branch"`
	Username             string `json:"username,omitempty"`
	Password             string `json:"password,omitempty"`
	PersonalAccessToken  string `json:"personalAccessToken,omitempty"`
	PrivateKeyPath       string `json:"privateKeyPath,omitempty"`
	PrivateKeyPassphrase string `json:"privateKeyPassphrase,omitempty"`
	LocalCachePath       string `json:"localCachePath"`
	AutoSync             bool   `json:"autoSync"`
	CommitAuthorName     string `json:"commitAuthorName"`
	CommitAuthorEmail    string `json:"commitAuthorEmail"`
	ConnectionTimeout    int    `json:"connectionTimeout"`
}

// DefaultCommitAuthorName is the fallback commit author used by
// NewGitConfig when callers supply no override. It is intentionally
// generic so the library does not imply any particular host project.
const DefaultCommitAuthorName = "vasic-config-bot"

// DefaultCommitAuthorEmail is the fallback commit author email used
// by NewGitConfig. Integrators should set their own identity via
// GitConfig.CommitAuthorName / CommitAuthorEmail after construction.
const DefaultCommitAuthorEmail = "noreply@vasic.digital"

// NewGitConfig creates a GitConfig with defaults. The commit-author
// fields default to the generic DefaultCommitAuthorName /
// DefaultCommitAuthorEmail so this library remains project-agnostic;
// callers override the identity to match their own release workflow.
func NewGitConfig(name, repositoryURL, localCachePath string) *GitConfig {
	return &GitConfig{
		CommonConfig:      NewCommonConfig(name, StorageTypeGit),
		RepositoryURL:     repositoryURL,
		Branch:            "main",
		LocalCachePath:    localCachePath,
		AutoSync:          true,
		CommitAuthorName:  DefaultCommitAuthorName,
		CommitAuthorEmail: DefaultCommitAuthorEmail,
		ConnectionTimeout: 30000,
	}
}

// StorageInfo represents lightweight storage metadata.
type StorageInfo struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	StorageType StorageType `json:"storageType"`
	BaseURL     string      `json:"baseUrl,omitempty"`
	IsConnected bool        `json:"isConnected"`
	LastSync    string      `json:"lastSync,omitempty"`
}

// QuotaInfo represents storage quota information.
type QuotaInfo struct {
	UsedBytes      int64   `json:"usedBytes"`
	TotalBytes     int64   `json:"totalBytes"`
	AvailableBytes int64   `json:"availableBytes"`
	UsedPercentage float64 `json:"usedPercentage"`
}

// FileInfo represents file or directory metadata.
type FileInfo struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Size         int64  `json:"size"`
	IsDirectory  bool   `json:"isDirectory"`
	LastModified string `json:"lastModified,omitempty"`
}

// MarshalConfig marshals any config struct to JSON.
func MarshalConfig(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// UnmarshalWebDav unmarshals JSON into a WebDavConfig.
func UnmarshalWebDav(data []byte) (*WebDavConfig, error) {
	var cfg WebDavConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalFtp unmarshals JSON into an FtpConfig.
func UnmarshalFtp(data []byte) (*FtpConfig, error) {
	var cfg FtpConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalSftp unmarshals JSON into an SftpConfig.
func UnmarshalSftp(data []byte) (*SftpConfig, error) {
	var cfg SftpConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalSmb unmarshals JSON into an SmbConfig.
func UnmarshalSmb(data []byte) (*SmbConfig, error) {
	var cfg SmbConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalGoogleDrive unmarshals JSON into a GoogleDriveConfig.
func UnmarshalGoogleDrive(data []byte) (*GoogleDriveConfig, error) {
	var cfg GoogleDriveConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalDropbox unmarshals JSON into a DropboxConfig.
func UnmarshalDropbox(data []byte) (*DropboxConfig, error) {
	var cfg DropboxConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalOneDrive unmarshals JSON into an OneDriveConfig.
func UnmarshalOneDrive(data []byte) (*OneDriveConfig, error) {
	var cfg OneDriveConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}

// UnmarshalGit unmarshals JSON into a GitConfig.
func UnmarshalGit(data []byte) (*GitConfig, error) {
	var cfg GitConfig
	err := json.Unmarshal(data, &cfg)
	return &cfg, err
}
