package servicesInfra

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/alessio/shellescape"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	dbModel "github.com/speedianet/os/src/infra/internalDatabase/model"
)

type ServicesCmdRepo struct {
	persistentDbSvc   *internalDbInfra.PersistentDatabaseService
	servicesQueryRepo *ServicesQueryRepo
}

func NewServicesCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *ServicesCmdRepo {
	return &ServicesCmdRepo{
		persistentDbSvc:   persistentDbSvc,
		servicesQueryRepo: NewServicesQueryRepo(persistentDbSvc),
	}
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
	err := repo.Stop(name)
	if err != nil {
		return err
	}

	time.Sleep(1 * time.Second)

	return repo.Start(name)
}

func (repo *ServicesCmdRepo) Reload() error {
	_, err := infraHelper.RunCmd("supervisorctl", "update")
	return err
}

func (repo *ServicesCmdRepo) replaceCmdStepsPlaceholders(
	cmdSteps []valueObject.UnixCommand,
	placeholders map[string]string,
) (usableCmdSteps []valueObject.UnixCommand, err error) {
	if len(cmdSteps) == 0 {
		return usableCmdSteps, nil
	}

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

		usableCmdStep, err := valueObject.NewUnixCommand(cmdStepStr)
		if err != nil {
			return nil, errors.New("InvalidCmdStep: " + cmdStepStr)
		}

		usableCmdSteps = append(usableCmdSteps, usableCmdStep)
	}

	return usableCmdSteps, nil
}

