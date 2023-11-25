package servicesInfra

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

func Uninstall(name valueObject.ServiceName) error {
	err := SupervisordFacade{}.RemoveConf(name.String())
	if err != nil {
		return err
	}

	var packages []string
	switch name.String() {
	case "php":
		packages = append(OlsPackages, "lsphp*")
	case "node":
		packages = NodePackages
	case "mysql":
		packages = MariaDbPackages
	case "redis":
		packages = RedisPackages
	default:
		log.Printf("ServiceNotImplemented: %s", name.String())
		return errors.New("ServiceNotImplemented")
	}

	err = SupervisordFacade{}.Stop(name)
	if err != nil {
		log.Printf("UninstallServiceError: %s", err.Error())
		return errors.New("UninstallServiceError")
	}

	err = SupervisordFacade{}.RemoveConf(name.String())
	if err != nil {
		log.Printf("UninstallServiceError: %s", err.Error())
		return errors.New("UninstallServiceError")
	}

	purgeEnvVars := map[string]string{
		"DEBIAN_FRONTEND": "noninteractive",
	}
	purgePackages := append([]string{"purge", "-y"}, packages...)
	_, err = infraHelper.RunCmdWithEnvVars("apt-get", purgeEnvVars, purgePackages...)
	if err != nil {
		log.Printf("UninstallServiceError: %s", err.Error())
		return errors.New("UninstallServiceError")
	}

	return nil
}
