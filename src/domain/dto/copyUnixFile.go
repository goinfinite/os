package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type CopyUnixFile struct {
	SourcePath        tkValueObject.UnixAbsoluteFilePath `json:"sourcePath"`
	DestinationPath   tkValueObject.UnixAbsoluteFilePath `json:"destinationPath"`
	ShouldOverwrite   bool                               `json:"shouldOverwrite"`
	OperatorAccountId tkValueObject.AccountId            `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress            `json:"-"`
}

func NewCopyUnixFile(
	sourcePath tkValueObject.UnixAbsoluteFilePath,
	destinationPath tkValueObject.UnixAbsoluteFilePath,
	shouldOverwrite bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CopyUnixFile {
	return CopyUnixFile{
		SourcePath:        sourcePath,
		DestinationPath:   destinationPath,
		ShouldOverwrite:   shouldOverwrite,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
