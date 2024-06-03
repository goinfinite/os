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

func removeRedis() error {
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

func removePostgres() error {
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
	case "php", "php-webserver":
		packages := append(OlsPackages, "lsphp*")
		return purgePkgs(packages)
	case "mariadb":
		return removeMariaDb()
	case "redis":
		return removeRedis()
	case "postgresql":
		return removePostgres()
	default:
		return nil
	}
}
