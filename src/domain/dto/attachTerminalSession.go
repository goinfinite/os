package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
)

type AttachTerminalSession struct {
	Id              valueObject.TerminalSessionId `json:"id"`
	AccountUsername valueObject.Username          `json:"-"`
}

func NewAttachTerminalSession(
	id valueObject.TerminalSessionId,
	accountUsername valueObject.Username,
) AttachTerminalSession {
	return AttachTerminalSession{
		Id:              id,
		AccountUsername: accountUsername,
	}
}
