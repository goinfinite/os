package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CopyUnixFile struct {
	SourcePath        valueObject.UnixFilePath `json:"sourcePath"`
	DestinationPath   valueObject.UnixFilePath `json:"destinationPath"`
	ShouldOverwrite   bool                     `json:"shouldOverwrite"`
	OperatorAccountId valueObject.AccountId    `json:"-"`
	OperatorIpAddress valueObject.IpAddress    `json:"-"`
}

func NewCopyUnixFile(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
	shouldOverwrite bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CopyUnixFile {
	return CopyUnixFile{
		SourcePath:        sourcePath,
		DestinationPath:   destinationPath,
		ShouldOverwrite:   shouldOverwrite,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
