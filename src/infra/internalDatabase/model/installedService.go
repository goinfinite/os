package dbModel

import (
	"log/slog"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type InstalledService struct {
	Name                 string `gorm:"primarykey;not null"`
	Nature               string `gorm:"not null"`
	Type                 string `gorm:"not null"`
	Version              string `gorm:"not null"`
	StartCmd             string `gorm:"not null"`
	Envs                 *string
	PortBindings         *string
	StopTimeoutSecs      int64
	StopCmdSteps         *string
	PreStartTimeoutSecs  int64
	PreStartCmdSteps     *string
	PostStartTimeoutSecs int64
	PostStartCmdSteps    *string
	PreStopTimeoutSecs   int64
	PreStopCmdSteps      *string
	PostStopTimeoutSecs  int64
	PostStopCmdSteps     *string
	ExecUser             *string
	WorkingDirectory     *string
	StartupFile          *string
	AutoStart            *bool
	AutoRestart          *bool
	TimeoutStartSecs     *uint
	MaxStartRetries      *uint
	LogOutputPath        *string
	LogErrorPath         *string
	AvatarUrl            *string
	CreatedAt            time.Time `gorm:"not null"`
	UpdatedAt            time.Time `gorm:"not null"`
}

func (InstalledService) TableName() string {
	return "installed_services"
}

func (InstalledService) InitialEntries() (entries []interface{}, err error) {
	osApiAvatarUrl := "https://goinfinite.github.io/os-services/system/os-api/assets/avatar.jpg"
	osWorkingDirectory := infraEnvs.InfiniteOsMainDir
	osLogOutputPath := "/dev/stdout"
	osLogErrorPath := "/dev/stderr"
	osApiPortBindings := infraEnvs.InfiniteOsApiHttpPublicPort + "/http"
	osApiService := InstalledService{
		Name:             "os-api",
		Nature:           valueObject.ServiceNatureSolo.String(),
		Type:             valueObject.ServiceTypeSystem.String(),
		Version:          infraEnvs.InfiniteOsVersion,
		StartCmd:         infraEnvs.InfiniteOsBinary + " serve",
		PortBindings:     &osApiPortBindings,
		WorkingDirectory: &osWorkingDirectory,
		LogOutputPath:    &osLogOutputPath,
		LogErrorPath:     &osLogErrorPath,
		AvatarUrl:        &osApiAvatarUrl,
	}

	cronAvatarUrl := "https://goinfinite.github.io/os-services/system/cron/assets/avatar.jpg"
	cronService := InstalledService{
		Name:      "cron",
		Nature:    valueObject.ServiceNatureSolo.String(),
		Type:      valueObject.ServiceTypeSystem.String(),
		Version:   "3.0",
		StartCmd:  "/usr/sbin/cron -f",
		AvatarUrl: &cronAvatarUrl,
	}

	nginxAvatarUrl := "https://goinfinite.github.io/os-services/system/nginx/assets/avatar.jpg"
	nginxPortBindings := "80/http;443/https"
	nginxAutoStart := false
	nginxService := InstalledService{
		Name:         "nginx",
		Nature:       valueObject.ServiceNatureSolo.String(),
		Type:         valueObject.ServiceTypeSystem.String(),
		Version:      "1.26.3",
		StartCmd:     "/usr/sbin/nginx",
		PortBindings: &nginxPortBindings,
		AutoStart:    &nginxAutoStart,
		AvatarUrl:    &nginxAvatarUrl,
	}

	return []interface{}{osApiService, cronService, nginxService}, nil
}

func (InstalledService) JoinCmdSteps(cmdSteps []tkValueObject.UnixCommand) string {
	cmdStepsStr := ""
	for _, cmdStep := range cmdSteps {
		cmdStepsStr += cmdStep.String() + "\n"
	}
	return strings.TrimSuffix(cmdStepsStr, "\n")
}

func (InstalledService) SplitCmdSteps(cmdStepsStr string) []tkValueObject.UnixCommand {
	rawCmdStepsList := strings.Split(cmdStepsStr, "\n")
	cmdSteps := []tkValueObject.UnixCommand{}
	for stepIndex, rawCmdStep := range rawCmdStepsList {
		if len(rawCmdStep) == 0 {
			continue
		}

		cmdStep, err := tkValueObject.NewUnixCommand(rawCmdStep)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("stepIndex", stepIndex))
			continue
		}
		cmdSteps = append(cmdSteps, cmdStep)
	}
	return cmdSteps
}

func (InstalledService) JoinEnvs(envs []valueObject.ServiceEnv) string {
	envsStr := ""
	for _, env := range envs {
		envsStr += env.String() + ";"
	}
	return strings.TrimSuffix(envsStr, ";")
}

