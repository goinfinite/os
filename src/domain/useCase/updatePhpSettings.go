package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

type UpdatePhpSettings struct {
	runtimeCmdRepo        repository.RuntimeCmdRepo
	vhostQueryRepo        repository.VirtualHostQueryRepo
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo
}

func NewUpdatePhpSettings(
	runtimeCmdRepo repository.RuntimeCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
) UpdatePhpSettings {
	return UpdatePhpSettings{
		runtimeCmdRepo:        runtimeCmdRepo,
		vhostQueryRepo:        vhostQueryRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
	}
}

func (uc UpdatePhpSettings) Execute(request dto.UpdatePhpSettingsRequest) error {
	_, err := uc.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &request.Hostname,
	})
	if err != nil {
		slog.Error("VirtualHostNotFound", slog.String("err", err.Error()))
		return errors.New("VirtualHostNotFound")
	}

	err = uc.runtimeCmdRepo.UpdatePhpSettings(
		request.Hostname, request.PhpSettings,
	)
	if err != nil {
		slog.Error("UpdatePhpSettingsError", slog.String("err", err.Error()))
		return errors.New("UpdatePhpSettingsInfraError")
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		UpdatePhpSettings(request)

	return nil
}
