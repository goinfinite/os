package dto

import "github.com/speedianet/os/src/domain/valueObject"

type AddUnixFileCopy struct {
	OriginPath      valueObject.UnixFilePath `json:"path"`
	DestinationPath valueObject.UnixFilePath `json:"destinationPath"`
}

func NewAddUnixFileCopy(
	originPath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) AddUnixFileCopy {
	return AddUnixFileCopy{
		OriginPath:      originPath,
		DestinationPath: destinationPath,
	}
}
