package infra

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
	servicesInfra "github.com/speedianet/sam/src/infra/services"
)

type ServicesCmdRepo struct {
}

const (
	supervisordCmd = "/usr/bin/supervisord"
)

func (repo ServicesCmdRepo) Start(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"start",
		name.String(),
	)
	if err != nil {
		log.Printf("StartServiceError: %s", err)
		return errors.New("StartServiceError")
	}

	return nil
}

func (repo ServicesCmdRepo) Stop(name valueObject.ServiceName) error {
	_, err := infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"stop",
		name.String(),
	)
	if err != nil {
		log.Printf("StopServiceError: %s", err)
		return errors.New("StopServiceError")
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

	_, err = infraHelper.RunCmd(
		supervisordCmd,
		"ctl",
		"reload",
	)
	if err != nil {
		log.Printf("ReloadSupervisorError: %s", err)
		return errors.New("ReloadSupervisorError")
	}

	return nil
}

func (repo ServicesCmdRepo) Uninstall(
	name valueObject.ServiceName,
) error {
	return servicesInfra.Uninstall(name)
}
