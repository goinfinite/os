package servicesInfra

import (
	"errors"
	"os"

	"github.com/speedianet/os/src/domain/valueObject"
	cronInfra "github.com/speedianet/os/src/infra/cron"
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

func uninstallPhpWebserver() error {
	cronQueryRepo := cronInfra.CronQueryRepo{}
	crons, err := cronQueryRepo.Get()
	if err != nil {
		return errors.New("ReadCronsError")
	}

	cronCmdRepo, err := cronInfra.NewCronCmdRepo()
	if err != nil {
		return errors.New("CreateCronCmdRepoError: " + err.Error())
	}

	for _, cron := range crons {
		if cron.Comment.String() != PhpWebserverAutoReloadCronComment {
			continue
		}

		err = cronCmdRepo.Delete(cron.Id)
		if err != nil {
			return errors.New("DeleteAutoReloadCronError")
		}
	}

	packages := append(OlsPackages, "lsphp*")
	return purgePkgs(packages)
}

func uninstallMariaDb() error {
	pathsToRemove := []string{
		"/etc/mysql",
		"/var/lib/mysql",
		"/var/log/mysql",
		"/etc/apt/sources.list.d/mariadb.list",
		"/root/.my.cnf",
	}

	for _, path := range pathsToRemove {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return purgePkgs(MariaDbPackages)
}

func uninstallRedis() error {
	pathsToRemove := []string{
		"/etc/redis",
		"/var/lib/redis",
		"/var/log/redis",
		"/etc/apt/sources.list.d/redis-server.list",
	}

	for _, path := range pathsToRemove {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return purgePkgs(RedisPackages)
}

func uninstallPostgreSql() error {
	pathsToRemove := []string{
		"/etc/postgresql",
		"/var/lib/postgresql",
		"/var/log/postgresql",
		"/usr/lib/postgresql",
		"/etc/apt/sources.list.d/pgdg.list",
		"/root/.pgpass",
	}

	for _, path := range pathsToRemove {
		err := os.RemoveAll(path)
		if err != nil {
			return err
		}
	}

	return purgePkgs([]string{"postgresql*"})
}

func Uninstall(name valueObject.ServiceName) error {
	err := SupervisordFacade{}.RemoveConf(name)
	if err != nil {
		return err
	}

	switch name.String() {
	case "php-webserver", "php":
		return uninstallPhpWebserver()
	case "mariadb":
		return uninstallMariaDb()
	case "redis":
		return uninstallRedis()
	case "postgresql", "postgres":
		return uninstallPostgreSql()
	default:
		return nil
	}
}
