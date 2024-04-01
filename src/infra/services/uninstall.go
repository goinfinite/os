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
	mariaDbPidFilePath := "/run/mysqld/mysqld.pid"
	mariaDbPid, err := infraHelper.GetFileContent(mariaDbPidFilePath)
	if err != nil {
		return err
	}

	mariaDbPidWithoutBreakLine := strings.Trim(mariaDbPid, "\n")
	_, err = infraHelper.RunCmd(
		"kill",
		"-15",
		mariaDbPidWithoutBreakLine,
	)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	_, err = infraHelper.RunCmd(
		"mariadbd-safe",
		"stop",
	)
	if err != nil {
		return err
	}

	mariaDbDataDirPath := "/var/lib/mysql"
	err = os.RemoveAll(mariaDbDataDirPath)
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
