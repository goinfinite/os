package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteTerminalSession struct {
	Id                valueObject.TerminalSessionId `json:"id"`
	AccountUsername   valueObject.Username          `json:"-"`
	OperatorAccountId tkValueObject.AccountId       `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress       `json:"-"`
}

func NewDeleteTerminalSession(
	id valueObject.TerminalSessionId,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteTerminalSession {
	return DeleteTerminalSession{
		Id:                id,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
