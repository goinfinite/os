package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type CreateCustomService struct {
	Name                           valueObject.ServiceName     `json:"name"`
	Type                           valueObject.ServiceType     `json:"type"`
	StartCmd                       tkValueObject.UnixCommand   `json:"startCmd"`
	Envs                           []valueObject.ServiceEnv    `json:"envs"`
	PortBindings                   []valueObject.PortBinding   `json:"portBindings"`
	StopCmdSteps                   []tkValueObject.UnixCommand `json:"stopCmdSteps"`
	PreStartCmdSteps               []tkValueObject.UnixCommand `json:"preStartCmdSteps"`
	PostStartCmdSteps              []tkValueObject.UnixCommand `json:"postStartCmdSteps"`
	PreStopCmdSteps                []tkValueObject.UnixCommand `json:"preStopCmdSteps"`
	PostStopCmdSteps               []tkValueObject.UnixCommand `json:"postStopCmdSteps"`
	Version                        *valueObject.ServiceVersion `json:"version"`
	ExecUser                       *tkValueObject.UnixUsername                `json:"execUser"`
	WorkingDirectory               *tkValueObject.UnixAbsoluteFilePath       `json:"workingDirectory"`
	AutoStart                      *bool                       `json:"autoStart"`
	AutoRestart                    *bool                       `json:"autoRestart"`
	TimeoutStartSecs               *uint                       `json:"timeoutStartSecs"`
	MaxStartRetries                *uint                       `json:"maxStartRetries"`
	LogOutputPath                  *tkValueObject.UnixAbsoluteFilePath       `json:"logOutputPath"`
	LogErrorPath                   *tkValueObject.UnixAbsoluteFilePath       `json:"logErrorPath"`
	AvatarUrl                      *tkValueObject.Url          `json:"avatarUrl"`
	AutoCreateMapping              *bool                       `json:"autoCreateMapping"`
	MappingHostname                *tkValueObject.Fqdn         `json:"mappingHostname"`
	MappingPath                    *valueObject.MappingPath    `json:"mappingPath"`
	MappingUpgradeInsecureRequests *bool                       `json:"mappingUpgradeInsecureRequests"`
	OperatorAccountId              tkValueObject.AccountId     `json:"-"`
	OperatorIpAddress              tkValueObject.IpAddress     `json:"-"`
}

func NewCreateCustomService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	startCmd tkValueObject.UnixCommand,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []tkValueObject.UnixCommand,
	version *valueObject.ServiceVersion,
	execUser *tkValueObject.UnixUsername,
	workingDirectory *tkValueObject.UnixAbsoluteFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *tkValueObject.UnixAbsoluteFilePath,
	avatarUrl *tkValueObject.Url,
	autoCreateMapping *bool,
	mappingHostname *tkValueObject.Fqdn,
	mappingPath *valueObject.MappingPath,
	mappingUpgradeInsecureRequests *bool,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) CreateCustomService {
	return CreateCustomService{
		Name:                           name,
		Type:                           serviceType,
		StartCmd:                       startCmd,
		Envs:                           envs,
		PortBindings:                   portBindings,
		StopCmdSteps:                   stopSteps,
		PreStartCmdSteps:               preStartSteps,
		PostStartCmdSteps:              postStartSteps,
		PreStopCmdSteps:                preStopSteps,
		PostStopCmdSteps:               postStopSteps,
		Version:                        version,
		ExecUser:                       execUser,
		WorkingDirectory:               workingDirectory,
		AutoStart:                      autoStart,
		TimeoutStartSecs:               timeoutStartSecs,
		AutoRestart:                    autoRestart,
		MaxStartRetries:                maxStartRetries,
		LogOutputPath:                  logOutputPath,
		LogErrorPath:                   logErrorPath,
		AvatarUrl:                      avatarUrl,
		AutoCreateMapping:              autoCreateMapping,
		MappingHostname:                mappingHostname,
		MappingPath:                    mappingPath,
		MappingUpgradeInsecureRequests: mappingUpgradeInsecureRequests,
		OperatorAccountId:              operatorAccountId,
		OperatorIpAddress:              operatorIpAddress,
	}
}
