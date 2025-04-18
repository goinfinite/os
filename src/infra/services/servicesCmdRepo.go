package servicesInfra

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	dbModel "github.com/goinfinite/os/src/infra/internalDatabase/model"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

const SupervisorCtlBin string = "/usr/bin/supervisorctl -c /infinite/supervisord.conf"

var defaultServiceDirectories []string = []string{"conf", "logs"}

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

func (repo *ServicesCmdRepo) runCmdSteps(
	stepsType string,
	steps []valueObject.UnixCommand,
	totalExecTimeoutSecs valueObject.UnixTime,
) error {
	if len(steps) == 0 {
		return nil
	}

	totalExecTimeoutSecsUint := uint64(totalExecTimeoutSecs.Int64())
	runCmdSettings := infraHelper.RunCmdSettings{
		ShouldRunWithSubShell: true,
		ExecutionTimeoutSecs:  totalExecTimeoutSecsUint,
	}

	totalExecRemainingTime := totalExecTimeoutSecsUint
	for stepIndex, step := range steps {
		stepStr := step.String()

		slog.Debug("Running"+stepsType+"Step", slog.String("step", stepStr))

		runCmdSettings.Command = stepStr

		execTimeStart := time.Now()
		stepOutput, err := infraHelper.RunCmd(runCmdSettings)
		if err != nil {
			errorMessage := stepOutput + " | " + err.Error()
			if infraHelper.IsRunCmdTimeout(err) {
				errorMessage = "Service" + stepsType + "TimeoutExceeded"
			}

			return fmt.Errorf(
				"%sCmdStepError (%s): %s",
				stepsType, strconv.Itoa(stepIndex), errorMessage,
			)
		}

		stepExecElapsedTimeSecs := uint64(time.Since(execTimeStart).Seconds())
		totalExecRemainingTime = totalExecRemainingTime - stepExecElapsedTimeSecs
		if totalExecRemainingTime == 0 {
			return errors.New("Service" + stepsType + "TimeoutExceeded")
		}

		runCmdSettings.ExecutionTimeoutSecs = totalExecRemainingTime
	}

	time.Sleep(1 * time.Second)

	return nil
}

func (repo *ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &name,
	}
	serviceEntity, err := repo.servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
	if err != nil {
		return err
	}

	err = repo.runCmdSteps(
		"PreStart", serviceEntity.PreStartCmdSteps, serviceEntity.PreStartTimeoutSecs,
	)
	if err != nil {
		return err
	}

	serviceNameStr := serviceEntity.Name.String()
	startOutput, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               SupervisorCtlBin + " start " + serviceNameStr,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		combinedOutput := startOutput + " " + err.Error()
		if !strings.Contains(combinedOutput, "no such process") {
			return errors.New("SupervisorStartError: " + combinedOutput)
		}

		addOutput, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command:               SupervisorCtlBin + " add " + serviceNameStr,
			ShouldRunWithSubShell: true,
		})
		if err != nil {
			combinedOutput = addOutput + " " + err.Error()
			return errors.New("SupervisorAddError: " + combinedOutput)
		}
	}

	time.Sleep(1 * time.Second)

	return repo.runCmdSteps(
		"PostStart", serviceEntity.PostStartCmdSteps, serviceEntity.PostStartTimeoutSecs,
	)
}

func (repo *ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &name,
	}
	serviceEntity, err := repo.servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
	if err != nil {
		return err
	}

	err = repo.runCmdSteps(
		"PreStop", serviceEntity.PreStopCmdSteps, serviceEntity.PreStopTimeoutSecs,
	)
	if err != nil {
		return err
	}

	stopOutput, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               SupervisorCtlBin + " stop " + serviceEntity.Name.String(),
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		combinedOutput := stopOutput + " " + err.Error()
		return errors.New("SupervisorStopError: " + combinedOutput)
	}

	time.Sleep(1 * time.Second)

	err = repo.runCmdSteps(
		"Stop", serviceEntity.StopCmdSteps, serviceEntity.StopTimeoutSecs,
	)
	if err != nil {
		return err
	}

	return repo.runCmdSteps(
		"PostStop", serviceEntity.PostStopCmdSteps, serviceEntity.PostStartTimeoutSecs,
	)
}

