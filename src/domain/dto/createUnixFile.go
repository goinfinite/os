package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateUnixFile struct {
	FilePath          tkValueObject.UnixAbsoluteFilePath `json:"filePath"`
	Permissions       valueObject.UnixFilePermissions    `json:"permissions"`
	MimeType          tkValueObject.MimeType             `json:"mimeType"`
	OperatorAccountId tkValueObject.AccountId            `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress            `json:"-"`
}

func NewCreateUnixFile(
	filePath tkValueObject.UnixAbsoluteFilePath,
	permissionsPtr *valueObject.UnixFilePermissions,
	mimeType tkValueObject.MimeType,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateUnixFile {
	permissions := valueObject.NewUnixFileDefaultPermissions()
	if mimeType.IsDir() {
		permissions = valueObject.NewUnixDirDefaultPermissions()
	}

	if permissionsPtr != nil {
		permissions = *permissionsPtr
	}

	return CreateUnixFile{
		FilePath:          filePath,
		Permissions:       permissions,
		MimeType:          mimeType,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
