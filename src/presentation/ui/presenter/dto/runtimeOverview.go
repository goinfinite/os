package presenterDto

import (
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type RuntimeOverview struct {
	RuntimeType             valueObject.RuntimeType `json:"type"`
	IsInstalled             bool                    `json:"-"`
	IsMappingAlreadyCreated bool                    `json:"-"`
	PhpConfigs              *entity.PhpConfigs      `json:"phpConfigs"`
}

func NewRuntimeOverview(
	runtimeType valueObject.RuntimeType,
	isInstalled, isMappingAlreadyCreated bool,
	phpConfigs *entity.PhpConfigs,
) RuntimeOverview {
	return RuntimeOverview{
		RuntimeType:             runtimeType,
		IsInstalled:             isInstalled,
		IsMappingAlreadyCreated: isMappingAlreadyCreated,
		PhpConfigs:              phpConfigs,
	}
}