func (repo *ServicesCmdRepo) Restart(name valueObject.ServiceName) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &name,
	}
	serviceEntity, err := repo.servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
	if err != nil {
		return err
	}

	if serviceEntity.Status.String() == "running" {
		err = repo.Stop(name)
		if err != nil {
			return err
		}
	}

	return repo.Start(name)
}

func (repo *ServicesCmdRepo) updateProcessManagerConf() error {
	shouldIncludeMetrics := false
	readInstalledItemsDto := dto.ReadInstalledServicesItemsRequest{
		Pagination: dto.Pagination{
			ItemsPerPage: 1000,
		},
		ShouldIncludeMetrics: &shouldIncludeMetrics,
	}
	readInstalledItemsResponseDto, err := repo.servicesQueryRepo.ReadInstalledItems(
		readInstalledItemsDto,
	)
	if err != nil {
		return err
	}
	if len(readInstalledItemsResponseDto.InstalledServices) == 0 {
		return errors.New("NoServicesFoundToUpdateProcessManager")
	}

	ctlPassword := infraHelper.GenStrongShortHash(
		readInstalledItemsResponseDto.InstalledServices[0].CreatedAt.String(),
	)

	// cSpell:disable
	fileTemplate := `# AUTO GENERATED FILE. DO NOT EDIT.
[unix_http_server]
file=/run/supervisor.sock
chmod=0700
username=supervisord
password=` + ctlPassword + `

[supervisord]
nodaemon=true
user=root
directory=/infinite
logfile=/dev/stdout
logfile_maxbytes=0
loglevel=ERROR
pidfile=/run/supervisord.pid

[supervisorctl]
serverurl=unix:///run/supervisor.sock
username=supervisord
password=` + ctlPassword + `

[rpcinterface:supervisor]
supervisor.rpcinterface_factory=supervisor.rpcinterface:make_main_rpcinterface
{{ range . }}
[program:{{.Name}}]
command={{.StartCmd}}
user={{ or .ExecUser "root" }}
{{- if .WorkingDirectory}}
directory={{.WorkingDirectory}}
{{- end}}
autostart={{ or .AutoStart "true" }}
autorestart={{ or .AutoRestart "true" }}
startretries={{ or .MaxStartRetries "3" }}
startsecs={{ or .TimeoutStartSecs "3" }}
{{- if .LogOutputPath}}
stdout_logfile={{.LogOutputPath}}
{{- if eq (printf "%s" .LogOutputPath) "/dev/stdout"}}
stdout_logfile_maxbytes=0
{{- else}}
stdout_logfile_maxbytes=10MB
{{end}}
{{- else}}
stdout_logfile=/app/logs/{{.Name}}/{{.Name}}.log
stdout_logfile_maxbytes=10MB
{{- end}}
{{- if .LogErrorPath}}
stderr_logfile={{.LogErrorPath}}
{{- if eq (printf "%s" .LogErrorPath) "/dev/stderr"}}
stderr_logfile_maxbytes=0
{{- else}}
stderr_logfile_maxbytes=10MB
{{end}}
{{- else}}
stderr_logfile=/app/logs/{{.Name}}/{{.Name}}_error.log
stderr_logfile_maxbytes=10MB
{{- end}}
{{- if .Envs}}
environment={{range $index, $envVar := .Envs}}{{if $index}},{{end}}{{$envVar}}{{end}}
{{- end}}
{{end}}
`
	// cSpell:enable

	templatePtr, err := template.New("supervisorConf").Parse(fileTemplate)
	if err != nil {
		return errors.New("TemplateParsingError: " + err.Error())
	}

	var supervisorConfFileContent strings.Builder
	err = templatePtr.Execute(
		&supervisorConfFileContent, readInstalledItemsResponseDto.InstalledServices,
	)
	if err != nil {
		return errors.New("TemplateExecutionError: " + err.Error())
	}

	err = infraHelper.UpdateFile(
		"/infinite/supervisord.conf", supervisorConfFileContent.String(), true,
	)
	if err != nil {
		return err
	}

	reReadOutput, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               SupervisorCtlBin + " reread",
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		combinedOutput := reReadOutput + " " + err.Error()
		return errors.New("SupervisorRereadError: " + combinedOutput)
	}

	return nil
}

