package entity

import "github.com/speedianet/os/src/domain/valueObject"

type UnixFile struct {
	Uid         valueObject.UnixUid             `json:"uid"`
	Gid         valueObject.GroupId             `json:"gid"`
	MimeType    valueObject.MimeType            `json:"mimeType"`
	Name        valueObject.UnixFileName        `json:"name"`
	Path        valueObject.UnixFilePath        `json:"path"`
	Extension   valueObject.UnixFileExtension   `json:"extension"`
	Permissions valueObject.UnixFilePermissions `json:"permissions"`
	Size        valueObject.Byte                `json:"size"`
	UpdatedAt   valueObject.UnixTime            `json:"updatedAt"`
	Owner       valueObject.Username            `json:"owner"`
	Group       valueObject.GroupName           `json:"group"`
}

func NewUnixFile(
	Uid valueObject.UnixUid,
	Gid valueObject.GroupId,
	MimeType valueObject.MimeType,
	Name valueObject.UnixFileName,
	Path valueObject.UnixFilePath,
	Extension valueObject.UnixFileExtension,
	Permissions valueObject.UnixFilePermissions,
	Size valueObject.Byte,
	UpdatedAt valueObject.UnixTime,
	Owner valueObject.Username,
	Group valueObject.GroupName,
) UnixFile {
	return UnixFile{
		Uid:         Uid,
		Gid:         Gid,
		MimeType:    MimeType,
		Name:        Name,
		Path:        Path,
		Extension:   Extension,
		Permissions: Permissions,
		Size:        Size,
		UpdatedAt:   UpdatedAt,
		Owner:       Owner,
		Group:       Group,
	}
}
