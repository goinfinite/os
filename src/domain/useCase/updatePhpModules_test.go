package useCase

import (
	"errors"
	"testing"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkRepository "github.com/goinfinite/tk/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type phpModulesRuntimeCmdRepo struct {
	responseDto            dto.UpdatePhpModulesResponse
	updatedModules         []entity.PhpModule
	updateCallMade         bool
	updateVersionCallMade  bool
	updateSettingsCallMade bool
}

func (repo *phpModulesRuntimeCmdRepo) CreatePhpVirtualHost(
	hostname tkValueObject.Fqdn,
) error {
	return nil
}

func (repo *phpModulesRuntimeCmdRepo) RunPhpCommand(
	requestDto dto.RunPhpCommandRequest,
) (dto.RunPhpCommandResponse, error) {
	return dto.RunPhpCommandResponse{}, nil
}

func (repo *phpModulesRuntimeCmdRepo) UpdatePhpVersion(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
) error {
	repo.updateVersionCallMade = true
	return nil
}

func (repo *phpModulesRuntimeCmdRepo) UpdatePhpSettings(
	hostname tkValueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	repo.updateSettingsCallMade = true
	return nil
}

func (repo *phpModulesRuntimeCmdRepo) UpdatePhpModules(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	modules []entity.PhpModule,
) (dto.UpdatePhpModulesResponse, error) {
	repo.updateCallMade = true
	repo.updatedModules = modules
	return repo.responseDto, nil
}

type phpModulesActivityRecordRepo struct{}

func (phpModulesActivityRecordRepo) Create(
	createDto tkDto.CreateActivityRecord,
) error {
	return nil
}

func (phpModulesActivityRecordRepo) Delete(
	deleteDto tkDto.DeleteActivityRecord,
) error {
	return nil
}

var _ tkRepository.ActivityRecordCmdRepo = phpModulesActivityRecordRepo{}

type installedPhpVersionReader struct{}

func (reader installedPhpVersionReader) ReadPhpVersionsInstalled() (
	[]valueObject.PhpVersion, error,
) {
	phpVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		return []valueObject.PhpVersion{}, err
	}

	return []valueObject.PhpVersion{phpVersion}, nil
}

func (reader installedPhpVersionReader) ReadPhpConfigs(
	hostname tkValueObject.Fqdn,
) (entity.PhpConfigs, error) {
	return entity.PhpConfigs{}, nil
}

func (reader installedPhpVersionReader) ReadPhpVersion(
	hostname tkValueObject.Fqdn,
) (entity.PhpVersion, error) {
	phpVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		return entity.PhpVersion{}, err
	}

	return entity.NewPhpVersion(phpVersion, []valueObject.PhpVersion{phpVersion}), nil
}

func (reader installedPhpVersionReader) ReadPhpModules(
	phpVersion valueObject.PhpVersion,
) ([]entity.PhpModule, error) {
	return []entity.PhpModule{}, nil
}

type supportedPhpModuleReader struct {
	installedPhpVersionReader
}

func (reader supportedPhpModuleReader) ReadPhpModules(
	phpVersion valueObject.PhpVersion,
) ([]entity.PhpModule, error) {
	phpModuleName, err := valueObject.NewPhpModuleName("ssh2")
	if err != nil {
		return []entity.PhpModule{}, err
	}

	return []entity.PhpModule{
		entity.NewPhpModule(phpModuleName, false),
	}, nil
}

type phpModulesReadFailureReader struct {
	installedPhpVersionReader
}

func (reader phpModulesReadFailureReader) ReadPhpModules(
	phpVersion valueObject.PhpVersion,
) ([]entity.PhpModule, error) {
	return nil, errors.New("ReadPhpModulesFailed")
}

type existingVirtualHostReader struct{}

func (reader existingVirtualHostReader) Read(
	requestDto dto.ReadVirtualHostsRequest,
) (dto.ReadVirtualHostsResponse, error) {
	return dto.ReadVirtualHostsResponse{}, nil
}