func (repo *ServicesCmdRepo) createDefaultDirectories(
	serviceName valueObject.ServiceName,
) error {
	for _, defaultDir := range defaultServiceDirectories {
		defaultDirPath := "/app/" + defaultDir + "/" + serviceName.String()

		err := infraHelper.MakeDir(defaultDirPath)
		if err != nil {
			return errors.New("CreateDefaultDirsError: " + err.Error())
		}

		deletionWarningFilePath := defaultDirPath + "/DONOTDELETE"
		if infraHelper.FileExists(deletionWarningFilePath) {
			continue
		}

		_, err = os.Create(deletionWarningFilePath)
		if err != nil {
			return errors.New("CreateDeletionWarningFileError: " + err.Error())
		}
	}

	return nil
}

func (repo *ServicesCmdRepo) updateDefaultDirectoriesPermissions(
	serviceName valueObject.ServiceName, execUser valueObject.UnixUsername,
) error {
	execUserStr := execUser.String()
	_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "id",
		Args:    []string{execUserStr},
	})
	if err != nil {
		return errors.New("EnsureExecUserExistenceError: " + err.Error())
	}

	for _, defaultDir := range defaultServiceDirectories {
		defaultDirPath := "/app/" + defaultDir + "/" + serviceName.String()

		_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command: "chown",
			Args:    []string{"-R", execUserStr, defaultDirPath},
		})
		if err != nil {
			return errors.New("ChownDefaultDirsError: " + err.Error())
		}
	}

	return nil
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
		stepPlaceholders := infraHelper.GetAllRegexGroupMatches(
			cmdStepStr, `%(\w{1,256})%`,
		)

		for _, stepPlaceholder := range stepPlaceholders {
			placeholderValue, exists := placeholders[stepPlaceholder]
			if !exists {
				slog.Debug(
					"MissingPlaceholder", slog.String("placeholder", stepPlaceholder),
				)
				placeholderValue = ""
			}

			printablePlaceholderValue := infraHelper.ShellEscape{}.StripUnsafe(placeholderValue)

			cmdStepStr = strings.ReplaceAll(
				cmdStepStr, "%"+stepPlaceholder+"%", printablePlaceholderValue,
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
	installableServiceEntity, err := repo.servicesQueryRepo.ReadFirstInstallableItem(
		dto.ReadInstallableServicesItemsRequest{ServiceName: &createDto.Name},
	)
	if err != nil {
		return installedServiceName, err
	}

	if installableServiceEntity.Nature == valueObject.ServiceNatureMulti {
		if createDto.StartupFile == nil {
			if installableServiceEntity.StartupFile == nil {
				return installedServiceName, errors.New("MultiNatureServiceRequiresStartupFile")
			}
			createDto.StartupFile = installableServiceEntity.StartupFile
		}

		startupFileHash := infraHelper.GenStrongShortHash(createDto.StartupFile.String())
		createDto.Name, err = valueObject.NewServiceName(
			createDto.Name.String() + "_" + startupFileHash,
		)
		if err != nil {
			return installedServiceName, errors.New(
				"AppendFileHashNameSuffixError: " + err.Error(),
			)
		}
	}
	installedServiceName = createDto.Name

	serviceVersion := installableServiceEntity.Versions[0]
	if createDto.Version != nil {
		serviceVersion = *createDto.Version
	}

	primaryHostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		return installedServiceName, err
	}

	synthesizer := tkInfra.Synthesizer{}
	stepsPlaceholders := map[string]string{
		"randomPassword":            synthesizer.PasswordFactory(16, false),
		"randomPasswordWithSymbols": synthesizer.PasswordFactory(16, true),
		"version":                   serviceVersion.String(),
		"primaryHostname":           primaryHostname.String(),
	}
	isSoloService := installableServiceEntity.Nature == valueObject.ServiceNatureSolo
	if createDto.StartupFile != nil && !isSoloService {
		stepsPlaceholders["startupFile"] = createDto.StartupFile.String()
	}

	installedServiceNameStr := strings.ToLower(installedServiceName.String())
	installablesAssetsDirPath := fmt.Sprintf(
		"%s/%s/%s/assets", infraEnvs.InstallableServicesItemsDir,
		installableServiceEntity.Type.String(), installedServiceNameStr,
	)
	stepsPlaceholders["installableServiceAssetsDirPath"] = installablesAssetsDirPath

	err = repo.createDefaultDirectories(installedServiceName)
	if err != nil {
		return installedServiceName, err
	}

	usableInstallCmdSteps, err := repo.replaceCmdStepsPlaceholders(
		installableServiceEntity.InstallCmdSteps, stepsPlaceholders,
	)
	if err != nil {
		return installedServiceName, err
	}

	err = repo.runCmdSteps(
		"Install", usableInstallCmdSteps, installableServiceEntity.InstallTimeoutSecs,
	)
	if err != nil {
		return installedServiceName, err
	}

	if installableServiceEntity.ExecUser != nil {
		err = repo.updateDefaultDirectoriesPermissions(
			installableServiceEntity.Name, *installableServiceEntity.ExecUser,
		)
		if err != nil {
			return installedServiceName, err
		}
	}

	startCmdSteps := []valueObject.UnixCommand{installableServiceEntity.StartCmd}
	usableCmdSteps := map[string][]valueObject.UnixCommand{
		"start":     startCmdSteps,
		"stop":      installableServiceEntity.StopCmdSteps,
		"preStart":  installableServiceEntity.PreStartCmdSteps,
		"postStart": installableServiceEntity.PostStartCmdSteps,
		"preStop":   installableServiceEntity.PreStopCmdSteps,
		"postStop":  installableServiceEntity.PostStopCmdSteps,
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

	if len(createDto.PortBindings) == 0 {
		createDto.PortBindings = installableServiceEntity.PortBindings
	}

	installedServiceModel := dbModel.NewInstalledService(
		installedServiceName.String(), installableServiceEntity.Nature.String(),
		installableServiceEntity.Type.String(), serviceVersion.String(),
		usableStartCmd.String(), createDto.Envs, createDto.PortBindings,
		usableCmdSteps["stop"], usableCmdSteps["preStart"], usableCmdSteps["postStart"],
		usableCmdSteps["preStop"], usableCmdSteps["postStop"], nil, nil, nil,
		createDto.AutoStart, createDto.AutoRestart, createDto.TimeoutStartSecs,
		createDto.MaxStartRetries, nil, nil, nil,
	)

	if installableServiceEntity.ExecUser != nil {
		execUserStr := installableServiceEntity.ExecUser.String()
		installedServiceModel.ExecUser = &execUserStr
	}

	if installableServiceEntity.WorkingDirectory != nil {
		workingDirectoryStr := installableServiceEntity.WorkingDirectory.String()
		installedServiceModel.WorkingDirectory = &workingDirectoryStr
	}
	if createDto.WorkingDir != nil {
		workingDirectoryStr := createDto.WorkingDir.String()
		installedServiceModel.WorkingDirectory = &workingDirectoryStr
	}

	if createDto.StartupFile != nil && isSoloService {
		startupFileStr := createDto.StartupFile.String()
		installedServiceModel.StartupFile = &startupFileStr
	}

	if installableServiceEntity.LogOutputPath != nil {
		logOutputPathStr := installableServiceEntity.LogOutputPath.String()
		installedServiceModel.LogOutputPath = &logOutputPathStr
	}

	if installableServiceEntity.LogErrorPath != nil {
		logErrorPathStr := installableServiceEntity.LogErrorPath.String()
		installedServiceModel.LogErrorPath = &logErrorPathStr
	}

	if installableServiceEntity.AvatarUrl != nil {
		avatarUrlStr := installableServiceEntity.AvatarUrl.String()
		installedServiceModel.AvatarUrl = &avatarUrlStr
	}

	err = repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return installedServiceName, err
	}

	err = repo.updateProcessManagerConf()
	if err != nil {
		return installedServiceName, err
	}

	return installedServiceName, repo.Start(installedServiceName)
}

