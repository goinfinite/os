package dbModel

import (
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
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
	StartupFile       *string
	ExecUser          *string
	WorkingDirectory  *string
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
	osWorkingDirectory := "/speedia"
	osLogOutputPath := "/dev/stdout"
	osLogErrorPath := "/dev/stderr"
	osApiPortBindings := "1618/http"
	osApiService := InstalledService{
		Name:             "os-api",
		Nature:           "solo",
		Type:             "system",
		Version:          infraEnvs.SpeediaOsVersion,
		StartCmd:         "/speedia/os serve",
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

func NewInstalledService(
	name, nature, serviceType, version, startCmd string,
	envs []valueObject.ServiceEnv, portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []valueObject.UnixCommand,
	startupFile, execUser, workingDirectory *string, autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint, logOutputPath, logErrorPath *string,
) InstalledService {
	var envsPtr *string
	if len(envs) > 0 {
		envsStr := ""
		for _, env := range envs {
			envsStr += env.String() + ";"
		}
		envsStr = strings.TrimSuffix(envsStr, ";")
		envsPtr = &envsStr
	}

	var portBindingsPtr *string
	if len(portBindings) > 0 {
		portBindingsStr := ""
		for _, portBinding := range portBindings {
			portBindingsStr += portBinding.String() + ";"
		}
		portBindingsStr = strings.TrimSuffix(portBindingsStr, ";")
		portBindingsPtr = &portBindingsStr
	}

	var stopStepsPtr *string
	if len(stopSteps) > 0 {
		stopStepsStr := ""
		for _, stopStep := range stopSteps {
			stopStepsStr += stopStep.String() + "\n"
		}
		stopStepsStr = strings.TrimSuffix(stopStepsStr, "\n")
		stopStepsPtr = &stopStepsStr
	}

	var preStartStepsPtr *string
	if len(preStartSteps) > 0 {
		preStartStepsStr := ""
		for _, preStartStep := range preStartSteps {
			preStartStepsStr += preStartStep.String() + "\n"
		}
		preStartStepsStr = strings.TrimSuffix(preStartStepsStr, "\n")
		preStartStepsPtr = &preStartStepsStr
	}

	var postStartStepsPtr *string
	if len(postStartSteps) > 0 {
		postStartStepsStr := ""
		for _, postStartStep := range postStartSteps {
			postStartStepsStr += postStartStep.String() + "\n"
		}
		postStartStepsStr = strings.TrimSuffix(postStartStepsStr, "\n")
		postStartStepsPtr = &postStartStepsStr
	}

	var preStopStepsPtr *string
	if len(preStopSteps) > 0 {
		preStopStepsStr := ""
		for _, preStopStep := range preStopSteps {
			preStopStepsStr += preStopStep.String() + "\n"
		}
		preStopStepsStr = strings.TrimSuffix(preStopStepsStr, "\n")
		preStopStepsPtr = &preStopStepsStr
	}

	var postStopStepsPtr *string
	if len(postStopSteps) > 0 {
		postStopStepsStr := ""
		for _, postStopStep := range postStopSteps {
			postStopStepsStr += postStopStep.String() + "\n"
		}
		postStopStepsStr = strings.TrimSuffix(postStopStepsStr, "\n")
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
		StartupFile:       startupFile,
		ExecUser:          execUser,
		WorkingDirectory:  workingDirectory,
		AutoStart:         autoStart,
		AutoRestart:       autoRestart,
		TimeoutStartSecs:  timeoutStartSecs,
		MaxStartRetries:   maxStartRetries,
		LogOutputPath:     logOutputPath,
		LogErrorPath:      logErrorPath,
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

	startCmd, err := valueObject.NewUnixCommand(model.StartCmd)
	if err != nil {
		return serviceEntity, err
	}

	status, _ := valueObject.NewServiceStatus("running")

	var envs []valueObject.ServiceEnv
	if model.Envs != nil {
		rawEnvsList := strings.Split(*model.Envs, ";")
		for _, rawEnv := range rawEnvsList {
			env, err := valueObject.NewServiceEnv(rawEnv)
			if err != nil {
				return serviceEntity, err
			}
			envs = append(envs, env)
		}
	}

	var portBindings []valueObject.PortBinding
	if model.PortBindings != nil {
		rawPortBindingsList := strings.Split(*model.PortBindings, ";")
		for _, rawPortBinding := range rawPortBindingsList {
			portBinding, err := valueObject.NewPortBinding(rawPortBinding)
			if err != nil {
				return serviceEntity, err
			}
			portBindings = append(portBindings, portBinding)
		}
	}

	var stopCmdSteps []valueObject.UnixCommand
	if model.StopCmdSteps != nil {
		rawStopCmdStepsList := strings.Split(*model.StopCmdSteps, "\n")
		for _, rawStopCmdStep := range rawStopCmdStepsList {
			stopCmdStep, err := valueObject.NewUnixCommand(rawStopCmdStep)
			if err != nil {
				return serviceEntity, err
			}
			stopCmdSteps = append(stopCmdSteps, stopCmdStep)
		}
	}

	var preStartCmdSteps []valueObject.UnixCommand
	if model.PreStartCmdSteps != nil {
		rawPreStartCmdStepsList := strings.Split(*model.PreStartCmdSteps, "\n")
		for _, rawPreStartCmdStep := range rawPreStartCmdStepsList {
			preStartCmdStep, err := valueObject.NewUnixCommand(rawPreStartCmdStep)
			if err != nil {
				return serviceEntity, err
			}
			preStartCmdSteps = append(preStartCmdSteps, preStartCmdStep)
		}
	}

	var postStartCmdSteps []valueObject.UnixCommand
	if model.PostStartCmdSteps != nil {
		rawPostStartCmdStepsList := strings.Split(*model.PostStartCmdSteps, "\n")
		for _, rawPostStartCmdStep := range rawPostStartCmdStepsList {
			postStartCmdStep, err := valueObject.NewUnixCommand(rawPostStartCmdStep)
			if err != nil {
				return serviceEntity, err
			}
			postStartCmdSteps = append(postStartCmdSteps, postStartCmdStep)
		}
	}

	var preStopCmdSteps []valueObject.UnixCommand
	if model.PreStopCmdSteps != nil {
		rawPreStopCmdStepsList := strings.Split(*model.PreStopCmdSteps, "\n")
		for _, rawPreStopCmdStep := range rawPreStopCmdStepsList {
			preStopCmdStep, err := valueObject.NewUnixCommand(rawPreStopCmdStep)
			if err != nil {
				return serviceEntity, err
			}
			preStopCmdSteps = append(preStopCmdSteps, preStopCmdStep)
		}
	}

	var postStopCmdSteps []valueObject.UnixCommand
	if model.PostStopCmdSteps != nil {
		rawPostStopCmdStepsList := strings.Split(*model.PostStopCmdSteps, "\n")
		for _, rawPostStopCmdStep := range rawPostStopCmdStepsList {
			postStopCmdStep, err := valueObject.NewUnixCommand(rawPostStopCmdStep)
			if err != nil {
				return serviceEntity, err
			}
			postStopCmdSteps = append(postStopCmdSteps, postStopCmdStep)
		}
	}

	var startupFilePtr *valueObject.UnixFilePath
	if model.StartupFile != nil {
		startupFile, err := valueObject.NewUnixFilePath(*model.StartupFile)
		if err != nil {
			return serviceEntity, err
		}
		startupFilePtr = &startupFile
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
		name, nature, serviceType, version, startCmd, status, envs, portBindings,
		stopCmdSteps, preStartCmdSteps, postStartCmdSteps, preStopCmdSteps, postStopCmdSteps,
		startupFilePtr, execUserPtr, workingDirectoryPtr, autoStart, autoRestart,
		timeoutStartSecs, maxStartRetries, logOutputPathPtr, logErrorPathPtr,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
