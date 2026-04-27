package storageconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllStorageTypes(t *testing.T) {
	types := AllStorageTypes()
	assert.Len(t, types, 8)
}

func TestStorageTypeDisplayNames(t *testing.T) {
	tests := []struct {
		st   StorageType
		want string
	}{
		{StorageTypeWebDAV, "WebDAV"},
		{StorageTypeFTP, "FTP"},
		{StorageTypeSFTP, "SFTP"},
		{StorageTypeSMB, "SMB/CIFS"},
		{StorageTypeGoogleDrive, "Google Drive"},
		{StorageTypeDropbox, "Dropbox"},
		{StorageTypeOneDrive, "OneDrive"},
		{StorageTypeGit, "Git"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.st.DisplayName())
	}
}

func TestStorageTypeDefaultPorts(t *testing.T) {
	assert.Equal(t, 443, StorageTypeWebDAV.DefaultPort())
	assert.Equal(t, 21, StorageTypeFTP.DefaultPort())
	assert.Equal(t, 22, StorageTypeSFTP.DefaultPort())
	assert.Equal(t, 445, StorageTypeSMB.DefaultPort())
	assert.Equal(t, 443, StorageTypeGoogleDrive.DefaultPort())
	assert.Equal(t, 443, StorageTypeDropbox.DefaultPort())
	assert.Equal(t, 443, StorageTypeOneDrive.DefaultPort())
	assert.Equal(t, 22, StorageTypeGit.DefaultPort())
}

func TestStorageTypeSupportsFolders(t *testing.T) {
	assert.True(t, StorageTypeWebDAV.SupportsFolders())
	assert.False(t, StorageTypeFTP.SupportsFolders())
	assert.True(t, StorageTypeSFTP.SupportsFolders())
	assert.True(t, StorageTypeSMB.SupportsFolders())
}

func TestStorageTypeSupportsEncryption(t *testing.T) {
	assert.True(t, StorageTypeWebDAV.SupportsEncryption())
	assert.False(t, StorageTypeFTP.SupportsEncryption())
	assert.True(t, StorageTypeSFTP.SupportsEncryption())
}

func TestUnknownStorageType(t *testing.T) {
	unknown := StorageType("UNKNOWN")
	assert.Equal(t, "UNKNOWN", unknown.DisplayName())
	assert.Equal(t, 0, unknown.DefaultPort())
	assert.True(t, unknown.SupportsFolders())
	assert.True(t, unknown.SupportsEncryption())
}

func TestNewWebDavConfig(t *testing.T) {
	cfg := NewWebDavConfig("My WebDAV", "https://dav.example.com", "user", "pass")
	assert.Equal(t, "My WebDAV", cfg.Name)
	assert.Equal(t, "https://dav.example.com", cfg.URL)
	assert.Equal(t, "user", cfg.Username)
	assert.Equal(t, "pass", cfg.Password)
	assert.Equal(t, StorageTypeWebDAV, cfg.Type)
	assert.True(t, cfg.IsEnabled)
	assert.Equal(t, 100, cfg.Priority)
	assert.Equal(t, WebDavAuthBasic, cfg.AuthenticationType)
	assert.True(t, cfg.SSLEnabled)
	assert.True(t, cfg.VerifyCertificate)
	assert.Equal(t, 30000, cfg.ConnectionTimeout)
	assert.Equal(t, 60000, cfg.ReadTimeout)
}

func TestNewFtpConfig(t *testing.T) {
	cfg := NewFtpConfig("My FTP", "ftp.example.com", "user", "pass")
	assert.Equal(t, StorageTypeFTP, cfg.Type)
	assert.Equal(t, 21, cfg.Port)
	assert.Equal(t, "/", cfg.RootPath)
	assert.True(t, cfg.PassiveMode)
	assert.False(t, cfg.SecureFtp)
	assert.Equal(t, "UTF-8", cfg.Encoding)
}

func TestNewSftpConfig(t *testing.T) {
	cfg := NewSftpConfig("My SFTP", "sftp.example.com")
	assert.Equal(t, StorageTypeSFTP, cfg.Type)
	assert.Equal(t, 22, cfg.Port)
	assert.Empty(t, cfg.Username)
	assert.Empty(t, cfg.Password)
	assert.True(t, cfg.StrictHostKeyChecking)
	assert.Equal(t, "/", cfg.RootPath)
	assert.True(t, cfg.UseSSL)
}

