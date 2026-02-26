package dto

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type DeleteUnixFiles struct {
	SourcePaths       []tkValueObject.UnixAbsoluteFilePath `json:"sourcePaths"`
	HardDelete        bool                                 `json:"hardDelete"`
	OperatorAccountId tkValueObject.AccountId              `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress              `json:"-"`
}

func NewDeleteUnixFiles(
	sourcePaths []tkValueObject.UnixAbsoluteFilePath,
	hardDelete bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteUnixFiles {
	return DeleteUnixFiles{
		SourcePaths:       sourcePaths,
		HardDelete:        hardDelete,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
