package entity

import (
	"os"

	"github.com/speedianet/os/src/domain/valueObject"
)

type UnixFile struct {
	Uid         valueObject.UnixUid             `json:"uid"`
	Owner       valueObject.Username            `json:"owner"`
	Gid         valueObject.GroupId             `json:"gid"`
	Group       valueObject.GroupName           `json:"group"`
	MimeType    valueObject.MimeType            `json:"mimeType"`
	Name        valueObject.UnixFileName        `json:"name"`
	Path        valueObject.UnixFilePath        `json:"path"`
	Extension   *valueObject.UnixFileExtension  `json:"extension"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	Size        valueObject.Byte                `json:"size"`
	UpdatedAt   valueObject.UnixTime            `json:"updatedAt"`
}

func NewUnixFile(
	Uid valueObject.UnixUid,
	Owner valueObject.Username,
	Gid valueObject.GroupId,
	Group valueObject.GroupName,
	MimeType valueObject.MimeType,
	Name valueObject.UnixFileName,
	Path valueObject.UnixFilePath,
	Extension *valueObject.UnixFileExtension,
	Permissions valueObject.UnixFilePermissions,
	Size valueObject.Byte,
	UpdatedAt valueObject.UnixTime,
) UnixFile {
	return UnixFile{
		Uid:         Uid,
		Owner:       Owner,
		Gid:         Gid,
		Group:       Group,
		MimeType:    MimeType,
		Name:        Name,
		Path:        Path,
		Extension:   Extension,
		Permissions: Permissions,
		Size:        Size,
		UpdatedAt:   UpdatedAt,
	}
}

func (unixFile UnixFile) Open() (*os.File, error) {
	return os.Open(unixFile.Path.String())
}