func (reader existingVirtualHostReader) ReadFirst(
	requestDto dto.ReadVirtualHostsRequest,
) (entity.VirtualHost, error) {
	return entity.VirtualHost{}, nil
}

func (reader existingVirtualHostReader) ReadFirstWithMappings(
	requestDto dto.ReadVirtualHostsRequest,
) (dto.VirtualHostWithMappings, error) {
	return dto.VirtualHostWithMappings{}, nil
}

type installedRuntimeReader struct {
	supportedPhpModuleReader
}

func (reader installedRuntimeReader) ReadPhpVersionsInstalled() (
	[]valueObject.PhpVersion, error,
) {
	installedVersions := []valueObject.PhpVersion{}
	for _, versionString := range []string{"5.6", "7.4"} {
		version, err := valueObject.NewPhpVersion(versionString)
		if err != nil {
			return installedVersions, err
		}
		installedVersions = append(installedVersions, version)
	}

	return installedVersions, nil
}

func (reader installedRuntimeReader) ReadPhpVersion(
	hostname tkValueObject.Fqdn,
) (entity.PhpVersion, error) {
	currentVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		return entity.PhpVersion{}, err
	}

	return entity.NewPhpVersion(
		currentVersion,
		[]valueObject.PhpVersion{currentVersion},
	), nil
}

func TestUpdatePhpModulesReturnsReport(test *testing.T) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	phpVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	moduleName, err := valueObject.NewPhpModuleName("ssh2")
	if err != nil {
		test.Fatalf("PhpModuleNameCreationFailed: %v", err)
	}

	failureReason, err := valueObject.NewFailureReason("ModuleInstallFailed")
	if err != nil {
		test.Fatalf("FailureReasonCreationFailed: %v", err)
	}
	expectedResponse := dto.NewUpdatePhpModulesResponse(
		[]dto.PhpModuleUpdate{},
		[]dto.PhpModuleUpdateFailure{
			dto.NewPhpModuleUpdateFailure(moduleName, true, failureReason),
		},
	)
	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{responseDto: expectedResponse}
	request := dto.NewUpdatePhpModulesRequest(
		hostname,
		phpVersion,
		[]dto.PhpModuleUpdate{dto.NewPhpModuleUpdate(moduleName, true)},
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)

	updatePhpModulesUseCase := NewUpdatePhpModules(
		supportedPhpModuleReader{},
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)
	response, err := updatePhpModulesUseCase.Execute(request)
	if err != nil {
		test.Fatalf("UpdatePhpModulesFailed: %v", err)
	}
	if !runtimeCmdRepo.updateCallMade {
		test.Fatal("RuntimeModulesUpdateWasNotCalled")
	}
	if len(response.FailedModulesWithReason) != 1 {
		test.Fatalf(
			"FailedModulesCountMismatch: expected 1, got %d",
			len(response.FailedModulesWithReason),
		)
	}
	if response.FailedModulesWithReason[0].Reason != failureReason {
		test.Errorf("FailureReasonMismatch")
	}
}

func TestUpdatePhpModulesRejectsChangedPhpVersion(test *testing.T) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	targetVersion, err := valueObject.NewPhpVersion("7.4")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	moduleName, err := valueObject.NewPhpModuleName("ssh2")
	if err != nil {
		test.Fatalf("PhpModuleNameCreationFailed: %v", err)
	}
	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{}
	request := dto.NewUpdatePhpModulesRequest(
		hostname,
		targetVersion,
		[]dto.PhpModuleUpdate{dto.NewPhpModuleUpdate(moduleName, true)},
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)

	updatePhpModulesUseCase := NewUpdatePhpModules(
		installedRuntimeReader{},
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)
	_, err = updatePhpModulesUseCase.Execute(request)
	if !errors.Is(err, repository.ErrPhpVersionChanged) {
		test.Errorf("UnexpectedError: %v", err)
	}
	if runtimeCmdRepo.updateCallMade {
		test.Error("RuntimeModulesUpdateShouldNotBeCalled")
	}
}

