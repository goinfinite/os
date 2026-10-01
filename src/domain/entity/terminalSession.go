package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type TerminalSession struct {
	Id              valueObject.TerminalSessionId      `json:"id"`
	Name            *valueObject.TerminalSessionName   `json:"name,omitempty"`
	AccountId       tkValueObject.AccountId            `json:"-"`
	AccountUsername valueObject.Username               `json:"accountUsername"`
	WorkingDir      tkValueObject.UnixAbsoluteFilePath `json:"workingDir"`
	Command         tkValueObject.UnixCommand          `json:"command"`
	CreatedAt       tkValueObject.UnixTime             `json:"createdAt"`
	AttachedClients uint16                             `json:"attachedClients"`
}

func NewTerminalSession(
	id valueObject.TerminalSessionId,
	name *valueObject.TerminalSessionName,
	accountId tkValueObject.AccountId,
	accountUsername valueObject.Username,
	workingDir tkValueObject.UnixAbsoluteFilePath,
	command tkValueObject.UnixCommand,
	createdAt tkValueObject.UnixTime,
	attachedClients uint16,
) TerminalSession {
	return TerminalSession{
		Id:              id,
		Name:            name,
		AccountId:       accountId,
		AccountUsername: accountUsername,
		WorkingDir:      workingDir,
		Command:         command,
		CreatedAt:       createdAt,
		AttachedClients: attachedClients,
	}
}

func (entity TerminalSession) ReadName() string {
	if entity.Name != nil {
		return entity.Name.String()
	}

	return entity.AccountUsername.String() + "@" + entity.Id.String()[:4]
}
