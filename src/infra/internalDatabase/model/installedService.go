package dbModel

import (
	"log"
	"strings"
	"time"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
)

type InstalledService struct {
	Name              string `gorm:"primarykey;not null"`
	Nature            string `gorm:"not null"`
	Type              string `gorm:"not null"`
	Version           string `gorm:"not null"`
	StartCmd          string `gorm:"not null"`
	Envs              *string
	PortBindings      *string
	StopCmdSteps      *string
	PreStartCmdSteps  *string
	PostStartCmdSteps *string
	PreStopCmdSteps   *string
	PostStopCmdSteps  *string
	ExecUser          *string
	WorkingDirectory  *string
	StartupFile       *string
	AutoStart         *bool
	AutoRestart       *bool
	TimeoutStartSecs  *uint
	MaxStartRetries   *uint
	LogOutputPath     *string
	LogErrorPath      *string
	CreatedAt         time.Time `gorm:"not null"`
	UpdatedAt         time.Time `gorm:"not null"`
}

func (InstalledService) TableName() string {
	return "installed_services"
}

func (InstalledService) InitialEntries() (entries []interface{}, err error) {
	osWorkingDirectory := "/infinite"
	osLogOutputPath := "/dev/stdout"
	osLogErrorPath := "/dev/stderr"
	osApiPortBindings := "1618/http"
	osApiService := InstalledService{
		Name:             "os-api",
		Nature:           "solo",
		Type:             "system",
		Version:          infraEnvs.InfiniteOsVersion,
		StartCmd:         infraEnvs.InfiniteOsBinary + " serve",
		PortBindings:     &osApiPortBindings,
		WorkingDirectory: &osWorkingDirectory,
		LogOutputPath:    &osLogOutputPath,
		LogErrorPath:     &osLogErrorPath,
	}

	cronService := InstalledService{
		Name:     "cron",
		Nature:   "solo",
		Type:     "system",
		Version:  "3.0",
		StartCmd: "/usr/sbin/cron -f",
	}

	nginxPortBindings := "80/http;443/https"
	nginxAutoStart := false
	nginxService := InstalledService{
		Name:         "nginx",
		Nature:       "solo",
		Type:         "system",
		Version:      "1.24.0",
		StartCmd:     "/usr/sbin/nginx",
		PortBindings: &nginxPortBindings,
		AutoStart:    &nginxAutoStart,
	}

	return []interface{}{osApiService, cronService, nginxService}, nil
}

func (InstalledService) JoinCmdSteps(cmdSteps []valueObject.UnixCommand) string {
	cmdStepsStr := ""
	for _, cmdStep := range cmdSteps {
		cmdStepsStr += cmdStep.String() + "\n"
	}
	return strings.TrimSuffix(cmdStepsStr, "\n")
}

