package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdateTerminalSession struct {
	Id                valueObject.TerminalSessionId    `json:"id"`
	Name              *valueObject.TerminalSessionName `json:"name"`
	AccountUsername   valueObject.Username             `json:"-"`
	RunAsUsername     valueObject.Username             `json:"-"`
	OperatorAccountId tkValueObject.AccountId          `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress          `json:"-"`
}

func NewUpdateTerminalSession(
	id valueObject.TerminalSessionId,
	name *valueObject.TerminalSessionName,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdateTerminalSession {
	return UpdateTerminalSession{
		Id:                id,
		Name:              name,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
