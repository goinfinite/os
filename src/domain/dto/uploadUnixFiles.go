package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UploadUnixFiles struct {
	DestinationPath    tkValueObject.UnixAbsoluteFilePath `json:"destinationPath"`
	FileStreamHandlers []valueObject.FileStreamHandler    `json:"fileStreamHandlers"`
	OperatorAccountId  tkValueObject.AccountId            `json:"-"`
	OperatorIpAddress  tkValueObject.IpAddress            `json:"-"`
}

func NewUploadUnixFiles(
	destinationPath tkValueObject.UnixAbsoluteFilePath,
	fileStreamHandlers []valueObject.FileStreamHandler,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UploadUnixFiles {
	return UploadUnixFiles{
		DestinationPath:    destinationPath,
		FileStreamHandlers: fileStreamHandlers,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
