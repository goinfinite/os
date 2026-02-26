package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type ExtractUnixFiles struct {
	SourcePath        tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	DestinationPath   tkValueObject.UnixAbsoluteFilePath `json:"destinationPath"`
	OperatorAccountId tkValueObject.AccountId            `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress            `json:"-"`
}

func NewExtractUnixFiles(
	sourcePath tkValueObject.UnixAbsoluteFilePath,
	destinationPath tkValueObject.UnixAbsoluteFilePath,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) ExtractUnixFiles {
	return ExtractUnixFiles{
		SourcePath:        sourcePath,
		DestinationPath:   destinationPath,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
