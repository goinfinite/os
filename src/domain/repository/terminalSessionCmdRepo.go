package repository

import (
	"errors"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

var (
	ErrTerminalSessionNotFound           = errors.New("TerminalSessionNotFound")
	ErrTerminalSessionAccountCapReached  = errors.New("TerminalSessionAccountCapReached")
	ErrTerminalSessionAccountRequired    = errors.New("TerminalSessionAccountRequired")
	ErrTerminalSessionWorkingDirNotFound = errors.New("TerminalSessionWorkingDirNotFound")
)

type TerminalSessionAttachHandle interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Resize(cols, rows uint16) error
	Wait() error
	Close() error
}

type TerminalSessionCmdRepo interface {
	Create(createDto dto.CreateTerminalSession) (valueObject.TerminalSessionId, error)
	Delete(deleteDto dto.DeleteTerminalSession) error
	Attach(attachDto dto.AttachTerminalSession) (TerminalSessionAttachHandle, error)
}
