package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFileContent struct {
	Path    valueObject.UnixFilePath    `json:"path"`
	Content valueObject.UnixFileContent `json:"content"`
}

func NewUpdateUnixFileContent(
	Path valueObject.UnixFilePath,
	Content valueObject.UnixFileContent,
) UpdateUnixFileContent {
	return UpdateUnixFileContent{
		Path:    Path,
		Content: Content,
	}
}
