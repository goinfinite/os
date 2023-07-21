package infra

import (
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
	servicesInfra "github.com/speedianet/sam/src/infra/services"
)

type ServicesCmdRepo struct {
}

func (repo ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	return servicesInfra.SupervisordFacade{}.Start(name)
}

func (repo ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	err := servicesInfra.SupervisordFacade{}.Stop(name)
	if err != nil {
		return err
	}

	switch name.String() {
	case "openlitespeed":
		infraHelper.RunCmd(
			"/usr/local/lsws/bin/lswsctrl",
			"stop",
		)
		infraHelper.RunCmd(
			"pkill",
			"lsphp",
		)
	}

	return nil
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
	return servicesInfra.Uninstall(name)
}
