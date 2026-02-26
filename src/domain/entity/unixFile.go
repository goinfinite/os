package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UnixFile struct {
	Name        tkValueObject.UnixFileName         `json:"name"`
	Path        tkValueObject.UnixAbsoluteFilePath `json:"path"`
	MimeType    tkValueObject.MimeType             `json:"mimeType"`
	Permissions valueObject.UnixFilePermissions    `json:"permissions"`
	Size        tkValueObject.Byte                 `json:"size"`
	Extension   *tkValueObject.UnixFileExtension   `json:"extension"`
	Content     *valueObject.UnixFileContent       `json:"content"`
	Uid         tkValueObject.UnixUserId           `json:"uid"`
	Owner       valueObject.Username               `json:"owner"`
	Gid         tkValueObject.UnixGroupId          `json:"gid"`
	Group       tkValueObject.UnixGroupName        `json:"group"`
	UpdatedAt   tkValueObject.UnixTime             `json:"updatedAt"`
}

func NewUnixFile(
	name tkValueObject.UnixFileName,
	path tkValueObject.UnixAbsoluteFilePath,
	mimeType tkValueObject.MimeType,
	permissions valueObject.UnixFilePermissions,
	size tkValueObject.Byte,
	extension *tkValueObject.UnixFileExtension,
	content *valueObject.UnixFileContent,
	uid tkValueObject.UnixUserId,
	owner valueObject.Username,
	gid tkValueObject.UnixGroupId,
	group tkValueObject.UnixGroupName,
	updatedAt tkValueObject.UnixTime,
) UnixFile {
	return UnixFile{
		Name:        name,
		Path:        path,
		MimeType:    mimeType,
		Permissions: permissions,
		Size:        size,
		Extension:   extension,
		Content:     content,
		Uid:         uid,
		Owner:       owner,
		Gid:         gid,
		Group:       group,
		UpdatedAt:   updatedAt,
	}
}

func (entity UnixFile) ToSimplified() SimplifiedUnixFile {
	return SimplifiedUnixFile{
		Name:     entity.Name,
		Path:     entity.Path,
		MimeType: entity.MimeType,
	}
}

type SimplifiedUnixFile struct {
	Name     tkValueObject.UnixFileName         `json:"name"`
	Path     tkValueObject.UnixAbsoluteFilePath `json:"path"`
	MimeType tkValueObject.MimeType             `json:"mimeType"`
}

func NewSimplifiedUnixFile(
	name tkValueObject.UnixFileName,
	path tkValueObject.UnixAbsoluteFilePath,
	mimeType tkValueObject.MimeType,
) SimplifiedUnixFile {
	return SimplifiedUnixFile{
		Name:     name,
		Path:     path,
		MimeType: mimeType,
	}
}
