package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CreateUnixFile struct {
	SourcePath  valueObject.UnixFilePath        `json:"sourcePath"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	MimeType    valueObject.MimeType            `json:"mimeType"`
}

func NewCreateUnixFile(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
	mimeType valueObject.MimeType,
) CreateUnixFile {
	return CreateUnixFile{
		SourcePath:  sourcePath,
		Permissions: permissions,
		MimeType:    mimeType,
	}
}
