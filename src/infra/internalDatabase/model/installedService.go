package dbModel

import (
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
)

type InstalledService struct {
	Name             string `gorm:"primarykey;not null"`
	Nature           string `gorm:"not null"`
	Type             string `gorm:"not null"`
	Version          string `gorm:"not null"`
	StartCmd         string `gorm:"not null"`
	Envs             *string
	PortBindings     *string
	StartupFile      *string
	AutoStart        *bool
	TimeoutStartSecs *uint
	AutoRestart      *bool
	MaxStartRetries  *uint
	CreatedAt        time.Time `gorm:"not null"`
	UpdatedAt        time.Time `gorm:"not null"`
}

func (InstalledService) TableName() string {
	return "installed_services"
}

func (InstalledService) InitialEntries() (entries []interface{}, err error) {
	autoStart := true
	timeoutStartSecs := uint(10)
	autoRestart := true
	maxStartRetries := uint(3)

	osApiPortBindings := "1618/http"
	osApiInstalledService := InstalledService{
		Name:             "os-api",
		Nature:           "solo",
		Type:             "system",
		Version:          infraEnvs.SpeediaOsVersion,
		StartCmd:         "/speedia/os serve",
		Envs:             nil,
		PortBindings:     &osApiPortBindings,
		StartupFile:      nil,
		AutoStart:        &autoStart,
		TimeoutStartSecs: &timeoutStartSecs,
		AutoRestart:      &autoRestart,
		MaxStartRetries:  &maxStartRetries,
	}

	cronInstalledService := InstalledService{
		Name:             "cron",
		Nature:           "solo",
		Type:             "system",
		Version:          "3.0",
		StartCmd:         "/usr/sbin/cron -f",
		Envs:             nil,
		PortBindings:     nil,
		StartupFile:      nil,
		AutoStart:        &autoStart,
		TimeoutStartSecs: &timeoutStartSecs,
		AutoRestart:      &autoRestart,
		MaxStartRetries:  &maxStartRetries,
	}

	nginxPortBindings := "80/http;443/https"
	nginxInstalledService := InstalledService{
		Name:             "nginx",
		Nature:           "solo",
		Type:             "system",
		Version:          "1.24.0",
		StartCmd:         "/usr/sbin/nginx",
		Envs:             nil,
		PortBindings:     &nginxPortBindings,
		StartupFile:      nil,
		AutoStart:        &autoStart,
		TimeoutStartSecs: &timeoutStartSecs,
		AutoRestart:      &autoRestart,
		MaxStartRetries:  &maxStartRetries,
	}

	return []interface{}{osApiInstalledService, cronInstalledService, nginxInstalledService}, nil
}

func NewInstalledService(
	name, nature, serviceType, version, startCmd string,
	envs []valueObject.ServiceEnv, portBindings []valueObject.PortBinding, startupFile *string,
	autoStart *bool, timeoutStartSecs *uint, autoRestart *bool, maxStartRetries *uint,
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

	return InstalledService{
		Name:             name,
		Nature:           nature,
		Type:             serviceType,
		Version:          version,
		StartCmd:         startCmd,
		Envs:             envsPtr,
		PortBindings:     portBindingsPtr,
		StartupFile:      startupFile,
		AutoStart:        autoStart,
		TimeoutStartSecs: timeoutStartSecs,
		AutoRestart:      autoRestart,
		MaxStartRetries:  maxStartRetries,
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

	var timeoutStartSecs *uint
	if model.TimeoutStartSecs != nil {
		timeoutStartSecs = model.TimeoutStartSecs
	}

	var autoRestart *bool
	if model.AutoRestart != nil {
		autoRestart = model.AutoRestart
	}

	var maxStartRetries *uint
	if model.MaxStartRetries != nil {
		maxStartRetries = model.MaxStartRetries
	}

	return entity.NewInstalledService(
		name,
		nature,
		serviceType,
		version,
		startCmd,
		status,
		envs,
		portBindings,
		startupFilePtr,
		autoStart,
		timeoutStartSecs,
		autoRestart,
		maxStartRetries,
		valueObject.NewUnixTimeWithGoTime(model.CreatedAt),
		valueObject.NewUnixTimeWithGoTime(model.UpdatedAt),
	), nil
}
