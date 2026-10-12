package terminalSessionInfra

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type TerminalSessionCmdRepo struct{}

func NewTerminalSessionCmdRepo() *TerminalSessionCmdRepo {
	return &TerminalSessionCmdRepo{}
}

func (repo *TerminalSessionCmdRepo) sessionIdFactory() (
	terminalSessionId valueObject.TerminalSessionId, err error,
) {
	randomBytes := make([]byte, 8)
	_, err = rand.Read(randomBytes)
	if err != nil {
		return terminalSessionId, errors.New(
			"GenerateTerminalSessionIdError: " + err.Error(),
		)
	}

	return valueObject.NewTerminalSessionId(hex.EncodeToString(randomBytes))
}

func (repo *TerminalSessionCmdRepo) killSessionAfterCreateFailure(
	client *TerminalMultiplexerClient,
	terminalSessionId valueObject.TerminalSessionId,
) {
	killErr := client.KillSession(terminalSessionId)
	if killErr != nil {
		slog.Error(
			"KillTerminalSessionAfterCreateError",
			slog.String("err", killErr.Error()),
		)
	}
}

func (repo *TerminalSessionCmdRepo) Create(
	createDto dto.CreateTerminalSession,
) (terminalSessionId valueObject.TerminalSessionId, err error) {
	workingDirStr := createDto.WorkingDir.String()
	workingDirInfo, err := os.Stat(workingDirStr)
	if err != nil || !workingDirInfo.IsDir() {
		return terminalSessionId, repository.ErrTerminalSessionWorkingDirNotFound
	}

	terminalSessionId, err = repo.sessionIdFactory()
	if err != nil {
		return terminalSessionId, err
	}

	client, err := NewTerminalMultiplexerClient(
		createDto.AccountUsername, createDto.RunAsUsername,
	)
	if err != nil {
		return terminalSessionId, err
	}

	err = client.CreateSession(terminalSessionId, *createDto.WorkingDir)
	if err != nil {
		return terminalSessionId, err
	}

	if createDto.Name != nil {
		err = client.SetSessionName(terminalSessionId, *createDto.Name)
		if err != nil {
			repo.killSessionAfterCreateFailure(client, terminalSessionId)
			return terminalSessionId, err
		}
	}

	if createDto.Command != nil {
		err = client.SendCommand(terminalSessionId, *createDto.Command)
		if err != nil {
			repo.killSessionAfterCreateFailure(client, terminalSessionId)
			return terminalSessionId, err
		}
	}

	return terminalSessionId, nil
}

func (repo *TerminalSessionCmdRepo) Update(
	updateDto dto.UpdateTerminalSession,
) error {
	client, err := NewTerminalMultiplexerClient(
		updateDto.AccountUsername, updateDto.RunAsUsername,
	)
	if err != nil {
		return err
	}

	if updateDto.Name == nil {
		return client.ClearSessionName(updateDto.Id)
	}

	return client.SetSessionName(updateDto.Id, *updateDto.Name)
}

func (repo *TerminalSessionCmdRepo) Delete(
	deleteDto dto.DeleteTerminalSession,
) error {
	client, err := NewTerminalMultiplexerClient(
		deleteDto.AccountUsername, deleteDto.RunAsUsername,
	)
	if err != nil {
		return err
	}

	return client.KillSession(deleteDto.Id)
}

func (repo *TerminalSessionCmdRepo) Attach(
	attachDto dto.AttachTerminalSession,
) (repository.TerminalSessionAttachHandle, error) {
	client, err := NewTerminalMultiplexerClient(
		attachDto.AccountUsername, attachDto.RunAsUsername,
	)
	if err != nil {
		return nil, err
	}

	return client.Attach(attachDto.Id)
}
