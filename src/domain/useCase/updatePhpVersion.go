package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	useCaseHelper "github.com/goinfinite/os/src/domain/useCase/helper"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

type UpdatePhpVersion struct {
	runtimeQueryRepo      repository.RuntimeQueryRepo
	runtimeCmdRepo        repository.RuntimeCmdRepo
	vhostQueryRepo        repository.VirtualHostQueryRepo
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo
}

func NewUpdatePhpVersion(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo,
) UpdatePhpVersion {
	return UpdatePhpVersion{
		runtimeQueryRepo:      runtimeQueryRepo,
		runtimeCmdRepo:        runtimeCmdRepo,
		vhostQueryRepo:        vhostQueryRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
	}
}

func (uc UpdatePhpVersion) Execute(request dto.UpdatePhpVersionRequest) error {
	if !useCaseHelper.IsPhpVersionInstalled(
		uc.runtimeQueryRepo, request.PhpVersion,
	) {
		return errors.New("PhpVersionNotInstalled")
	}

	_, err := uc.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &request.Hostname,
	})
	if err != nil {
		slog.Error("VirtualHostNotFound", slog.String("err", err.Error()))
		return errors.New("VirtualHostNotFound")
	}

	err = uc.runtimeCmdRepo.UpdatePhpVersion(
		request.Hostname, request.PhpVersion,
	)
	if err != nil {
		slog.Error("UpdatePhpVersionError", slog.String("err", err.Error()))
		return errors.New("UpdatePhpVersionInfraError")
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		UpdatePhpVersion(request)

	return nil
}
