package servicesInfra

import (
	"errors"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type ServicesCmdRepo struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewServicesCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesCmdRepo {
	return &ServicesCmdRepo{persistentDbSvc: persistentDbSvc}
}

func (repo *ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd("supervisorctl", "start", name.String())
	return err
}

func (repo *ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd("supervisorctl", "stop", name.String())
	return err
}

func (repo *ServicesCmdRepo) Restart(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd("supervisorctl", "restart", name.String())
	return err
}

func (repo *ServicesCmdRepo) Reload() error {
	_, err := infraHelper.RunCmd("supervisorctl", "reload")
	return err
}

func (repo *ServicesCmdRepo) replaceCmdStepsPlaceholders(
	cmdSteps []valueObject.UnixCommand,
	placeholders map[string]string,
) (finalCmdSteps []valueObject.UnixCommand, err error) {
	for _, cmdStep := range cmdSteps {
		cmdStepStr := cmdStep.String()
		stepPlaceholders, _ := infraHelper.GetAllRegexGroupMatches(cmdStepStr, `%(.*?)%`)

		for _, stepPlaceholder := range stepPlaceholders {
			placeholderValue, exists := placeholders[stepPlaceholder]
			if !exists {
				return nil, errors.New("MissingPlaceholder: " + stepPlaceholder)
			}

			escapedPlaceholderValue := shellescape.Quote(placeholderValue)

			cmdStepStr = strings.ReplaceAll(
				cmdStepStr, "%"+stepPlaceholder+"%", escapedPlaceholderValue,
			)
		}

		finalCmdStep, err := valueObject.NewUnixCommand(cmdStepStr)
		if err != nil {
			return nil, errors.New("InvalidCmdStep: " + cmdStepStr)
		}

		finalCmdSteps = append(finalCmdSteps, finalCmdStep)
	}

	return finalCmdSteps, nil
}

func (repo *ServicesCmdRepo) CreateInstallable(
	createDto dto.CreateInstallableService,
) error {
	servicesQueryRepo := NewServicesQueryRepo(repo.persistentDbSvc)
	installableService, err := servicesQueryRepo.ReadInstallableByName(createDto.Name)
	if err != nil {
		return err
	}

	if installableService.Nature.String() == "multi" && createDto.StartupFile == nil {
		return errors.New("MultiNatureServicesRequiresStartupFile")
	}

	serviceVersion := installableService.Versions[0]
	if createDto.Version != nil {
		serviceVersion = *createDto.Version
	}

	primaryHostname, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return err
	}

	stepsPlaceholders := map[string]string{
		"randomPassword":  infraHelper.GenPass(16),
		"version":         serviceVersion.String(),
		"primaryHostname": primaryHostname.String(),
	}

	if createDto.StartupFile != nil {
		stepsPlaceholders["startupFile"] = createDto.StartupFile.String()
	}

	serviceNameStr := createDto.Name.String()
	defaultDirectories := []string{"conf", "logs"}
	for _, defaultDir := range defaultDirectories {
		err = infraHelper.MakeDir("/app/" + defaultDir + "/" + serviceNameStr)
		if err != nil {
			return errors.New("CreateDefaultDirsError: " + err.Error())
		}
	}

	finalInstallCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		installableService.InstallCmdSteps, stepsPlaceholders,
	)
	if err != nil {
		return err
	}

	for stepIndex, cmdStep := range finalInstallCmdSteps {
		_, err = infraHelper.RunCmdWithSubShell(cmdStep.String())
		if err != nil {
			stepIndexStr := strconv.Itoa(stepIndex)
			return errors.New(
				"RunCmdStepError (" + stepIndexStr + "): " + err.Error(),
			)
		}
	}

	finalStartCmd, err := repo.replaceCmdStepsPlaceholders(
		[]valueObject.UnixCommand{installableService.StartCmd}, stepsPlaceholders,
	)
	if err != nil {
		return err
	}

	if len(createDto.PortBindings) == 0 {
		createDto.PortBindings = installableService.PortBindings
	}

	installedServiceModel := dbModel.NewInstalledService(
		createDto.Name.String(),
		installableService.Nature.String(),
		installableService.Type.String(),
		serviceVersion.String(),
		finalStartCmd[0].String(),
		createDto.Envs,
		createDto.PortBindings,
		nil,
		createDto.AutoStart,
		createDto.TimeoutStartSecs,
		createDto.AutoRestart,
		createDto.MaxStartRetries,
	)

	if createDto.StartupFile != nil {
		startupFileStr := createDto.StartupFile.String()
		installedServiceModel.StartupFile = &startupFileStr
	}

	err = repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) CreateCustom(createDto dto.CreateCustomService) error {
	customNature, _ := valueObject.NewServiceNature("custom")

	installedServiceModel := dbModel.NewInstalledService(
		createDto.Name.String(),
		customNature.String(),
		createDto.Type.String(),
		createDto.Version.String(),
		createDto.StartCmd.String(),
		createDto.Envs,
		createDto.PortBindings,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	err := repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) Update(updateDto dto.UpdateService) error {
	return repo.Reload()
}

func (repo *ServicesCmdRepo) Delete(name valueObject.ServiceName) error {
	return repo.Reload()
}
