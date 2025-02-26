package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type UpdateUnixFiles struct {
	SourcePaths       []valueObject.UnixFilePath       `json:"sourcePaths"`
	DestinationPath   *valueObject.UnixFilePath        `json:"destinationPath"`
	Permissions       *valueObject.UnixFilePermissions `json:"permissions"`
	EncodedContent    *valueObject.EncodedContent      `json:"encodedContent"`
	Ownership         *valueObject.UnixFileOwnership   `json:"ownership"`
	FixPermissions    *bool                            `json:"fixPermissions"`
	OperatorAccountId valueObject.AccountId            `json:"-"`
	OperatorIpAddress valueObject.IpAddress            `json:"-"`
}

func NewUpdateUnixFiles(
	sourcePaths []valueObject.UnixFilePath,
	destinationPath *valueObject.UnixFilePath,
	permissions *valueObject.UnixFilePermissions,
	encodedContent *valueObject.EncodedContent,
	ownership *valueObject.UnixFileOwnership,
	fixPermissions *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UpdateUnixFiles {
	return UpdateUnixFiles{
		SourcePaths:       sourcePaths,
		DestinationPath:   destinationPath,
		Permissions:       permissions,
		EncodedContent:    encodedContent,
		Ownership:         ownership,
		FixPermissions:    fixPermissions,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
