package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesCmdRepo interface {
	AddInstallable(addDto dto.CreateInstallableService) error
	AddCustom(addDto dto.CreateCustomService) error
	Start(name valueObject.ServiceName) error
	Stop(name valueObject.ServiceName) error
	Update(updateDto dto.UpdateService) error
	Uninstall(name valueObject.ServiceName) error
}
