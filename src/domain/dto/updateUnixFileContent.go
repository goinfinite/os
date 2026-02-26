package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateUnixFileContent struct {
	SourcePath tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	Content    valueObject.EncodedContent         `json:"content"`
}

func NewUpdateUnixFileContent(
	sourcePath tkValueObject.UnixAbsoluteFilePath,
	content valueObject.EncodedContent,
) UpdateUnixFileContent {
	return UpdateUnixFileContent{
		SourcePath: sourcePath,
		Content:    content,
	}
}
