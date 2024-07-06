package dbModel

import (
	"strings"
	"time"

	"github.com/speedianet/os/src/domain/valueObject"
	infraEnvs "github.com/speedianet/os/src/infra/envs"
)

type InstalledService struct {
	Name            string `gorm:"primarykey;not null"`
	Nature          string `gorm:"not null"`
	Type            string `gorm:"not null"`
	Version         string `gorm:"not null"`
	Command         string `gorm:"not null"`
	AutoStart       bool   `gorm:"not null"`
	TimeoutStartSec uint   `gorm:"not null"`
	AutoRestart     bool   `gorm:"not null"`
	MaxStartRetries uint   `gorm:"not null"`
	StartupFile     *string
	Envs            *string
	PortBindings    *string
	CreatedAt       time.Time `gorm:"not null"`
	UpdatedAt       time.Time `gorm:"not null"`
}

func (InstalledService) TableName() string {
	return "installed_services"
}

func (InstalledService) InitialEntries() (entries []interface{}, err error) {
	osApiPortBindings := "1618/http"
	osApiInstalledService := InstalledService{
		Name:            "os-api",
		Nature:          "solo",
		Type:            "system",
		Version:         infraEnvs.SpeediaOsVersion,
		Command:         "/speedia/os serve",
		AutoStart:       true,
		TimeoutStartSec: 10,
		AutoRestart:     true,
		MaxStartRetries: 3,
		StartupFile:     nil,
		PortBindings:    &osApiPortBindings,
	}

	cronInstalledService := InstalledService{
		Name:            "cron",
		Nature:          "solo",
		Type:            "system",
		Version:         "3.0",
		Command:         "/usr/sbin/cron -f",
		AutoStart:       true,
		TimeoutStartSec: 10,
		AutoRestart:     true,
		MaxStartRetries: 3,
		StartupFile:     nil,
		PortBindings:    nil,
	}

	nginxPortBindings := "80/http;443/https"
	nginxInstalledService := InstalledService{
		Name:            "nginx",
		Nature:          "solo",
		Type:            "system",
		Version:         "1.24.0",
		Command:         "/usr/sbin/nginx",
		AutoStart:       true,
		TimeoutStartSec: 10,
		AutoRestart:     true,
		MaxStartRetries: 3,
		StartupFile:     nil,
		PortBindings:    &nginxPortBindings,
	}

	return []interface{}{osApiInstalledService, cronInstalledService, nginxInstalledService}, nil
}

func NewInstalledService(
	name, nature, serviceType, version, command string,
	autoStart bool,
	timeoutStartSec uint,
	autoRestart bool,
	maxStartRetries uint,
	startupFile *string,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
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
		Name:            name,
		Nature:          nature,
		Type:            serviceType,
		Version:         version,
		Command:         command,
		AutoStart:       autoStart,
		TimeoutStartSec: timeoutStartSec,
		AutoRestart:     autoRestart,
		MaxStartRetries: maxStartRetries,
		StartupFile:     startupFile,
		Envs:            envsPtr,
		PortBindings:    portBindingsPtr,
	}
}
