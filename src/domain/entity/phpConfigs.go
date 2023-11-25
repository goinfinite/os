package entity

import (
	"github.com/speedianet/os/src/domain/valueObject"
)

type PhpConfigs struct {
	Hostname valueObject.Fqdn `json:"hostname"`
	Version  PhpVersion       `json:"version"`
	Settings []PhpSetting     `json:"settings"`
	Modules  []PhpModule      `json:"modules"`
}

func NewPhpConfigs(
	hostname valueObject.Fqdn,
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
