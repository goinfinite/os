package servicesInfra

import (
	"os"
	"strings"
	"time"

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
	pidFilePath := "/run/mysqld/mysqld.pid"
	pid, err := infraHelper.GetFileContent(pidFilePath)
	if err != nil {
		return err
	}

	pidWithoutBreakLine := strings.Trim(pid, "\n")
	_, err = infraHelper.RunCmd(
		"kill",
		"-15",
		pidWithoutBreakLine,
	)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	dbDataDirPath := "/var/lib/mysql"
	err = os.RemoveAll(dbDataDirPath)
	if err != nil {
		return err
	}

	packages := append(MariaDbPackages, "mariadb*", "mysql*")
	return purgePkgs(packages)
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
