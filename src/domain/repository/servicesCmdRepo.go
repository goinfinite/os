package repository

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type ServicesCmdRepo interface {
	CreateInstallable(createDto dto.CreateInstallableService) (valueObject.ServiceName, error)
	CreateCustom(createDto dto.CreateCustomService) error
	Update(updateDto dto.UpdateService) error
	Delete(name valueObject.ServiceName) error
	RefreshInstallableItems() error
}
