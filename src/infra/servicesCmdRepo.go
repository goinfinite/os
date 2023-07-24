package infra

import (
	"github.com/speedianet/sam/src/domain/valueObject"
	servicesInfra "github.com/speedianet/sam/src/infra/services"
)

type ServicesCmdRepo struct {
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

func (repo ServicesCmdRepo) Install(
	name valueObject.ServiceName,
	version *valueObject.ServiceVersion,
) error {
	err := servicesInfra.Install(name, version)
	if err != nil {
		return err
	}

	err = servicesInfra.SupervisordFacade{}.Reload()
	if err != nil {
		return err
	}

	return nil
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
