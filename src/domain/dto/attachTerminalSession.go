package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type AttachTerminalSession struct {
	Id                valueObject.TerminalSessionId `json:"id"`
	AccountUsername   valueObject.Username          `json:"-"`
	OperatorAccountId tkValueObject.AccountId       `json:"-"`
}

func NewAttachTerminalSession(
	id valueObject.TerminalSessionId,
	operatorAccountId tkValueObject.AccountId,
) AttachTerminalSession {
	return AttachTerminalSession{
		Id:                id,
		OperatorAccountId: operatorAccountId,
	}
}
