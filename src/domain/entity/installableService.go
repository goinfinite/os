package entity

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type InstallableService struct {
	ManifestVersion      valueObject.ServiceManifestVersion    `json:"manifestVersion"`
	Name                 valueObject.ServiceName               `json:"name"`
	Nature               valueObject.ServiceNature             `json:"nature"`
	Type                 valueObject.ServiceType               `json:"type"`
	StartCmd             tkValueObject.UnixCommand             `json:"startCmd"`
	Description          valueObject.ServiceDescription        `json:"description"`
	Versions             []valueObject.ServiceVersion          `json:"versions"`
	Envs                 []valueObject.ServiceEnv              `json:"envs"`
	PortBindings         []valueObject.PortBinding             `json:"portBindings"`
	StopTimeoutSecs      tkValueObject.UnixTime                `json:"-"`
	StopCmdSteps         []tkValueObject.UnixCommand           `json:"-"`
	InstallTimeoutSecs   tkValueObject.UnixTime                `json:"-"`
	InstallCmdSteps      []tkValueObject.UnixCommand           `json:"-"`
	UninstallTimeoutSecs tkValueObject.UnixTime                `json:"-"`
	UninstallCmdSteps    []tkValueObject.UnixCommand           `json:"-"`
	UninstallFilePaths   []tkValueObject.UnixAbsoluteFilePath  `json:"-"`
	PreStartTimeoutSecs  tkValueObject.UnixTime                `json:"-"`
	PreStartCmdSteps     []tkValueObject.UnixCommand           `json:"-"`
	PostStartTimeoutSecs tkValueObject.UnixTime                `json:"-"`
	PostStartCmdSteps    []tkValueObject.UnixCommand           `json:"-"`
	PreStopTimeoutSecs   tkValueObject.UnixTime                `json:"-"`
	PreStopCmdSteps      []tkValueObject.UnixCommand           `json:"-"`
	PostStopTimeoutSecs  tkValueObject.UnixTime                `json:"-"`
	PostStopCmdSteps     []tkValueObject.UnixCommand           `json:"-"`
	ExecUser             *tkValueObject.UnixUsername            `json:"execUser"`
	WorkingDirectory     *tkValueObject.UnixAbsoluteFilePath   `json:"workingDirectory"`
	StartupFile          *tkValueObject.UnixAbsoluteFilePath   `json:"startupFile"`
	LogOutputPath        *tkValueObject.UnixAbsoluteFilePath   `json:"logOutputPath"`
	LogErrorPath         *tkValueObject.UnixAbsoluteFilePath   `json:"logErrorPath"`
	AvatarUrl            *tkValueObject.Url                    `json:"avatarUrl"`
	EstimatedSizeBytes   *tkValueObject.Byte                   `json:"estimatedSizeBytes"`
}

func NewInstallableService(
	manifestVersion valueObject.ServiceManifestVersion,
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	startCmd tkValueObject.UnixCommand,
	description valueObject.ServiceDescription,
	versions []valueObject.ServiceVersion,
	envs []valueObject.ServiceEnv,
	portBindings []valueObject.PortBinding,
	stopTimeoutSecs tkValueObject.UnixTime,
	stopSteps []tkValueObject.UnixCommand,
	installTimeoutSecs tkValueObject.UnixTime,
	installSteps []tkValueObject.UnixCommand,
	uninstallTimeoutSecs tkValueObject.UnixTime,
	uninstallSteps []tkValueObject.UnixCommand,
	uninstallFilePaths []tkValueObject.UnixAbsoluteFilePath,
	preStartTimeoutSecs tkValueObject.UnixTime,
	preStartSteps []tkValueObject.UnixCommand,
	postStartTimeoutSecs tkValueObject.UnixTime,
	postStartSteps []tkValueObject.UnixCommand,
	preStopTimeoutSecs tkValueObject.UnixTime,
	preStopSteps []tkValueObject.UnixCommand,
	postStopTimeoutSecs tkValueObject.UnixTime,
	postStopSteps []tkValueObject.UnixCommand,
	execUser *tkValueObject.UnixUsername,
	workingDirectory, startupFile, logOutputPath, logErrorPath *tkValueObject.UnixAbsoluteFilePath,
	avatarUrl *tkValueObject.Url,
	estimatedSizeBytes *tkValueObject.Byte,
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
