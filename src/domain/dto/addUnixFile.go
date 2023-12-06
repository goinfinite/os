package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddUnixFile struct {
	Path        valueObject.UnixFilePath        `json:"path"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	Type        valueObject.UnixFileType        `json:"type"`
}

func NewAddUnixFile(
	Path valueObject.UnixFilePath,
	Permissions valueObject.UnixFilePermissions,
	Type valueObject.UnixFileType,
) AddUnixFile {
	return AddUnixFile{
		Path:        Path,
		Permissions: Permissions,
		Type:        Type,
	}
}
