package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteUnixFiles struct {
	SourcePaths       []valueObject.UnixFilePath `json:"sourcePaths"`
	HardDelete        bool                       `json:"hardDelete"`
	OperatorAccountId valueObject.AccountId      `json:"-"`
	OperatorIpAddress valueObject.IpAddress      `json:"-"`
}

func NewDeleteUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	hardDelete bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteUnixFiles {
	return DeleteUnixFiles{
		SourcePaths:       sourcePaths,
		HardDelete:        hardDelete,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
