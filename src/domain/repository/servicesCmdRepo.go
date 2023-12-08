package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesCmdRepo interface {
	AddInstallable(addDto dto.AddInstallableService) error
	AddCustom(addDto dto.AddCustomService) error
	Start(name valueObject.ServiceName) error
	Stop(name valueObject.ServiceName) error
	Uninstall(name valueObject.ServiceName) error
}
