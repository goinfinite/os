package servicesInfra

import (
	"errors"
	"log"

	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

func Uninstall(name valueObject.ServiceName) error {
	err := SupervisordFacade{}.RemoveConf(name.String())
	if err != nil {
		return err
	}

	var packages []string
	var supervisordConfigName string
	switch name.String() {
	case "openlitespeed", "litespeed":
		packages = OlsPackages
		supervisordConfigName = "openlitespeed"
	case "mysql", "mysqld", "maria", "mariadb", "percona", "perconadb":
		packages = MariaDbPackages
		supervisordConfigName = "mysql"
	case "node", "nodejs":
		packages = NodePackages
		supervisordConfigName = "node"
	case "redis", "redis-server":
		packages = RedisPackages
		supervisordConfigName = "redis"
	default:
		log.Printf("ServiceNotImplemented: %s", name.String())
		return errors.New("ServiceNotImplemented")
	}

	err = SupervisordFacade{}.Stop(name)
	if err != nil {
		log.Printf("UninstallServiceError: %s", err.Error())
		return errors.New("UninstallServiceError")
	}

	err = SupervisordFacade{}.RemoveConf(supervisordConfigName)
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
