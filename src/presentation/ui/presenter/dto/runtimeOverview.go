package presenterDto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type RuntimeOverview struct {
	VirtualHostHostname            valueObject.Fqdn        `json:"vhostHostname"`
	Type                           valueObject.RuntimeType `json:"type"`
	IsInstalled                    bool                    `json:"-"`
	IsServiceMappingAlreadyCreated bool                    `json:"-"`
	PhpConfigs                     *entity.PhpConfigs      `json:"phpConfigs"`
}

func NewRuntimeOverview(
	virtualHostHostname valueObject.Fqdn,
	runtimeType valueObject.RuntimeType,
	isInstalled, isServiceMappingAlreadyCreated bool,
	phpConfigs *entity.PhpConfigs,
) RuntimeOverview {
	return RuntimeOverview{
		VirtualHostHostname:            virtualHostHostname,
		Type:                           runtimeType,
		IsInstalled:                    isInstalled,
		IsServiceMappingAlreadyCreated: isServiceMappingAlreadyCreated,
		PhpConfigs:                     phpConfigs,
	}
}