func (InstalledService) SplitEnvs(envsStr string) []valueObject.ServiceEnv {
	rawEnvsList := strings.Split(envsStr, ";")
	envs := []valueObject.ServiceEnv{}
	for envIndex, rawEnv := range rawEnvsList {
		if len(rawEnv) == 0 {
			continue
		}

		env, err := valueObject.NewServiceEnv(rawEnv)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("envIndex", envIndex))
			continue
		}
		envs = append(envs, env)
	}
	return envs
}

func (InstalledService) JoinPortBindings(portBindings []valueObject.PortBinding) string {
	portBindingsStr := ""
	for _, portBinding := range portBindings {
		portBindingsStr += portBinding.String() + ";"
	}
	return strings.TrimSuffix(portBindingsStr, ";")
}

func (InstalledService) SplitPortBindings(portBindingsStr string) []valueObject.PortBinding {
	rawPortBindingsList := strings.Split(portBindingsStr, ";")
	portBindings := []valueObject.PortBinding{}
	for portIndex, rawPortBinding := range rawPortBindingsList {
		if len(rawPortBinding) == 0 {
			continue
		}

		portBinding, err := valueObject.NewPortBinding(rawPortBinding)
		if err != nil {
			slog.Debug(err.Error(), slog.Int("portIndex", portIndex))
			continue
		}
		portBindings = append(portBindings, portBinding)
	}
	return portBindings
}

func NewInstalledService(
	name, nature, serviceType, version, startCmd string,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []tkValueObject.UnixCommand,
	execUser, workingDirectory, startupFile *string,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *string,
	avatarUrl *string,
) InstalledService {
	var envsPtr *string
	if len(envs) > 0 {
		envsStr := InstalledService{}.JoinEnvs(envs)
		envsPtr = &envsStr
	}

	var portBindingsPtr *string
	if len(portBindings) > 0 {
		portBindingsStr := InstalledService{}.JoinPortBindings(portBindings)
		portBindingsPtr = &portBindingsStr
	}

	var stopStepsPtr *string
	if len(stopSteps) > 0 {
		stopStepsStr := InstalledService{}.JoinCmdSteps(stopSteps)
		stopStepsPtr = &stopStepsStr
	}

	var preStartStepsPtr *string
	if len(preStartSteps) > 0 {
		preStartStepsStr := InstalledService{}.JoinCmdSteps(preStartSteps)
		preStartStepsPtr = &preStartStepsStr
	}

	var postStartStepsPtr *string
	if len(postStartSteps) > 0 {
		postStartStepsStr := InstalledService{}.JoinCmdSteps(postStartSteps)
		postStartStepsPtr = &postStartStepsStr
	}

	var preStopStepsPtr *string
	if len(preStopSteps) > 0 {
		preStopStepsStr := InstalledService{}.JoinCmdSteps(preStopSteps)
		preStopStepsPtr = &preStopStepsStr
	}

	var postStopStepsPtr *string
	if len(postStopSteps) > 0 {
		postStopStepsStr := InstalledService{}.JoinCmdSteps(postStopSteps)
		postStopStepsPtr = &postStopStepsStr
	}

	return InstalledService{
		Name:              name,
		Nature:            nature,
		Type:              serviceType,
		Version:           version,
		StartCmd:          startCmd,
		Envs:              envsPtr,
		PortBindings:      portBindingsPtr,
		StopCmdSteps:      stopStepsPtr,
		PreStartCmdSteps:  preStartStepsPtr,
		PostStartCmdSteps: postStartStepsPtr,
		PreStopCmdSteps:   preStopStepsPtr,
		PostStopCmdSteps:  postStopStepsPtr,
		ExecUser:          execUser,
		WorkingDirectory:  workingDirectory,
		StartupFile:       startupFile,
		AutoStart:         autoStart,
		AutoRestart:       autoRestart,
		TimeoutStartSecs:  timeoutStartSecs,
		MaxStartRetries:   maxStartRetries,
		LogOutputPath:     logOutputPath,
		LogErrorPath:      logErrorPath,
		AvatarUrl:         avatarUrl,
	}
}

