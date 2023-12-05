package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddUnixFile struct {
	MimeType    valueObject.MimeType            `json:"mimeType"`
	Name        valueObject.UnixFileName        `json:"name"`
	Path        valueObject.UnixFilePath        `json:"path"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
}

func NewAddUnixFile(
	MimeType valueObject.MimeType,
	Name valueObject.UnixFileName,
	Path valueObject.UnixFilePath,
	Permissions valueObject.UnixFilePermissions,
) AddUnixFile {
	return AddUnixFile{
		MimeType:    MimeType,
		Name:        Name,
		Path:        Path,
		Permissions: Permissions,
	}
}
