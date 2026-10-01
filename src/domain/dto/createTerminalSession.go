package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateTerminalSession struct {
	AccountUsername   valueObject.Username                `json:"-"`
	Name              *valueObject.TerminalSessionName    `json:"name"`
	WorkingDir        *tkValueObject.UnixAbsoluteFilePath `json:"workingDir"`
	Command           *tkValueObject.UnixCommand          `json:"command"`
	AccountId         *tkValueObject.AccountId            `json:"accountId"`
	OperatorAccountId tkValueObject.AccountId             `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress             `json:"-"`
}

func NewCreateTerminalSession(
	name *valueObject.TerminalSessionName,
	workingDir *tkValueObject.UnixAbsoluteFilePath,
	command *tkValueObject.UnixCommand,
	accountId *tkValueObject.AccountId,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateTerminalSession {
	return CreateTerminalSession{
		Name:              name,
		WorkingDir:        workingDir,
		Command:           command,
		AccountId:         accountId,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
