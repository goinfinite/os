package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CompressUnixFiles struct {
	SourcePaths       []tkValueObject.UnixAbsoluteFilePath `json:"sourcePaths"`
	DestinationPath   tkValueObject.UnixAbsoluteFilePath   `json:"destinationPath"`
	CompressionType   *valueObject.UnixCompressionType     `json:"compressionType"`
	OperatorAccountId tkValueObject.AccountId              `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress              `json:"-"`
}

func NewCompressUnixFiles(
	sourcePaths []tkValueObject.UnixAbsoluteFilePath,
	destinationPath tkValueObject.UnixAbsoluteFilePath,
	compressionType *valueObject.UnixCompressionType,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CompressUnixFiles {
	return CompressUnixFiles{
		SourcePaths:       sourcePaths,
		DestinationPath:   destinationPath,
		CompressionType:   compressionType,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