func TestNewSmbConfig(t *testing.T) {
	cfg := NewSmbConfig("My SMB", "smb.example.com", "documents", "user", "pass")
	assert.Equal(t, StorageTypeSMB, cfg.Type)
	assert.Equal(t, 445, cfg.Port)
	assert.Equal(t, "documents", cfg.Share)
	assert.Empty(t, cfg.Domain)
	assert.True(t, cfg.Encryption)
	assert.True(t, cfg.Signing)
	assert.False(t, cfg.UseSSL)
}

func TestNewGoogleDriveConfig(t *testing.T) {
	cfg := NewGoogleDriveConfig("GDrive", "id123", "secret456")
	assert.Equal(t, StorageTypeGoogleDrive, cfg.Type)
	assert.Equal(t, "id123", cfg.ClientID)
	assert.Equal(t, "secret456", cfg.ClientSecret)
	assert.Empty(t, cfg.RefreshToken)
	assert.Empty(t, cfg.RootFolderID)
}

func TestNewDropboxConfig(t *testing.T) {
	cfg := NewDropboxConfig("Dropbox", "token", "key", "secret")
	assert.Equal(t, StorageTypeDropbox, cfg.Type)
	assert.Equal(t, "token", cfg.AccessToken)
	assert.Equal(t, "key", cfg.AppKey)
	assert.Equal(t, "", cfg.RootPath)
}

func TestNewOneDriveConfig(t *testing.T) {
	cfg := NewOneDriveConfig("OneDrive", "id", "secret")
	assert.Equal(t, StorageTypeOneDrive, cfg.Type)
	assert.Equal(t, OneDriveDriveMe, cfg.DriveType)
	assert.Empty(t, cfg.DriveID)
}

func TestNewGitConfig(t *testing.T) {
	cfg := NewGitConfig("GitHub", "https://github.com/user/repo.git", "/tmp/cache")
	assert.Equal(t, StorageTypeGit, cfg.Type)
	assert.Equal(t, "main", cfg.Branch)
	assert.True(t, cfg.AutoSync)
	// The library ships generic defaults so it stays project-
	// agnostic. Integrators override these after construction to
	// match their own release workflow.
	assert.Equal(t, DefaultCommitAuthorName, cfg.CommitAuthorName)
	assert.Equal(t, DefaultCommitAuthorEmail, cfg.CommitAuthorEmail)
	assert.Equal(t, 30000, cfg.ConnectionTimeout)
}

func TestWebDavConfigJSON(t *testing.T) {
	cfg := NewWebDavConfig("Test", "https://dav.test.com", "u", "p")
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalWebDav(data)
	require.NoError(t, err)
	assert.Equal(t, cfg.Name, restored.Name)
	assert.Equal(t, cfg.URL, restored.URL)
	assert.Equal(t, cfg.Type, restored.Type)
	assert.Equal(t, cfg.AuthenticationType, restored.AuthenticationType)
}

func TestFtpConfigJSON(t *testing.T) {
	cfg := NewFtpConfig("Test", "ftp.test.com", "u", "p")
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalFtp(data)
	require.NoError(t, err)
	assert.Equal(t, cfg.Name, restored.Name)
	assert.Equal(t, cfg.Host, restored.Host)
	assert.Equal(t, cfg.Port, restored.Port)
}

func TestSftpConfigJSON(t *testing.T) {
	cfg := NewSftpConfig("Test", "sftp.test.com")
	cfg.Username = "user"
	cfg.PrivateKeyPath = "/home/.ssh/id_rsa"
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalSftp(data)
	require.NoError(t, err)
	assert.Equal(t, "user", restored.Username)
	assert.Equal(t, "/home/.ssh/id_rsa", restored.PrivateKeyPath)
}

func TestSmbConfigJSON(t *testing.T) {
	cfg := NewSmbConfig("Test", "smb.test.com", "share", "u", "p")
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalSmb(data)
	require.NoError(t, err)
	assert.Equal(t, "share", restored.Share)
	assert.Equal(t, 445, restored.Port)
}

