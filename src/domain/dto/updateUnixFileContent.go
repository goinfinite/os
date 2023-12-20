package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFileContent struct {
	Path    valueObject.UnixFilePath   `json:"path"`
	Content valueObject.EncodedContent `json:"content"`
}

func NewUpdateUnixFileContent(
	path valueObject.UnixFilePath,
	content valueObject.EncodedContent,
) UpdateUnixFileContent {
	return UpdateUnixFileContent{
		Path:    path,
		Content: content,
	}
}
