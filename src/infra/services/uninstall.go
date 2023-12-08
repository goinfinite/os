package servicesInfra

import (
	"errors"

	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

func Uninstall(name valueObject.ServiceName) error {
	err := SupervisordFacade{}.Stop(name)
	if err != nil {
		return err
	}

	err = SupervisordFacade{}.RemoveConf(name)
	if err != nil {
		return err
	}

	var packages []string
	switch name.String() {
	case "php":
		packages = append(OlsPackages, "lsphp*")
	case "node":
		packages = NodePackages
	case "mariadb":
		packages = MariaDbPackages
	case "redis":
		packages = RedisPackages
	default:
		return errors.New("ServiceUnknown")
	}

	purgeEnvVars := map[string]string{
		"DEBIAN_FRONTEND": "noninteractive",
	}
	purgePackages := append([]string{"purge", "-y"}, packages...)
	_, err = infraHelper.RunCmdWithEnvVars("apt-get", purgeEnvVars, purgePackages...)
	if err != nil {
		return err
	}

	return nil
}
