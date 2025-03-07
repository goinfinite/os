package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateUnixFile struct {
	FilePath          valueObject.UnixFilePath        `json:"filePath"`
	Permissions       valueObject.UnixFilePermissions `json:"permissions"`
	MimeType          valueObject.MimeType            `json:"mimeType"`
	OperatorAccountId valueObject.AccountId           `json:"-"`
	OperatorIpAddress valueObject.IpAddress           `json:"-"`
}

func NewCreateUnixFile(
	filePath valueObject.UnixFilePath,
	permissionsPtr *valueObject.UnixFilePermissions,
	mimeType valueObject.MimeType,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
