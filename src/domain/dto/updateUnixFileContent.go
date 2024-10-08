package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateUnixFileContent struct {
	SourcePath valueObject.UnixFilePath   `json:"sourcePath"`
	Content    valueObject.EncodedContent `json:"content"`
}

func NewUpdateUnixFileContent(
	sourcePath valueObject.UnixFilePath,
	content valueObject.EncodedContent,
) UpdateUnixFileContent {
	return UpdateUnixFileContent{
		SourcePath: sourcePath,
		Content:    content,
	}
}
