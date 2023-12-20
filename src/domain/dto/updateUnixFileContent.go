package dto

import "github.com/speedianet/os/src/domain/valueObject"

type UpdateUnixFileContent struct {
	Path    valueObject.UnixFilePath         `json:"path"`
	Content valueObject.EncodedBase64Content `json:"content"`
}

func NewUpdateUnixFileContent(
	path valueObject.UnixFilePath,
	content valueObject.EncodedBase64Content,
) UpdateUnixFileContent {
	return UpdateUnixFileContent{
		Path:    path,
		Content: content,
	}
}
