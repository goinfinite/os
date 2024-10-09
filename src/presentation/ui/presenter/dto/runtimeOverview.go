package presenterDto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type RuntimeOverview struct {
	VirtualHostHostname       valueObject.Fqdn        `json:"vhostHostname"`
	Type                      valueObject.RuntimeType `json:"type"`
	IsInstalled               bool                    `json:"-"`
	IsVirtualHostUsingRuntime bool                    `json:"-"`
	PhpConfigs                *entity.PhpConfigs      `json:"phpConfigs"`
}

func NewRuntimeOverview(
	virtualHostHostname valueObject.Fqdn,
	runtimeType valueObject.RuntimeType,
	isInstalled, isVirtualHostUsingRuntime bool,
	phpConfigs *entity.PhpConfigs,
) RuntimeOverview {
	return RuntimeOverview{
		VirtualHostHostname:       virtualHostHostname,
		Type:                      runtimeType,
		IsInstalled:               isInstalled,
		IsVirtualHostUsingRuntime: isVirtualHostUsingRuntime,
		PhpConfigs:                phpConfigs,
	}
}
