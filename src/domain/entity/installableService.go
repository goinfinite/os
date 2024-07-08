package entity

import "github.com/speedianet/os/src/domain/valueObject"

type InstallableService struct {
	Name               valueObject.ServiceName        `json:"name"`
	Nature             valueObject.ServiceNature      `json:"nature"`
	Type               valueObject.ServiceType        `json:"type"`
	Command            valueObject.UnixCommand        `json:"command"`
	Description        valueObject.ServiceDescription `json:"description"`
	Versions           []valueObject.ServiceVersion   `json:"versions"`
	PortBindings       []valueObject.PortBinding      `json:"portBindings"`
	InstallCmdSteps    []valueObject.UnixCommand      `json:"-"`
	UninstallCmdSteps  []valueObject.UnixCommand      `json:"-"`
	UninstallFileNames []valueObject.UnixFileName     `json:"-"`
	PreStartCmdSteps   []valueObject.UnixCommand      `json:"-"`
	PostStartCmdSteps  []valueObject.UnixCommand      `json:"-"`
	PreStopCmdSteps    []valueObject.UnixCommand      `json:"-"`
	PostStopCmdSteps   []valueObject.UnixCommand      `json:"-"`
	StartupFile        *valueObject.UnixFilePath      `json:"startupFile"`
	EstimatedSizeBytes *valueObject.Byte              `json:"estimatedSizeBytes"`
	AvatarUrl          *valueObject.Url               `json:"avatarUrl"`
}

func NewInstallableService(
	name valueObject.ServiceName,
	nature valueObject.ServiceNature,
	serviceType valueObject.ServiceType,
	command valueObject.UnixCommand,
	description valueObject.ServiceDescription,
	versions []valueObject.ServiceVersion,
	portBindings []valueObject.PortBinding,
	installCmdSteps []valueObject.UnixCommand,
	uninstallCmdSteps []valueObject.UnixCommand,
	uninstallFileNames []valueObject.UnixFileName,
	preStartCmdSteps []valueObject.UnixCommand,
	postStartCmdSteps []valueObject.UnixCommand,
	preStopCmdSteps []valueObject.UnixCommand,
	postStopCmdSteps []valueObject.UnixCommand,
	startupFile *valueObject.UnixFilePath,
	estimatedSizeBytes *valueObject.Byte,
	avatarUrl *valueObject.Url,
) InstallableService {
	return InstallableService{
		Name:               name,
		Nature:             nature,
		Type:               serviceType,
		Command:            command,
		Description:        description,
		Versions:           versions,
		PortBindings:       portBindings,
		InstallCmdSteps:    installCmdSteps,
		UninstallCmdSteps:  uninstallCmdSteps,
		UninstallFileNames: uninstallFileNames,
		PreStartCmdSteps:   preStartCmdSteps,
		PostStartCmdSteps:  postStartCmdSteps,
		PreStopCmdSteps:    preStopCmdSteps,
		PostStopCmdSteps:   postStopCmdSteps,
		StartupFile:        startupFile,
		EstimatedSizeBytes: estimatedSizeBytes,
		AvatarUrl:          avatarUrl,
	}
}
