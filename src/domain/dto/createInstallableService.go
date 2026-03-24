package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateInstallableService struct {
	Name                           valueObject.ServiceName              `json:"name"`
	Envs                           []valueObject.ServiceEnv             `json:"envs"`
	PortBindings                   []valueObject.PortBinding            `json:"portBindings"`
	Version                        *valueObject.ServiceVersion          `json:"version"`
	StartupFile                    *tkValueObject.UnixAbsoluteFilePath  `json:"startupFile"`
	WorkingDir                     *tkValueObject.UnixAbsoluteFilePath  `json:"workingDir"`
	AutoStart                      *bool                                `json:"autoStart"`
	TimeoutStartSecs               *uint                                `json:"timeoutStartSecs"`
	AutoRestart                    *bool                                `json:"autoRestart"`
	MaxStartRetries                *uint                                `json:"maxStartRetries"`
	AutoCreateMapping              *bool                                `json:"autoCreateMapping"`
	MappingHostname                *tkValueObject.Fqdn                  `json:"mappingHostname"`
	MappingPath                    *valueObject.MappingPath             `json:"mappingPath"`
	MappingUpgradeInsecureRequests *bool                                `json:"mappingUpgradeInsecureRequests"`
	OperatorAccountId              tkValueObject.AccountId              `json:"-"`
	OperatorIpAddress              tkValueObject.IpAddress              `json:"-"`
}

func NewCreateInstallableService(
	name valueObject.ServiceName,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	version *valueObject.ServiceVersion,
	startupFile *tkValueObject.UnixAbsoluteFilePath,
	workingDir *tkValueObject.UnixAbsoluteFilePath,
	autoStart *bool,
	timeoutStartSecs *uint,
	autoRestart *bool,
	maxStartRetries *uint,
	autoCreateMapping *bool,
	mappingHostname *tkValueObject.Fqdn,
	mappingPath *valueObject.MappingPath,
	mappingUpgradeInsecureRequests *bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateInstallableService {
	return CreateInstallableService{
		Name:                           name,
		Envs:                           envs,
		PortBindings:                   portBindings,
		Version:                        version,
		StartupFile:                    startupFile,
		WorkingDir:                     workingDir,
		AutoStart:                      autoStart,
		TimeoutStartSecs:               timeoutStartSecs,
		AutoRestart:                    autoRestart,
		MaxStartRetries:                maxStartRetries,
		AutoCreateMapping:              autoCreateMapping,
		MappingHostname:                mappingHostname,
		MappingPath:                    mappingPath,
		MappingUpgradeInsecureRequests: mappingUpgradeInsecureRequests,
		OperatorAccountId:              operatorAccountId,
		OperatorIpAddress:              operatorIpAddress,
	}
}
