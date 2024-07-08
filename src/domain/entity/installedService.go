package entity

import "github.com/speedianet/os/src/domain/valueObject"

type InstalledService struct {
	Name             valueObject.ServiceName    `json:"name"`
	Nature           valueObject.ServiceNature  `json:"nature"`
	Type             valueObject.ServiceType    `json:"type"`
	Version          valueObject.ServiceVersion `json:"version"`
	StartCmd         valueObject.UnixCommand    `json:"startCmd"`
	Status           valueObject.ServiceStatus  `json:"status"`
	Envs             []valueObject.ServiceEnv   `json:"envs"`
	PortBindings     []valueObject.PortBinding  `json:"portBindings"`
	StartupFile      *valueObject.UnixFilePath  `json:"startupFile"`
	AutoStart        *bool                      `json:"autoStart"`
	TimeoutStartSecs *uint                      `json:"timeoutStartSecs"`
	AutoRestart      *bool                      `json:"autoRestart"`
	MaxStartRetries  *uint                      `json:"maxStartRetries"`
	CreatedAt        valueObject.UnixTime       `json:"createdAt"`
	UpdatedAt        valueObject.UnixTime       `json:"updatedAt"`
}

func NewInstalledService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	version valueObject.ServiceVersion,
	startCmd valueObject.UnixCommand,
	status valueObject.ServiceStatus,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	startupFile *valueObject.UnixFilePath,
	autoStart *bool,
	timeoutStartSecs *uint,
	autoRestart *bool,
	maxStartRetries *uint,
	createdAt valueObject.UnixTime,
	updatedAt valueObject.UnixTime,
) InstalledService {
	return InstalledService{
		Name:             name,
		Nature:           nature,
		Type:             serviceType,
		Version:          version,
		StartCmd:         startCmd,
		Status:           status,
		Envs:             envs,
		PortBindings:     portBindings,
		StartupFile:      startupFile,
		AutoStart:        autoStart,
		TimeoutStartSecs: timeoutStartSecs,
		AutoRestart:      autoRestart,
		MaxStartRetries:  maxStartRetries,
		CreatedAt:        createdAt,
		UpdatedAt:        updatedAt,
	}
}
