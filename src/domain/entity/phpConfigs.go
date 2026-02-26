package entity

import tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"

type PhpConfigs struct {
	Hostname tkValueObject.Fqdn `json:"hostname"`
	Version  PhpVersion         `json:"version"`
	Settings []PhpSetting       `json:"settings"`
	Modules  []PhpModule        `json:"modules"`
}

func NewPhpConfigs(
	hostname tkValueObject.Fqdn,
	version PhpVersion,
	settings []PhpSetting,
	modules []PhpModule,
) PhpConfigs {
	return PhpConfigs{
		Hostname: hostname,
		Version:  version,
		Settings: settings,
		Modules:  modules,
	}
}
