package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type ReadTerminalSessionsRequest struct {
	Pagination        tkDto.Pagination               `json:"pagination"`
	TerminalSessionId *valueObject.TerminalSessionId `json:"id,omitempty"`
	AccountId         *tkValueObject.AccountId       `json:"accountId,omitempty"`
	OperatorAccountId tkValueObject.AccountId        `json:"-"`
}

type ReadTerminalSessionsResponse struct {
	Pagination       tkDto.Pagination         `json:"pagination"`
	TerminalSessions []entity.TerminalSession `json:"terminalSessions"`
}
