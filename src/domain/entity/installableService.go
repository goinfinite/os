package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type InstallableService struct {
	ManifestVersion      valueObject.ServiceManifestVersion `json:"manifestVersion"`
	Name                 valueObject.ServiceName            `json:"name"`
	Nature               valueObject.ServiceNature          `json:"nature"`
	Type                 valueObject.ServiceType            `json:"type"`
	StartCmd             valueObject.UnixCommand            `json:"startCmd"`
	Description          valueObject.ServiceDescription     `json:"description"`
	Versions             []valueObject.ServiceVersion       `json:"versions"`
	Envs                 []valueObject.ServiceEnv           `json:"envs"`
	PortBindings         []valueObject.PortBinding          `json:"portBindings"`
	StopTimeoutSecs      valueObject.UnixTime               `json:"-"`
	StopCmdSteps         []valueObject.UnixCommand          `json:"-"`
	InstallTimeoutSecs   valueObject.UnixTime               `json:"-"`
	InstallCmdSteps      []valueObject.UnixCommand          `json:"-"`
	UninstallTimeoutSecs valueObject.UnixTime               `json:"-"`
	UninstallCmdSteps    []valueObject.UnixCommand          `json:"-"`
	UninstallFilePaths   []valueObject.UnixFilePath         `json:"-"`
	PreStartTimeoutSecs  valueObject.UnixTime               `json:"-"`
	PreStartCmdSteps     []valueObject.UnixCommand          `json:"-"`
	PostStartTimeoutSecs valueObject.UnixTime               `json:"-"`
	PostStartCmdSteps    []valueObject.UnixCommand          `json:"-"`
	PreStopTimeoutSecs   valueObject.UnixTime               `json:"-"`
	PreStopCmdSteps      []valueObject.UnixCommand          `json:"-"`
	PostStopTimeoutSecs  valueObject.UnixTime               `json:"-"`
	PostStopCmdSteps     []valueObject.UnixCommand          `json:"-"`
	ExecUser             *valueObject.UnixUsername          `json:"execUser"`
	WorkingDirectory     *valueObject.UnixFilePath          `json:"workingDirectory"`
	StartupFile          *valueObject.UnixFilePath          `json:"startupFile"`
	LogOutputPath        *valueObject.UnixFilePath          `json:"logOutputPath"`
	LogErrorPath         *valueObject.UnixFilePath          `json:"logErrorPath"`
	AvatarUrl            *valueObject.Url                   `json:"avatarUrl"`
	EstimatedSizeBytes   *valueObject.Byte                  `json:"estimatedSizeBytes"`
}

func NewInstallableService(
	manifestVersion valueObject.ServiceManifestVersion,
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	startCmd valueObject.UnixCommand,
	description valueObject.ServiceDescription,
	versions []valueObject.ServiceVersion,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopTimeoutSecs valueObject.UnixTime,
	stopSteps []valueObject.UnixCommand,
	installTimeoutSecs valueObject.UnixTime,
	installSteps []valueObject.UnixCommand,
	uninstallTimeoutSecs valueObject.UnixTime,
	uninstallSteps []valueObject.UnixCommand,
	uninstallFilePaths []valueObject.UnixFilePath,
	preStartTimeoutSecs valueObject.UnixTime,
	preStartSteps []valueObject.UnixCommand,
	postStartTimeoutSecs valueObject.UnixTime,
	postStartSteps []valueObject.UnixCommand,
	preStopTimeoutSecs valueObject.UnixTime,
	preStopSteps []valueObject.UnixCommand,
	postStopTimeoutSecs valueObject.UnixTime,
	postStopSteps []valueObject.UnixCommand,
	execUser *valueObject.UnixUsername,
	workingDirectory, startupFile, logOutputPath, logErrorPath *valueObject.UnixFilePath,
	avatarUrl *valueObject.Url,
	estimatedSizeBytes *valueObject.Byte,
) InstallableService {
	return InstallableService{
		ManifestVersion:      manifestVersion,
		Name:                 name,
		Nature:               nature,
		Type:                 serviceType,
		StartCmd:             startCmd,
		Description:          description,
		Versions:             versions,
		Envs:                 envs,
		PortBindings:         portBindings,
		StopTimeoutSecs:      stopTimeoutSecs,
		StopCmdSteps:         stopSteps,
		InstallTimeoutSecs:   installTimeoutSecs,
		InstallCmdSteps:      installSteps,
		UninstallTimeoutSecs: uninstallTimeoutSecs,
		UninstallCmdSteps:    uninstallSteps,
		UninstallFilePaths:   uninstallFilePaths,
		PreStartTimeoutSecs:  preStartTimeoutSecs,
		PreStartCmdSteps:     preStartSteps,
		PostStartTimeoutSecs: postStartTimeoutSecs,
		PostStartCmdSteps:    postStartSteps,
		PreStopTimeoutSecs:   preStopTimeoutSecs,
		PreStopCmdSteps:      preStopSteps,
		PostStopTimeoutSecs:  postStopTimeoutSecs,
		PostStopCmdSteps:     postStopSteps,
		ExecUser:             execUser,
		WorkingDirectory:     workingDirectory,
		StartupFile:          startupFile,
		LogOutputPath:        logOutputPath,
		LogErrorPath:         logErrorPath,
		AvatarUrl:            avatarUrl,
		EstimatedSizeBytes:   estimatedSizeBytes,
	}
}
