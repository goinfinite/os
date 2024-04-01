package servicesInfra

import (
	"os"

	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

func purgePkgs(packages []string) error {
	purgePackages := append([]string{"purge", "-y"}, packages...)
	_, err := infraHelper.RunCmd("apt-get", purgePackages...)
	if err != nil {
		return err
	}

	return nil
}

func removeMariaDb() error {
	dbDataDirPath := "/var/lib/mysql"
	err := os.RemoveAll(dbDataDirPath)
	if err != nil {
		return err
	}

	return purgePkgs(MariaDbPackages)
}

func Uninstall(name valueObject.ServiceName) error {
	err := SupervisordFacade{}.RemoveConf(name)
	if err != nil {
		return err
	}

	switch name.String() {
	case "php":
		packages := append(OlsPackages, "lsphp*")
		return purgePkgs(packages)
	case "mariadb":
		return removeMariaDb()
	case "redis":
		return purgePkgs(RedisPackages)
	default:
		return nil
	}
}
