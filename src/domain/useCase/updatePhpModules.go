package useCase

import (
	"errors"
	"log/slog"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	useCaseHelper "github.com/goinfinite/os/src/domain/useCase/helper"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
)

type UpdatePhpModules struct {
	runtimeQueryRepo      repository.RuntimeQueryRepo
	runtimeCmdRepo        repository.RuntimeCmdRepo
	vhostQueryRepo        repository.VirtualHostQueryRepo
	activityRecordCmdRepo tkRepository.ActivityRecordCmdRepo
}

func NewUpdatePhpModules(
	runtimeQueryRepo repository.RuntimeQueryRepo,
	runtimeCmdRepo repository.RuntimeCmdRepo,
	vhostQueryRepo repository.VirtualHostQueryRepo,
	activityCmdRepo tkRepository.ActivityRecordCmdRepo,
) UpdatePhpModules {
	return UpdatePhpModules{
		runtimeQueryRepo:      runtimeQueryRepo,
		runtimeCmdRepo:        runtimeCmdRepo,
		vhostQueryRepo:        vhostQueryRepo,
		activityRecordCmdRepo: activityCmdRepo,
	}
}

func (uc UpdatePhpModules) phpModuleEntitiesFactory(
	moduleUpdates []dto.PhpModuleUpdate,
) []entity.PhpModule {
	modules := []entity.PhpModule{}
	for _, moduleUpdate := range moduleUpdates {
		modules = append(
			modules,
			entity.NewPhpModule(moduleUpdate.Name, moduleUpdate.Status),
		)
	}

	return modules
}

func (uc UpdatePhpModules) discardDuplicatePhpModules(
	requestedModules []entity.PhpModule,
) []entity.PhpModule {
	requestedModuleNames := map[string]struct{}{}
	uniqueModules := []entity.PhpModule{}
	for _, module := range requestedModules {
		moduleName := module.Name.String()
		if _, exists := requestedModuleNames[moduleName]; exists {
			slog.Debug("DuplicatePhpModule", slog.String("moduleName", moduleName))
			continue
		}

		requestedModuleNames[moduleName] = struct{}{}
		uniqueModules = append(uniqueModules, module)
	}

	return uniqueModules
}

func (uc UpdatePhpModules) ensurePhpModuleUpdatesAreAllowed(
	phpVersion valueObject.PhpVersion,
	requestedModules []entity.PhpModule,
) ([]dto.PhpModuleUpdateFailure, error) {
	supportedPhpModules, err := uc.runtimeQueryRepo.ReadPhpModules(phpVersion)
	if err != nil {
		slog.Error("ReadPhpModulesError", slog.String("err", err.Error()))
		return nil, errors.New("ReadPhpModulesInfraError")
	}

	supportedModuleNames := map[string]struct{}{}
	for _, module := range supportedPhpModules {
		supportedModuleNames[module.Name.String()] = struct{}{}
	}

	unsupportedModuleFailures := []dto.PhpModuleUpdateFailure{}
	for _, module := range requestedModules {
		moduleName := module.Name.String()
		if _, isSupported := supportedModuleNames[moduleName]; isSupported {
			continue
		}

		slog.Warn(
			"PhpModuleNotSupportedForVersion",
			slog.String("moduleName", moduleName),
			slog.String("phpVersion", phpVersion.String()),
		)
		failureReason := valueObject.NewFailureReason(
			repository.ErrPhpModuleNotSupportedForVersion.Error(),
		)
		unsupportedModuleFailures = append(
			unsupportedModuleFailures,
			dto.NewPhpModuleUpdateFailure(
				module.Name, module.Status, failureReason,
			),
		)
	}

	if len(unsupportedModuleFailures) > 0 {
		return unsupportedModuleFailures, repository.ErrPhpModuleNotSupportedForVersion
	}

	return nil, nil
}

func (uc UpdatePhpModules) Execute(
	requestDto dto.UpdatePhpModulesRequest,
) (responseDto dto.UpdatePhpModulesResponse, err error) {
	if !useCaseHelper.IsPhpVersionInstalled(
		uc.runtimeQueryRepo, requestDto.PhpVersion,
	) {
		return responseDto, errors.New("PhpVersionNotInstalled")
	}

	_, err = uc.vhostQueryRepo.ReadFirst(dto.ReadVirtualHostsRequest{
		Hostname: &requestDto.Hostname,
	})
	if err != nil {
		slog.Error("VirtualHostNotFound", slog.String("err", err.Error()))
		return responseDto, errors.New("VirtualHostNotFound")
	}

	currentPhpVersion, err := uc.runtimeQueryRepo.ReadPhpVersion(
		requestDto.Hostname,
	)
	if err != nil {
		slog.Error("ReadCurrentPhpVersionError", slog.String("err", err.Error()))
		return responseDto, errors.New("ReadCurrentPhpVersionInfraError")
	}
	if currentPhpVersion.Value != requestDto.PhpVersion {
		return responseDto, repository.ErrPhpVersionChanged
	}

	requestedModules := uc.phpModuleEntitiesFactory(requestDto.PhpModules)
	uniqueModules := uc.discardDuplicatePhpModules(requestedModules)
	unsupportedModuleFailures, err := uc.ensurePhpModuleUpdatesAreAllowed(
		requestDto.PhpVersion, uniqueModules,
	)
	if err != nil {
		responseDto.FailedModulesWithReason = unsupportedModuleFailures
		return responseDto, err
	}

	responseDto, err = uc.runtimeCmdRepo.UpdatePhpModules(
		requestDto.Hostname, requestDto.PhpVersion, uniqueModules,
	)
	if err != nil {
		if errors.Is(err, repository.ErrPhpVersionChanged) {
			return responseDto, err
		}
		slog.Error("UpdatePhpModulesError", slog.String("err", err.Error()))
		return responseDto, errors.New("UpdatePhpModulesInfraError")
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		UpdatePhpModules(requestDto, responseDto)

	return responseDto, nil
}
