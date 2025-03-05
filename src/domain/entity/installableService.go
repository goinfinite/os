package entity

import "github.com/goinfinite/os/src/domain/valueObject"

type InstallableService struct {
	ManifestVersion    valueObject.ServiceManifestVersion `json:"manifestVersion"`
	Name               valueObject.ServiceName            `json:"name"`
	Nature             valueObject.ServiceNature          `json:"nature"`
	Type               valueObject.ServiceType            `json:"type"`
	StartCmd           valueObject.UnixCommand            `json:"startCmd"`
	Description        valueObject.ServiceDescription     `json:"description"`
	Versions           []valueObject.ServiceVersion       `json:"versions"`
	Envs               []valueObject.ServiceEnv           `json:"envs"`
	PortBindings       []valueObject.PortBinding          `json:"portBindings"`
	StopCmdSteps       []valueObject.UnixCommand          `json:"-"`
	InstallTimeoutSecs valueObject.UnixTime               `json:"installTimeoutSecs"`
	InstallCmdSteps    []valueObject.UnixCommand          `json:"-"`
	UninstallCmdSteps  []valueObject.UnixCommand          `json:"-"`
	UninstallFilePaths []valueObject.UnixFilePath         `json:"-"`
	PreStartCmdSteps   []valueObject.UnixCommand          `json:"-"`
	PostStartCmdSteps  []valueObject.UnixCommand          `json:"-"`
	PreStopCmdSteps    []valueObject.UnixCommand          `json:"-"`
	PostStopCmdSteps   []valueObject.UnixCommand          `json:"-"`
	ExecUser           *valueObject.UnixUsername          `json:"execUser"`
	WorkingDirectory   *valueObject.UnixFilePath          `json:"workingDirectory"`
	StartupFile        *valueObject.UnixFilePath          `json:"startupFile"`
	LogOutputPath      *valueObject.UnixFilePath          `json:"logOutputPath"`
	LogErrorPath       *valueObject.UnixFilePath          `json:"logErrorPath"`
	AvatarUrl          *valueObject.Url                   `json:"avatarUrl"`
	EstimatedSizeBytes *valueObject.Byte                  `json:"estimatedSizeBytes"`
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
	stopSteps []valueObject.UnixCommand,
	installTimeoutSecs valueObject.UnixTime,
	installSteps, uninstallSteps []valueObject.UnixCommand,
	uninstallFilePaths []valueObject.UnixFilePath,
	preStartSteps, postStartSteps, preStopSteps, postStopSteps []valueObject.UnixCommand,
	execUser *valueObject.UnixUsername,
	workingDirectory, startupFile, logOutputPath, logErrorPath *valueObject.UnixFilePath,
	avatarUrl *valueObject.Url,
	estimatedSizeBytes *valueObject.Byte,
) InstallableService {
	return InstallableService{
		ManifestVersion:    manifestVersion,
		Name:               name,
		Nature:             nature,
		Type:               serviceType,
		StartCmd:           startCmd,
		Description:        description,
		Versions:           versions,
		Envs:               envs,
		PortBindings:       portBindings,
		StopCmdSteps:       stopSteps,
		InstallTimeoutSecs: installTimeoutSecs,
		InstallCmdSteps:    installSteps,
		UninstallCmdSteps:  uninstallSteps,
		UninstallFilePaths: uninstallFilePaths,
		PreStartCmdSteps:   preStartSteps,
		PostStartCmdSteps:  postStartSteps,
		PreStopCmdSteps:    preStopSteps,
		PostStopCmdSteps:   postStopSteps,
		ExecUser:           execUser,
		WorkingDirectory:   workingDirectory,
		StartupFile:        startupFile,
		LogOutputPath:      logOutputPath,
		LogErrorPath:       logErrorPath,
		AvatarUrl:          avatarUrl,
		EstimatedSizeBytes: estimatedSizeBytes,
	}
}
