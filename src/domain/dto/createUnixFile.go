package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateUnixFile struct {
	FilePath    valueObject.UnixFilePath        `json:"filePath"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	MimeType    valueObject.MimeType            `json:"mimeType"`
}

func NewCreateUnixFile(
	filePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
	mimeType valueObject.MimeType,
) CreateUnixFile {
	return CreateUnixFile{
		FilePath:    filePath,
		Permissions: permissions,
		MimeType:    mimeType,
	}
}
