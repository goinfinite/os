package dto

import "github.com/speedianet/os/src/domain/valueObject"

type CopyUnixFile struct {
	OriginPath      valueObject.UnixFilePath `json:"path"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewCopyUnixFile(
	originPath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) CopyUnixFile {
	return CopyUnixFile{
		OriginPath:      originPath,
		DestinationPath: destinationPath,
	}
}