func (model InstalledService) ToEntity() (serviceEntity entity.InstalledService, err error) {
	name, err := valueObject.NewServiceName(model.Name)
	if err != nil {
		return serviceEntity, err
	}

	nature, err := valueObject.NewServiceNature(model.Nature)
	if err != nil {
		return serviceEntity, err
	}

	serviceType, err := valueObject.NewServiceType(model.Type)
	if err != nil {
		return serviceEntity, err
	}

	version, err := valueObject.NewServiceVersion(model.Version)
	if err != nil {
		return serviceEntity, err
	}

	startCmd, err := tkValueObject.NewUnixCommand(model.StartCmd)
	if err != nil {
		return serviceEntity, err
	}

	status, _ := valueObject.NewServiceStatus("running")

	envs := []valueObject.ServiceEnv{}
	if model.Envs != nil {
		envs = model.SplitEnvs(*model.Envs)
	}

	portBindings := []valueObject.PortBinding{}
	if model.PortBindings != nil {
		portBindings = model.SplitPortBindings(*model.PortBindings)
	}

	stopTimeoutSecs, err := tkValueObject.NewUnixTime(model.StopTimeoutSecs)
	if err != nil {
		return serviceEntity, err
	}

	stopCmdSteps := []tkValueObject.UnixCommand{}
	if model.StopCmdSteps != nil {
		stopCmdSteps = model.SplitCmdSteps(*model.StopCmdSteps)
	}

	preStartTimeoutSecs, err := tkValueObject.NewUnixTime(model.PreStartTimeoutSecs)
	if err != nil {
		return serviceEntity, err
	}

	preStartCmdSteps := []tkValueObject.UnixCommand{}
	if model.PreStartCmdSteps != nil {
		preStartCmdSteps = model.SplitCmdSteps(*model.PreStartCmdSteps)
	}

	postStartTimeoutSecs, err := tkValueObject.NewUnixTime(model.PostStartTimeoutSecs)
	if err != nil {
		return serviceEntity, err
	}

	postStartCmdSteps := []tkValueObject.UnixCommand{}
	if model.PostStartCmdSteps != nil {
		postStartCmdSteps = model.SplitCmdSteps(*model.PostStartCmdSteps)
	}

	preStopTimeoutSecs, err := tkValueObject.NewUnixTime(model.PreStopTimeoutSecs)
	if err != nil {
		return serviceEntity, err
	}

	preStopCmdSteps := []tkValueObject.UnixCommand{}
	if model.PreStopCmdSteps != nil {
		preStopCmdSteps = model.SplitCmdSteps(*model.PreStopCmdSteps)
	}

	postStopTimeoutSecs, err := tkValueObject.NewUnixTime(model.PostStopTimeoutSecs)
	if err != nil {
		return serviceEntity, err
	}

	postStopCmdSteps := []tkValueObject.UnixCommand{}
	if model.PostStopCmdSteps != nil {
		postStopCmdSteps = model.SplitCmdSteps(*model.PostStopCmdSteps)
	}

	var execUserPtr *tkValueObject.UnixUsername
	if model.ExecUser != nil {
		execUser, err := tkValueObject.NewUnixUsername(*model.ExecUser)
		if err != nil {
			return serviceEntity, err
		}
		execUserPtr = &execUser
	}

	var workingDirectoryPtr *tkValueObject.UnixAbsoluteFilePath
	if model.WorkingDirectory != nil {
		workingDirectory, err := valueObject.NewUnixFilePath(*model.WorkingDirectory)
		if err != nil {
			return serviceEntity, err
		}
		workingDirectoryPtr = &workingDirectory
	}

	var startupFilePtr *tkValueObject.UnixAbsoluteFilePath
	if model.StartupFile != nil {
		startupFile, err := valueObject.NewUnixFilePath(*model.StartupFile)
		if err != nil {
			return serviceEntity, err
		}
		startupFilePtr = &startupFile
	}

	var autoStart *bool
	if model.AutoStart != nil {
		autoStart = model.AutoStart
	}

	var autoRestart *bool
	if model.AutoRestart != nil {
		autoRestart = model.AutoRestart
	}

	var timeoutStartSecs *uint
	if model.TimeoutStartSecs != nil {
		timeoutStartSecs = model.TimeoutStartSecs
	}

	var maxStartRetries *uint
	if model.MaxStartRetries != nil {
		maxStartRetries = model.MaxStartRetries
	}

	var logOutputPathPtr *tkValueObject.UnixAbsoluteFilePath
	if model.LogOutputPath != nil {
		logOutputPath, err := valueObject.NewUnixFilePath(*model.LogOutputPath)
		if err != nil {
			return serviceEntity, err
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *tkValueObject.UnixAbsoluteFilePath
	if model.LogErrorPath != nil {
		logErrorPath, err := valueObject.NewUnixFilePath(*model.LogErrorPath)
		if err != nil {
			return serviceEntity, err
		}
		logErrorPathPtr = &logErrorPath
	}

	var avatarUrlPtr *tkValueObject.Url
	if model.AvatarUrl != nil {
		avatarUrl, err := tkValueObject.NewUrl(*model.AvatarUrl)
		if err != nil {
			return serviceEntity, err
		}
		avatarUrlPtr = &avatarUrl
	}

	return entity.NewInstalledService(
		name, nature, serviceType, version, startCmd, status, envs,
		portBindings, stopTimeoutSecs, stopCmdSteps, preStartTimeoutSecs,
		preStartCmdSteps, postStartTimeoutSecs, postStartCmdSteps, preStopTimeoutSecs,
		preStopCmdSteps, postStopTimeoutSecs, postStopCmdSteps, execUserPtr,
		workingDirectoryPtr, startupFilePtr, autoStart, autoRestart, timeoutStartSecs,
		maxStartRetries, logOutputPathPtr, logErrorPathPtr, avatarUrlPtr,
		tkValueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		tkValueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
