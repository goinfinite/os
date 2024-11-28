package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type UploadUnixFiles struct {
	DestinationPath    valueObject.UnixFilePath        `json:"destinationPath"`
	FileStreamHandlers []valueObject.FileStreamHandler `json:"fileStreamHandlers"`
	OperatorAccountId  valueObject.AccountId           `json:"-"`
	OperatorIpAddress  valueObject.IpAddress           `json:"-"`
}

func NewUploadUnixFiles(
	destinationPath valueObject.UnixFilePath,
	fileStreamHandlers []valueObject.FileStreamHandler,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) UploadUnixFiles {
	return UploadUnixFiles{
		DestinationPath:    destinationPath,
		FileStreamHandlers: fileStreamHandlers,
		OperatorAccountId:  operatorAccountId,
		OperatorIpAddress:  operatorIpAddress,
	}
}