func TestUpdatePhpModulesReportsUnsupportedModules(test *testing.T) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	phpVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	supportedModuleName, err := valueObject.NewPhpModuleName("ssh2")
	if err != nil {
		test.Fatalf("SupportedPhpModuleNameCreationFailed: %v", err)
	}
	firstUnsupportedModuleName, err := valueObject.NewPhpModuleName("curl")
	if err != nil {
		test.Fatalf("FirstUnsupportedPhpModuleNameCreationFailed: %v", err)
	}
	secondUnsupportedModuleName, err := valueObject.NewPhpModuleName("mysqli")
	if err != nil {
		test.Fatalf("SecondUnsupportedPhpModuleNameCreationFailed: %v", err)
	}
	firstUnsupportedModuleUpdate := dto.NewPhpModuleUpdate(
		firstUnsupportedModuleName, true,
	)
	secondUnsupportedModuleUpdate := dto.NewPhpModuleUpdate(
		secondUnsupportedModuleName, false,
	)
	unsupportedModuleUpdates := []dto.PhpModuleUpdate{
		firstUnsupportedModuleUpdate,
		secondUnsupportedModuleUpdate,
	}
	moduleUpdates := []dto.PhpModuleUpdate{
		dto.NewPhpModuleUpdate(supportedModuleName, true),
		firstUnsupportedModuleUpdate,
		dto.NewPhpModuleUpdate(firstUnsupportedModuleName, false),
		secondUnsupportedModuleUpdate,
	}
	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{}
	request := dto.NewUpdatePhpModulesRequest(
		hostname,
		phpVersion,
		moduleUpdates,
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)
	updatePhpModulesUseCase := NewUpdatePhpModules(
		supportedPhpModuleReader{},
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)

	response, err := updatePhpModulesUseCase.Execute(request)
	if !errors.Is(err, repository.ErrPhpModuleNotSupportedForVersion) {
		test.Fatalf("UnexpectedError: %v", err)
	}
	if len(response.FailedModulesWithReason) != len(unsupportedModuleUpdates) {
		test.Fatalf(
			"FailedModulesCountMismatch: expected %d, got %d",
			len(unsupportedModuleUpdates),
			len(response.FailedModulesWithReason),
		)
	}

	for failureIndex, expectedFailure := range unsupportedModuleUpdates {
		failure := response.FailedModulesWithReason[failureIndex]
		if failure.Name != expectedFailure.Name {
			test.Errorf("FailedModuleNameMismatch: index %d", failureIndex)
		}
		if failure.Status != expectedFailure.Status {
			test.Errorf("FailedModuleStatusMismatch: index %d", failureIndex)
		}
		if failure.Reason.String() != repository.ErrPhpModuleNotSupportedForVersion.Error() {
			test.Errorf("FailureReasonMismatch: index %d", failureIndex)
		}
	}
	if runtimeCmdRepo.updateCallMade {
		test.Error("RuntimeModulesUpdateShouldNotBeCalled")
	}
}

func TestUpdatePhpModulesDoesNotReportFailuresWhenModulesCannotBeRead(
	test *testing.T,
) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	phpVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	moduleName, err := valueObject.NewPhpModuleName("ssh2")
	if err != nil {
		test.Fatalf("PhpModuleNameCreationFailed: %v", err)
	}

	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{}
	request := dto.NewUpdatePhpModulesRequest(
		hostname,
		phpVersion,
		[]dto.PhpModuleUpdate{dto.NewPhpModuleUpdate(moduleName, true)},
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)
	updatePhpModulesUseCase := NewUpdatePhpModules(
		phpModulesReadFailureReader{},
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)

	response, err := updatePhpModulesUseCase.Execute(request)
	if err == nil {
		test.Fatal("MissingExpectedError")
	}
	if err.Error() != "ReadPhpModulesInfraError" {
		test.Errorf("UnexpectedError: %v", err)
	}
	if len(response.FailedModulesWithReason) != 0 {
		test.Errorf(
			"FailedModulesCountMismatch: expected 0, got %d",
			len(response.FailedModulesWithReason),
		)
	}
	if runtimeCmdRepo.updateCallMade {
		test.Error("RuntimeModulesUpdateShouldNotBeCalled")
	}
}