func TestGoogleDriveConfigJSON(t *testing.T) {
	cfg := NewGoogleDriveConfig("Test", "cid", "csecret")
	cfg.RefreshToken = "refresh123"
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalGoogleDrive(data)
	require.NoError(t, err)
	assert.Equal(t, "refresh123", restored.RefreshToken)
}

func TestDropboxConfigJSON(t *testing.T) {
	cfg := NewDropboxConfig("Test", "tok", "key", "sec")
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalDropbox(data)
	require.NoError(t, err)
	assert.Equal(t, "tok", restored.AccessToken)
}

func TestOneDriveConfigJSON(t *testing.T) {
	cfg := NewOneDriveConfig("Test", "cid", "csecret")
	cfg.DriveType = OneDriveDriveBusiness
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalOneDrive(data)
	require.NoError(t, err)
	assert.Equal(t, OneDriveDriveBusiness, restored.DriveType)
}

func TestGitConfigJSON(t *testing.T) {
	cfg := NewGitConfig("Test", "https://github.com/u/r.git", "/tmp/c")
	cfg.PersonalAccessToken = "ghp_xxx"
	data, err := MarshalConfig(cfg)
	require.NoError(t, err)

	restored, err := UnmarshalGit(data)
	require.NoError(t, err)
	assert.Equal(t, "ghp_xxx", restored.PersonalAccessToken)
	assert.Equal(t, "main", restored.Branch)
}

func TestStorageInfo(t *testing.T) {
	info := StorageInfo{
		ID:          "1",
		Name:        "test",
		StorageType: StorageTypeWebDAV,
	}
	assert.Equal(t, "1", info.ID)
	assert.Empty(t, info.BaseURL)
	assert.False(t, info.IsConnected)
}

func TestQuotaInfo(t *testing.T) {
	quota := QuotaInfo{
		UsedBytes:      100,
		TotalBytes:     1000,
		AvailableBytes: 900,
		UsedPercentage: 10.0,
	}
	assert.Equal(t, int64(100), quota.UsedBytes)
	assert.Equal(t, 10.0, quota.UsedPercentage)
}

func TestFileInfo(t *testing.T) {
	file := FileInfo{
		Name: "test.txt",
		Path: "/docs/test.txt",
		Size: 1024,
	}
	assert.Equal(t, "test.txt", file.Name)
	assert.False(t, file.IsDirectory)
	assert.Empty(t, file.LastModified)
}

func TestFileInfoDirectory(t *testing.T) {
	dir := FileInfo{
		Name:        "docs",
		Path:        "/docs",
		Size:        0,
		IsDirectory: true,
	}
	assert.True(t, dir.IsDirectory)
}

func TestCommonConfigDefaults(t *testing.T) {
	cc := NewCommonConfig("test", StorageTypeWebDAV)
	assert.Equal(t, "test", cc.Name)
	assert.Equal(t, StorageTypeWebDAV, cc.Type)
	assert.True(t, cc.IsEnabled)
	assert.Equal(t, 100, cc.Priority)
	assert.Nil(t, cc.Metadata)
}

func TestCommonConfigWithMetadata(t *testing.T) {
	cc := NewCommonConfig("test", StorageTypeFTP)
	cc.Metadata = map[string]string{"key": "value"}
	assert.Equal(t, "value", cc.Metadata["key"])
}

func TestWebDavAuthTypes(t *testing.T) {
	types := []WebDavAuthType{WebDavAuthBasic, WebDavAuthDigest, WebDavAuthOAuth, WebDavAuthNone}
	assert.Len(t, types, 4)
}

func TestOneDriveDriveTypes(t *testing.T) {
	types := []OneDriveDriveType{OneDriveDriveMe, OneDriveDriveBusiness, OneDriveDriveSharePoint, OneDriveDriveGroup}
	assert.Len(t, types, 4)
}

func TestConfigDisable(t *testing.T) {
	cfg := NewWebDavConfig("Test", "https://test.com", "u", "p")
	cfg.IsEnabled = false
	assert.False(t, cfg.IsEnabled)
}

func TestConfigPriority(t *testing.T) {
	cfg := NewFtpConfig("Test", "ftp.test.com", "u", "p")
	cfg.Priority = 50
	assert.Equal(t, 50, cfg.Priority)
}
