package dbModel

import (
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
)

type InstalledService struct {
	Name             string `gorm:"primarykey;not null"`
	Nature           string `gorm:"not null"`
	Type             string `gorm:"not null"`
	Version          string `gorm:"not null"`
	Command          string `gorm:"not null"`
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
		Command:          "/speedia/os serve",
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
		Command:          "/usr/sbin/cron -f",
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
		Command:          "/usr/sbin/nginx",
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
	name, nature, serviceType, version, command string,
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
		Command:          command,
		Envs:             envsPtr,
		PortBindings:     portBindingsPtr,
		StartupFile:      startupFile,
		AutoStart:        autoStart,
		TimeoutStartSecs: timeoutStartSecs,
		AutoRestart:      autoRestart,
		MaxStartRetries:  maxStartRetries,
	}
}