func TestUpdatePhpModulesDiscardsDuplicateModules(test *testing.T) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	phpVersion, err := valueObject.NewPhpVersion("5.6")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	moduleName, err := valueObject.NewPhpModuleName("ssh2")
	if err != nil {
		test.Fatalf("PhpModuleNameCreationFailed: %v", err)
	}

	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{}
	request := dto.NewUpdatePhpModulesRequest(
		hostname,
		phpVersion,
		[]dto.PhpModuleUpdate{
			dto.NewPhpModuleUpdate(moduleName, true),
			dto.NewPhpModuleUpdate(moduleName, false),
		},
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)
	updatePhpModulesUseCase := NewUpdatePhpModules(
		supportedPhpModuleReader{},
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)

	response, err := updatePhpModulesUseCase.Execute(request)
	if err != nil {
		test.Fatalf("UpdatePhpModulesFailed: %v", err)
	}
	if !runtimeCmdRepo.updateCallMade {
		test.Fatal("RuntimeModulesUpdateWasNotCalled")
	}
	if len(response.FailedModulesWithReason) != 0 {
		test.Errorf(
			"FailedModulesCountMismatch: expected 0, got %d",
			len(response.FailedModulesWithReason),
		)
	}
	if len(runtimeCmdRepo.updatedModules) != 1 {
		test.Fatalf(
			"UpdatedModulesCountMismatch: expected 1, got %d",
			len(runtimeCmdRepo.updatedModules),
		)
	}

	updatedModule := runtimeCmdRepo.updatedModules[0]
	if updatedModule.Name != moduleName {
		test.Error("UpdatedModuleNameMismatch")
	}
	if !updatedModule.Status {
		test.Error("UpdatedModuleStatusMismatch")
	}
}

func TestUpdatePhpVersionUpdatesVersion(test *testing.T) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	phpVersion, err := valueObject.NewPhpVersion("7.4")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{}
	request := dto.NewUpdatePhpVersionRequest(
		hostname,
		phpVersion,
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)

	updatePhpVersionUseCase := NewUpdatePhpVersion(
		installedRuntimeReader{},
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)
	err = updatePhpVersionUseCase.Execute(request)
	if err != nil {
		test.Fatalf("UpdatePhpVersionFailed: %v", err)
	}
	if !runtimeCmdRepo.updateVersionCallMade {
		test.Fatal("PhpVersionUpdateWasNotCalled")
	}
}

func TestUpdatePhpSettingsUpdatesSettings(test *testing.T) {
	hostname, err := tkValueObject.NewFqdn("example.com")
	if err != nil {
		test.Fatalf("HostnameCreationFailed: %v", err)
	}
	setting, err := entity.NewPhpSettingFromString("memory_limit:128M")
	if err != nil {
		test.Fatalf("PhpSettingCreationFailed: %v", err)
	}
	runtimeCmdRepo := &phpModulesRuntimeCmdRepo{}
	request := dto.NewUpdatePhpSettingsRequest(
		hostname,
		[]entity.PhpSetting{setting},
		tkValueObject.AccountIdSystem,
		tkValueObject.IpAddressLocal,
	)

	updatePhpSettingsUseCase := NewUpdatePhpSettings(
		runtimeCmdRepo,
		existingVirtualHostReader{},
		phpModulesActivityRecordRepo{},
	)
	err = updatePhpSettingsUseCase.Execute(request)
	if err != nil {
		test.Fatalf("UpdatePhpSettingsFailed: %v", err)
	}
	if !runtimeCmdRepo.updateSettingsCallMade {
		test.Fatal("PhpSettingsUpdateWasNotCalled")
	}
}
