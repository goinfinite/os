package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type CreateCustomService struct {
	Name                           valueObject.ServiceName     `json:"name"`
	Type                           valueObject.ServiceType     `json:"type"`
	StartCmd                       valueObject.UnixCommand     `json:"startCmd"`
	Envs                           []valueObject.ServiceEnv    `json:"envs"`
	PortBindings                   []valueObject.PortBinding   `json:"portBindings"`
	StopCmdSteps                   []valueObject.UnixCommand   `json:"stopCmdSteps"`
	PreStartCmdSteps               []valueObject.UnixCommand   `json:"preStartCmdSteps"`
	PostStartCmdSteps              []valueObject.UnixCommand   `json:"postStartCmdSteps"`
	PreStopCmdSteps                []valueObject.UnixCommand   `json:"preStopCmdSteps"`
	PostStopCmdSteps               []valueObject.UnixCommand   `json:"postStopCmdSteps"`
	Version                        *valueObject.ServiceVersion `json:"version"`
	ExecUser                       *valueObject.UnixUsername   `json:"execUser"`
	WorkingDirectory               *valueObject.UnixFilePath   `json:"workingDirectory"`
	AutoStart                      *bool                       `json:"autoStart"`
	AutoRestart                    *bool                       `json:"autoRestart"`
	TimeoutStartSecs               *uint                       `json:"timeoutStartSecs"`
	MaxStartRetries                *uint                       `json:"maxStartRetries"`
	LogOutputPath                  *valueObject.UnixFilePath   `json:"logOutputPath"`
	LogErrorPath                   *valueObject.UnixFilePath   `json:"logErrorPath"`
	AvatarUrl                      *valueObject.Url            `json:"avatarUrl"`
	AutoCreateMapping              *bool                       `json:"autoCreateMapping"`
	MappingHostname                *valueObject.Fqdn           `json:"mappingHostname"`
	MappingPath                    *valueObject.MappingPath    `json:"mappingPath"`
	MappingUpgradeInsecureRequests *bool                       `json:"mappingUpgradeInsecureRequests"`
	OperatorAccountId              valueObject.AccountId       `json:"-"`
	OperatorIpAddress              valueObject.IpAddress       `json:"-"`
}

func NewCreateCustomService(
	name valueObject.ServiceName,
	serviceType valueObject.ServiceType,
	startCmd valueObject.UnixCommand,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopSteps, preStartSteps, postStartSteps, preStopSteps, postStopSteps []valueObject.UnixCommand,
	version *valueObject.ServiceVersion,
	execUser *valueObject.UnixUsername,
	workingDirectory *valueObject.UnixFilePath,
	autoStart, autoRestart *bool,
	timeoutStartSecs, maxStartRetries *uint,
	logOutputPath, logErrorPath *valueObject.UnixFilePath,
	avatarUrl *valueObject.Url,
	autoCreateMapping *bool,
	mappingHostname *valueObject.Fqdn,
	mappingPath *valueObject.MappingPath,
	mappingUpgradeInsecureRequests *bool,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
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
