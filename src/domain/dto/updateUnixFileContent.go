package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFileContent struct {
	Path    valueObject.UnixFilePath    `json:"path"`
	Content valueObject.UnixFileContent `json:"content"`
}

func NewUpdateUnixFileContent(
	path valueObject.UnixFilePath,
	content valueObject.UnixFileContent,
) UpdateUnixFileContent {
	return UpdateUnixFileContent{
		Path:    path,
		Content: content,
	}
}
