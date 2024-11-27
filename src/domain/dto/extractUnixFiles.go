package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type ExtractUnixFiles struct {
	SourcePath        valueObject.UnixFilePath `json:"sourcePath"`
	DestinationPath   valueObject.UnixFilePath `json:"destinationPath"`
	OperatorAccountId valueObject.AccountId    `json:"-"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewExtractUnixFiles(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) ExtractUnixFiles {
	return ExtractUnixFiles{
		SourcePath:        sourcePath,
		DestinationPath:   destinationPath,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
