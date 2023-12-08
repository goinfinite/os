package infra

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

type ServicesCmdRepo struct {
}

func (repo ServicesCmdRepo) AddInstallable(
	addDto dto.AddInstallableService,
) error {
	return servicesInfra.AddInstallable(addDto)
}

func (repo ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	return servicesInfra.SupervisordFacade{}.Start(name)
}

func (repo ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	return servicesInfra.SupervisordFacade{}.Stop(name)
}

func (repo ServicesCmdRepo) Restart(name valueObject.ServiceName) error {
	return servicesInfra.SupervisordFacade{}.Restart(name)
}

func (repo ServicesCmdRepo) Uninstall(
	name valueObject.ServiceName,
) error {
	err := repo.Stop(name)
	if err != nil {
		return err
	}

	err = servicesInfra.Uninstall(name)
	if err != nil {
		return err
	}

	err = servicesInfra.SupervisordFacade{}.Reload()
	if err != nil {
		return err
	}

	return nil
}
