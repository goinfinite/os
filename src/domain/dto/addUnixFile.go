package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddUnixFile struct {
	Path        valueObject.UnixFilePath        `json:"path"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	Type        valueObject.UnixFileType        `json:"type"`
}

func NewAddUnixFile(
	path valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
	fileType valueObject.UnixFileType,
) AddUnixFile {
	return AddUnixFile{
		Path:        path,
		Permissions: permissions,
		Type:        fileType,
	}
}