func (repo *ServicesCmdRepo) CreateCustom(createDto dto.CreateCustomService) error {
	customNature, _ := valueObject.NewServiceNature("custom")

	installedServiceModel := dbModel.NewInstalledService(
		createDto.Name.String(), customNature.String(), createDto.Type.String(),
		createDto.Version.String(), createDto.StartCmd.String(),
		createDto.Envs, createDto.PortBindings, createDto.StopCmdSteps,
		createDto.PreStartCmdSteps, createDto.PostStartCmdSteps,
		createDto.PreStopCmdSteps, createDto.PostStopCmdSteps, nil, nil, nil,
		createDto.AutoStart, createDto.AutoRestart, createDto.TimeoutStartSecs,
		createDto.MaxStartRetries, nil, nil, nil,
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

	if createDto.AvatarUrl != nil {
		avatarUrlStr := createDto.AvatarUrl.String()
		installedServiceModel.AvatarUrl = &avatarUrlStr
	}

	err := repo.persistentDbSvc.Handler.Create(&installedServiceModel).Error
	if err != nil {
		return err
	}

	err = repo.createDefaultDirectories(createDto.Name)
	if err != nil {
		return err
	}

	if createDto.ExecUser != nil {
		err = repo.updateDefaultDirectoriesPermissions(
			createDto.Name, *createDto.ExecUser,
		)
		if err != nil {
			return err
		}
	}

	err = repo.updateProcessManagerConf()
	if err != nil {
		return err
	}

	return repo.Start(createDto.Name)
}

func (repo *ServicesCmdRepo) Update(updateDto dto.UpdateService) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &updateDto.Name,
	}
	serviceEntity, err := repo.servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
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

	installedServiceModel := dbModel.InstalledService{}
	updateMap := map[string]interface{}{}
	if updateDto.Type != nil {
		updateMap["type"] = updateDto.Type.String()
	}

	if updateDto.StartCmd != nil {
		updateMap["start_cmd"] = updateDto.StartCmd.String()
	}

	if updateDto.Version != nil {
		updateMap["version"] = updateDto.Version.String()
	}

	isSoloService := serviceEntity.Nature == valueObject.ServiceNatureSolo
	if updateDto.StartupFile != nil && isSoloService {
		startupFileStr := updateDto.StartupFile.String()
		updateMap["startup_file"] = &startupFileStr
	}

	if updateDto.Envs != nil {
		updateMap["envs"] = installedServiceModel.JoinEnvs(updateDto.Envs)
	}

	if updateDto.PortBindings != nil {
		updateMap["port_bindings"] = installedServiceModel.JoinPortBindings(
			updateDto.PortBindings,
		)
	}

	if updateDto.StopCmdSteps != nil {
		updateMap["stop_cmd_steps"] = installedServiceModel.JoinCmdSteps(
			updateDto.StopCmdSteps,
		)
	}

	if updateDto.PreStartCmdSteps != nil {
		updateMap["pre_start_cmd_steps"] = installedServiceModel.JoinCmdSteps(
			updateDto.PreStartCmdSteps,
		)
	}

	if updateDto.PostStartCmdSteps != nil {
		updateMap["post_start_cmd_steps"] = installedServiceModel.JoinCmdSteps(
			updateDto.PostStartCmdSteps,
		)
	}

	if updateDto.PreStopCmdSteps != nil {
		updateMap["pre_stop_cmd_steps"] = installedServiceModel.JoinCmdSteps(
			updateDto.PreStopCmdSteps,
		)
	}

	if updateDto.PostStopCmdSteps != nil {
		updateMap["post_stop_cmd_steps"] = installedServiceModel.JoinCmdSteps(
			updateDto.PostStopCmdSteps,
		)
	}

	if updateDto.ExecUser != nil {
		execUserStr := updateDto.ExecUser.String()
		updateMap["exec_user"] = &execUserStr
	}

	if updateDto.WorkingDirectory != nil {
		workingDirectoryStr := updateDto.WorkingDirectory.String()
		updateMap["working_directory"] = &workingDirectoryStr
	}

	if updateDto.StartupFile != nil && isSoloService {
		startupFileStr := updateDto.StartupFile.String()
		updateMap["startup_file"] = &startupFileStr
	}

	if updateDto.AutoStart != nil {
		updateMap["auto_start"] = updateDto.AutoStart
	}

	if updateDto.TimeoutStartSecs != nil {
		updateMap["timeout_start_secs"] = updateDto.TimeoutStartSecs
	}

	if updateDto.AutoRestart != nil {
		updateMap["auto_restart"] = updateDto.AutoRestart
	}

	if updateDto.MaxStartRetries != nil {
		updateMap["max_start_retries"] = updateDto.MaxStartRetries
	}

	if updateDto.LogOutputPath != nil {
		logOutputPathStr := updateDto.LogOutputPath.String()
		updateMap["log_output_path"] = &logOutputPathStr
	}

	if updateDto.LogErrorPath != nil {
		logErrorPathStr := updateDto.LogErrorPath.String()
		updateMap["log_error_path"] = &logErrorPathStr
	}

	if updateDto.AvatarUrl != nil && serviceEntity.Nature == valueObject.ServiceNatureCustom {
		avatarUrlStr := updateDto.AvatarUrl.String()
		updateMap["avatar_url"] = &avatarUrlStr
	}

	err = repo.persistentDbSvc.Handler.
		Model(&installedServiceModel).
		Where("name = ?", updateDto.Name.String()).
		Updates(updateMap).Error
	if err != nil {
		return err
	}

	err = repo.updateProcessManagerConf()
	if err != nil {
		return err
	}

	return repo.Restart(updateDto.Name)
}