func (InstalledService) SplitCmdSteps(cmdStepsStr string) []valueObject.UnixCommand {
	rawCmdStepsList := strings.Split(cmdStepsStr, "\n")
	var cmdSteps []valueObject.UnixCommand
	for stepIndex, rawCmdStep := range rawCmdStepsList {
		cmdStep, err := valueObject.NewUnixCommand(rawCmdStep)
		if err != nil {
			log.Printf("[index %d] %s", stepIndex, err)
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
	var envs []valueObject.ServiceEnv
	for envIndex, rawEnv := range rawEnvsList {
		env, err := valueObject.NewServiceEnv(rawEnv)
		if err != nil {
			log.Printf("[index %d] %s", envIndex, err)
			continue
		}
		envs = append(envs, env)
	}
	return envs
}

func (InstalledService) JoinPortBindings(
	portBindings []valueObject.PortBinding,
) string {
	portBindingsStr := ""
	for _, portBinding := range portBindings {
		portBindingsStr += portBinding.String() + ";"
	}
	return strings.TrimSuffix(portBindingsStr, ";")
}

func (InstalledService) SplitPortBindings(
	portBindingsStr string,
) []valueObject.PortBinding {
	rawPortBindingsList := strings.Split(portBindingsStr, ";")
	var portBindings []valueObject.PortBinding
	for portIndex, rawPortBinding := range rawPortBindingsList {
		portBinding, err := valueObject.NewPortBinding(rawPortBinding)
		if err != nil {
			log.Printf("[index %d] %s", portIndex, err)
			continue
		}
		portBindings = append(portBindings, portBinding)
	}
	return portBindings
}

func NewInstalledService(
	name, nature, serviceType, version, startCmd string,
	envs []valueObject.ServiceEnv, portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []valueObject.UnixCommand,
	execUser, workingDirectory, startupFile *string, autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint, logOutputPath, logErrorPath *string,
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
	}
}

func (model InstalledService) ToEntity() (
	serviceEntity entity.InstalledService, err error,
) {
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

	startCmd, err := valueObject.NewUnixCommand(model.StartCmd)
	if err != nil {
		return serviceEntity, err
	}

	status, _ := valueObject.NewServiceStatus("running")

	var envs []valueObject.ServiceEnv
	if model.Envs != nil {
		envs = model.SplitEnvs(*model.Envs)
	}

	var portBindings []valueObject.PortBinding
	if model.PortBindings != nil {
		portBindings = model.SplitPortBindings(*model.PortBindings)
	}

	var stopCmdSteps []valueObject.UnixCommand
	if model.StopCmdSteps != nil {
		stopCmdSteps = model.SplitCmdSteps(*model.StopCmdSteps)
	}

	var preStartCmdSteps []valueObject.UnixCommand
	if model.PreStartCmdSteps != nil {
		preStartCmdSteps = model.SplitCmdSteps(*model.PreStartCmdSteps)
	}

	var postStartCmdSteps []valueObject.UnixCommand
	if model.PostStartCmdSteps != nil {
		postStartCmdSteps = model.SplitCmdSteps(*model.PostStartCmdSteps)
	}

	var preStopCmdSteps []valueObject.UnixCommand
	if model.PreStopCmdSteps != nil {
		preStopCmdSteps = model.SplitCmdSteps(*model.PreStopCmdSteps)
	}

	var postStopCmdSteps []valueObject.UnixCommand
	if model.PostStopCmdSteps != nil {
		postStopCmdSteps = model.SplitCmdSteps(*model.PostStopCmdSteps)
	}

	var execUserPtr *valueObject.UnixUsername
	if model.ExecUser != nil {
		execUser, err := valueObject.NewUnixUsername(*model.ExecUser)
		if err != nil {
			return serviceEntity, err
		}
		execUserPtr = &execUser
	}

	var workingDirectoryPtr *valueObject.UnixFilePath
	if model.WorkingDirectory != nil {
		workingDirectory, err := valueObject.NewUnixFilePath(*model.WorkingDirectory)
		if err != nil {
			return serviceEntity, err
		}
		workingDirectoryPtr = &workingDirectory
	}

	var startupFilePtr *valueObject.UnixFilePath
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

	var logOutputPathPtr *valueObject.UnixFilePath
	if model.LogOutputPath != nil {
		logOutputPath, err := valueObject.NewUnixFilePath(*model.LogOutputPath)
		if err != nil {
			return serviceEntity, err
		}
		logOutputPathPtr = &logOutputPath
	}

	var logErrorPathPtr *valueObject.UnixFilePath
	if model.LogErrorPath != nil {
		logErrorPath, err := valueObject.NewUnixFilePath(*model.LogErrorPath)
		if err != nil {
			return serviceEntity, err
		}
		logErrorPathPtr = &logErrorPath
	}

	return entity.NewInstalledService(
		name, nature, serviceType, version, startCmd, status, envs, portBindings, nil,
		stopCmdSteps, preStartCmdSteps, postStartCmdSteps, preStopCmdSteps,
		postStopCmdSteps, execUserPtr, workingDirectoryPtr, startupFilePtr, autoStart,
		autoRestart, timeoutStartSecs, maxStartRetries, logOutputPathPtr,
		logErrorPathPtr, valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
