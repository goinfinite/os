package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateInstallableService struct {
	Name                        valueObject.ServiceName     `json:"name"`
	Envs                        []valueObject.ServiceEnv    `json:"envs"`
	PortBindings                []valueObject.PortBinding   `json:"portBindings"`
	Version                     *valueObject.ServiceVersion `json:"version"`
	StartupFile                 *valueObject.UnixFilePath   `json:"startupFile"`
	AutoStart                   *bool                       `json:"autoStart"`
	TimeoutStartSecs            *uint                       `json:"timeoutStartSecs"`
	AutoRestart                 *bool                       `json:"autoRestart"`
	MaxStartRetries             *uint                       `json:"maxStartRetries"`
	AutoCreateMapping           *bool                       `json:"autoCreateMapping"`
	ScheduledInstallTimeoutSecs *uint                       `json:"scheduledInstallTimeoutSecs"`
	OperatorAccountId           valueObject.AccountId       `json:"-"`
	OperatorIpAddress           valueObject.IpAddress       `json:"-"`
}

func NewCreateInstallableService(
	name valueObject.ServiceName,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	version *valueObject.ServiceVersion,
	startupFile *valueObject.UnixFilePath,
	autoStart *bool,
	timeoutStartSecs *uint,
	autoRestart *bool,
	maxStartRetries *uint,
	autoCreateMapping *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) CreateInstallableService {
	return CreateInstallableService{
		Name:              name,
		Envs:              envs,
		PortBindings:      portBindings,
		Version:           version,
		StartupFile:       startupFile,
		AutoStart:         autoStart,
		TimeoutStartSecs:  timeoutStartSecs,
		AutoRestart:       autoRestart,
		MaxStartRetries:   maxStartRetries,
		AutoCreateMapping: autoCreateMapping,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
