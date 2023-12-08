package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddUnixFileCopy struct {
	OriginPath      valueObject.UnixFilePath `json:"path"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewAddUnixFileCopy(
	OriginPath valueObject.UnixFilePath,
	DestinationPath valueObject.UnixFilePath,
) AddUnixFileCopy {
	return AddUnixFileCopy{
		OriginPath:      OriginPath,
		DestinationPath: DestinationPath,
	}
}