func (repo *ServicesCmdRepo) Delete(name valueObject.ServiceName) error {
	readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
		ServiceName: &name,
	}
	serviceEntity, err := repo.servicesQueryRepo.ReadFirstInstalledItem(
		readFirstInstalledRequestDto,
	)
	if err != nil {
		return err
	}

	if serviceEntity.Status.String() != "stopped" {
		err = repo.Stop(serviceEntity.Name)
		if err != nil {
			return err
		}
	}

	serviceNameStr := serviceEntity.Name.String()
	removeOutput, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               SupervisorCtlBin + " remove " + serviceNameStr,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		combinedOutput := removeOutput + " " + err.Error()
		return errors.New("SupervisorRemoveError: " + combinedOutput)
	}

	err = repo.persistentDbSvc.Handler.
		Where("name = ?", serviceNameStr).
		Delete(dbModel.InstalledService{}).Error
	if err != nil {
		return err
	}

	err = repo.updateProcessManagerConf()
	if err != nil {
		return err
	}

	if serviceEntity.Nature == valueObject.ServiceNatureCustom {
		return nil
	}

	if serviceEntity.Nature == valueObject.ServiceNatureMulti {
		nameWithoutHashStr := strings.Split(serviceNameStr, "_")[0]
		nameWithoutHash, err := valueObject.NewServiceName(nameWithoutHashStr)
		if err != nil {
			return errors.New("CreateCustomServiceMultiNameError: " + err.Error())
		}
		name = nameWithoutHash
	}

	readInstallableDto := dto.ReadInstallableServicesItemsRequest{
		ServiceName: &name,
	}
	installableEntity, err := repo.servicesQueryRepo.ReadFirstInstallableItem(
		readInstallableDto,
	)
	if err != nil {
		return errors.New("ReadInstallableEntityError: " + err.Error())
	}

	err = repo.runCmdSteps(
		"Uninstall", installableEntity.UninstallCmdSteps,
		installableEntity.UninstallTimeoutSecs,
	)
	if err != nil {
		return err
	}

	for fileIndex, uninstallFilePath := range installableEntity.UninstallFilePaths {
		filePathStr := uninstallFilePath.String()

		slog.Debug("RemovingFilePath", slog.String("filePath", filePathStr))

		_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command: "rm",
			Args:    []string{"-rf", filePathStr},
		})
		if err != nil {
			fileIndexStr := strconv.Itoa(fileIndex)
			return errors.New(
				"RemoveFilePathError (" + fileIndexStr + "): " + err.Error(),
			)
		}
	}

	return nil
}

func (repo *ServicesCmdRepo) RefreshInstallableItems() error {
	_, err := os.Stat(infraEnvs.InstallableServicesItemsDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		repoCloneCmd := fmt.Sprintf(
			"cd %s; git clone --single-branch --branch %s %s services",
			infraEnvs.InfiniteOsMainDir, infraEnvs.InstallableServicesItemsRepoBranch,
			infraEnvs.InstallableServicesItemsRepoUrl,
		)
		_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command:               repoCloneCmd,
			ShouldRunWithSubShell: true,
		})
		if err != nil {
			return errors.New("CloneServicesItemsRepoError: " + err.Error())
		}
	}

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "cd " + infraEnvs.InstallableServicesItemsDir + ";" +
			"git clean -f -d; git reset --hard HEAD; git pull",
		ShouldRunWithSubShell: true,
	})
	return err
}
