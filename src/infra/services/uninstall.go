package servicesInfra

import (
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
		return nil
	}

	purgePackages := append([]string{"purge", "-y"}, packages...)
	_, err = infraHelper.RunCmd("apt-get", purgePackages...)
	if err != nil {
		return err
	}

	return nil
}
