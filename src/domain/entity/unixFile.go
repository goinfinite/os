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
	Stream      *os.File                        `json:"-"`
}

func NewUnixFile(
	uid valueObject.UnixUid,
	owner valueObject.Username,
	gid valueObject.GroupId,
	group valueObject.GroupName,
	mimeType valueObject.MimeType,
	name valueObject.UnixFileName,
	path valueObject.UnixFilePath,
	extension *valueObject.UnixFileExtension,
	permissions valueObject.UnixFilePermissions,
	size valueObject.Byte,
	updatedAt valueObject.UnixTime,
	stream *os.File,
) UnixFile {
	return UnixFile{
		Uid:         uid,
		Owner:       owner,
		Gid:         gid,
		Group:       group,
		MimeType:    mimeType,
		Name:        name,
		Path:        path,
		Extension:   extension,
		Permissions: permissions,
		Size:        size,
		UpdatedAt:   updatedAt,
		Stream:      stream,
	}
}
