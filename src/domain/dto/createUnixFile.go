package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateUnixFile struct {
	FilePath    valueObject.UnixFilePath        `json:"filePath"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	MimeType    valueObject.MimeType            `json:"mimeType"`
}

func NewCreateUnixFile(
	filePath valueObject.UnixFilePath,
	permissionsPtr *valueObject.UnixFilePermissions,
	mimeType valueObject.MimeType,
) CreateUnixFile {
	permissions, _ := valueObject.NewUnixFilePermissions("644")
	if mimeType.IsDir() {
		permissions, _ = valueObject.NewUnixFilePermissions("755")
	}

	if permissionsPtr != nil {
		permissions = *permissionsPtr
	}

	return CreateUnixFile{
		FilePath:    filePath,
		Permissions: permissions,
		MimeType:    mimeType,
	}
}
