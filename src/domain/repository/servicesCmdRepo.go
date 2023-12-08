package repository

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
)

type ServicesCmdRepo interface {
	AddInstallable(addInstallableService dto.AddInstallableService) error
	Start(name valueObject.ServiceName) error
	Stop(name valueObject.ServiceName) error
	Install(name valueObject.ServiceName, version *valueObject.ServiceVersion) error
	Uninstall(name valueObject.ServiceName) error
}