func (repo *ServicesCmdRepo) CreateInstallable(
	createDto dto.CreateInstallableService,
) (installedServiceName valueObject.ServiceName, err error) {
	installableService, err := repo.servicesQueryRepo.ReadInstallableByName(createDto.Name)
	if err != nil {
		return installedServiceName, err
	}

	if installableService.Nature.String() == "multi" {
		if createDto.StartupFile == nil {
			if installableService.StartupFile == nil {
				return installedServiceName, errors.New("MissingStartupFile")
			}
			createDto.StartupFile = installableService.StartupFile
		}

		startupFileHash := infraHelper.GenStrongShortHash(createDto.StartupFile.String())
		createDto.Name, err = valueObject.NewServiceName(
			createDto.Name.String() + "-" + startupFileHash,
		)
		if err != nil {
			return installedServiceName, errors.New(
				"AddFileHashNameSuffixError: " + err.Error(),
			)
		}
	}
	installedServiceName = createDto.Name

	serviceVersion := installableService.Versions[0]
	if createDto.Version != nil {
		serviceVersion = *createDto.Version
	}

	primaryHostname, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return installedServiceName, err
	}

	stepsPlaceholders := map[string]string{
		"randomPassword":  infraHelper.GenPass(16),
		"version":         serviceVersion.String(),
		"primaryHostname": primaryHostname.String(),
	}

	if createDto.StartupFile != nil {
		stepsPlaceholders["startupFile"] = createDto.StartupFile.String()
	}

	installedServiceNameStr := installedServiceName.String()
	defaultDirectories := []string{"conf", "logs"}
	for _, defaultDir := range defaultDirectories {
		err = infraHelper.MakeDir("/app/" + defaultDir + "/" + installedServiceNameStr)
		if err != nil {
			return installedServiceName, errors.New("CreateDefaultDirsError: " + err.Error())
		}
	}

	usableInstallCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		installableService.InstallCmdSteps, stepsPlaceholders,
	)
	if err != nil {
		return installedServiceName, err
	}

	for stepIndex, cmdStep := range usableInstallCmdSteps {
		_, err = infraHelper.RunCmdWithSubShell(cmdStep.String())
		if err != nil {
			stepIndexStr := strconv.Itoa(stepIndex)
			return installedServiceName, errors.New(
				"RunCmdStepError (" + stepIndexStr + "): " + err.Error(),
			)
		}
	}

	startCmdSteps := []valueObject.UnixCommand{installableService.StartCmd}
	usableCmdSteps := map[string][]valueObject.UnixCommand{
		"start":     startCmdSteps,
		"stop":      installableService.StopCmdSteps,
		"preStart":  installableService.PreStartCmdSteps,
		"postStart": installableService.PostStartCmdSteps,
		"preStop":   installableService.PreStopCmdSteps,
		"postStop":  installableService.PostStopCmdSteps,
	}
	for cmdStepType, cmdSteps := range usableCmdSteps {
		usableCmdSteps[cmdStepType], err = repo.replaceCmdStepsPlaceholders(
			cmdSteps, stepsPlaceholders,
		)
		if err != nil {
			return installedServiceName, err
		}
	}

	usableStartCmdSteps := usableCmdSteps["start"]
	if len(usableStartCmdSteps) == 0 {
		return installedServiceName, errors.New("MissingStartCmdStep")
	}
	usableStartCmd := usableStartCmdSteps[0]

	installedServiceModel := dbModel.NewInstalledService(
		installedServiceNameStr, installableService.Nature.String(),
		installableService.Type.String(), serviceVersion.String(),
		usableStartCmd.String(), createDto.Envs, createDto.PortBindings,
		usableCmdSteps["stop"], usableCmdSteps["preStart"], usableCmdSteps["postStart"],
		usableCmdSteps["preStop"], usableCmdSteps["postStop"], nil, nil, nil,
		createDto.AutoStart, createDto.AutoRestart, createDto.TimeoutStartSecs,
		createDto.MaxStartRetries, nil, nil,
	)

	if installableService.ExecUser != nil {
		execUserStr := installableService.ExecUser.String()
		installedServiceModel.ExecUser = &execUserStr
	}

	if installableService.WorkingDirectory != nil {
		workingDirectoryStr := installableService.WorkingDirectory.String()
		installedServiceModel.WorkingDirectory = &workingDirectoryStr
	}

	if createDto.StartupFile != nil {
		startupFileStr := createDto.StartupFile.String()
		installedServiceModel.StartupFile = &startupFileStr
	}

	if installableService.LogOutputPath != nil {
		logOutputPathStr := installableService.LogOutputPath.String()
		installedServiceModel.LogOutputPath = &logOutputPathStr
	}

	if installableService.LogErrorPath != nil {
		logErrorPathStr := installableService.LogErrorPath.String()
		installedServiceModel.LogErrorPath = &logErrorPathStr
	}

	err = repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return installedServiceName, err
	}

	return installedServiceName, repo.Reload()
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
		createDto.StopCmdSteps,
		createDto.PreStartCmdSteps,
		createDto.PostStartCmdSteps,
		createDto.PreStopCmdSteps,
		createDto.PostStopCmdSteps,
		nil,
		nil,
		nil,
		createDto.AutoStart,
		createDto.AutoRestart,
		createDto.TimeoutStartSecs,
		createDto.MaxStartRetries,
		nil,
		nil,
	)

	if createDto.ExecUser != nil {
		execUserStr := createDto.ExecUser.String()
		installedServiceModel.ExecUser = &execUserStr
	}

	if createDto.WorkingDirectory != nil {
		workingDirectoryStr := createDto.WorkingDirectory.String()
		installedServiceModel.WorkingDirectory = &workingDirectoryStr
	}

	if createDto.LogOutputPath != nil {
		logOutputPathStr := createDto.LogOutputPath.String()
		installedServiceModel.LogOutputPath = &logOutputPathStr
	}

	if createDto.LogErrorPath != nil {
		logErrorPathStr := createDto.LogErrorPath.String()
		installedServiceModel.LogErrorPath = &logErrorPathStr
	}

	err := repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) Update(updateDto dto.UpdateService) error {
	serviceEntity, err := repo.servicesQueryRepo.ReadByName(updateDto.Name)
	if err != nil {
		return err
	}

	if updateDto.Status != nil {
		desiredStatusStr := updateDto.Status.String()
		isSameStatus := serviceEntity.Status.String() == desiredStatusStr
		if isSameStatus {
			return nil
		}

		switch desiredStatusStr {
		case "running":
			return repo.Start(updateDto.Name)
		case "stopped":
			return repo.Stop(updateDto.Name)
		case "restarting":
			return repo.Restart(updateDto.Name)
		default:
			return errors.New("InvalidStatus: " + desiredStatusStr)
		}
	}

	updateMap := map[string]interface{}{}
	if updateDto.Type != nil {
		updateMap["type"] = updateDto.Type.String()
	}

	if updateDto.StartCmd != nil {
		updateMap["startCmd"] = updateDto.StartCmd.String()
	}

	if updateDto.Version != nil {
		updateMap["version"] = updateDto.Version.String()
	}

	if updateDto.StartupFile != nil {
		startupFileStr := updateDto.StartupFile.String()
		updateMap["startupFile"] = &startupFileStr
	}

	if updateDto.Envs != nil {
		envsStr := ""
		for _, env := range updateDto.Envs {
			envsStr += env.String() + ";"
		}
		envsStr = strings.TrimSuffix(envsStr, ";")
		updateMap["envs"] = &envsStr
	}

	if updateDto.PortBindings != nil {
		portBindingsStr := ""
		for _, portBinding := range updateDto.PortBindings {
			portBindingsStr += portBinding.String() + ";"
		}
		portBindingsStr = strings.TrimSuffix(portBindingsStr, ";")
		updateMap["portBindings"] = &portBindingsStr
	}

	if updateDto.AutoStart != nil {
		updateMap["autoStart"] = updateDto.AutoStart
	}

	if updateDto.TimeoutStartSecs != nil {
		updateMap["timeoutStartSecs"] = updateDto.TimeoutStartSecs
	}

	if updateDto.AutoRestart != nil {
		updateMap["autoRestart"] = updateDto.AutoRestart
	}

	if updateDto.MaxStartRetries != nil {
		updateMap["maxStartRetries"] = updateDto.MaxStartRetries
	}

	err = repo.persistentDbSvc.Handler.
		Model(&dbModel.InstalledService{}).
		Where("name = ?", updateDto.Name.String()).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	return repo.Reload()
}

func (repo *ServicesCmdRepo) Delete(name valueObject.ServiceName) error {
	err := repo.Stop(name)
	if err != nil {
		return err
	}

	err = repo.persistentDbSvc.Handler.
		Where("name = ?", name.String()).
		Delete(dbModel.InstalledService{}).Error
	if err != nil {
		return err
	}

	return repo.Reload()
}
