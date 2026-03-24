package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateUnixFiles struct {
	SourcePaths          []tkValueObject.UnixAbsoluteFilePath `json:"sourcePaths"`
	DestinationPath      *tkValueObject.UnixAbsoluteFilePath  `json:"destinationPath"`
	Permissions          *valueObject.UnixFilePermissions     `json:"permissions"`
	EncodedContent       *valueObject.EncodedContent          `json:"encodedContent"`
	Ownership            *tkValueObject.UnixFileOwnership     `json:"ownership"`
	ShouldFixPermissions *bool                                `json:"shouldFixPermissions"`
	OperatorAccountId    tkValueObject.AccountId              `json:"-"`
	OperatorIpAddress    tkValueObject.IpAddress              `json:"-"`
}

func NewUpdateUnixFiles(
	sourcePaths []tkValueObject.UnixAbsoluteFilePath,
	destinationPath *tkValueObject.UnixAbsoluteFilePath,
	permissions *valueObject.UnixFilePermissions,
	encodedContent *valueObject.EncodedContent,
	ownership *tkValueObject.UnixFileOwnership,
	shouldFixPermissions *bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdateUnixFiles {
	return UpdateUnixFiles{
		SourcePaths:          sourcePaths,
		DestinationPath:      destinationPath,
		Permissions:          permissions,
		EncodedContent:       encodedContent,
		Ownership:            ownership,
		ShouldFixPermissions: shouldFixPermissions,
		OperatorAccountId:    operatorAccountId,
		OperatorIpAddress:    operatorIpAddress,
	}
}
