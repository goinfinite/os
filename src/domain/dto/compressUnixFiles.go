package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CompressUnixFiles struct {
	SourcePaths       []valueObject.UnixFilePath       `json:"sourcePaths"`
	DestinationPath   valueObject.UnixFilePath         `json:"destinationPath"`
	CompressionType   *valueObject.UnixCompressionType `json:"compressionType"`
	OperatorAccountId valueObject.AccountId            `json:"-"`
	OperatorIpAddress valueObject.IpAddress            `json:"-"`
}

func NewCompressUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	compressionType *valueObject.UnixCompressionType,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CompressUnixFiles {
	return CompressUnixFiles{
		SourcePaths:       sourcePaths,
		DestinationPath:   destinationPath,
		CompressionType:   compressionType,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
