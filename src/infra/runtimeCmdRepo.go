package infra

import (
	"errors"
	"log"
	"os"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type RuntimeCmdRepo struct {
}

func (r RuntimeCmdRepo) UpdatePhpVersion(
	hostname valueObject.Fqdn,
	version valueObject.PhpVersion,
) error {
	vhconfFile := "/app/conf/vhconf.conf"
	mainVirtualHost := valueObject.NewFqdnPanic(os.Getenv("VIRTUAL_HOST"))
	if hostname != mainVirtualHost {
		vhconfFile = "/app/domains/" + string(hostname) + "/conf/vhconf.conf"
	}

	currentPhpVersionStr, err := infraHelper.RunCmd(
		"awk",
		"/lsapi:lsphp/ {gsub(/[^0-9]/, \"\", $2); print $2}",
		vhconfFile,
	)
	if err != nil {
		log.Printf("FailedToGetPhpVersion: %v", err)
		return errors.New("FailedToGetPhpVersion")
	}

	currentPhpVersion, err := valueObject.NewPhpVersion(currentPhpVersionStr)
	if err != nil {
		return errors.New("FailedToGetPhpVersion")
	}

	if currentPhpVersion == version {
		return nil
	}

	newLsapiLine := "lsapi:lsphp" + version.GetWithoutDots()
	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/lsapi:lsphp[0-9][0-9]/"+newLsapiLine+"/g",
		vhconfFile,
	)
	if err != nil {
		return errors.New("FailedToUpdatePhpVersion")
	}

	err = ServicesCmdRepo{}.Restart(valueObject.NewServiceNamePanic("openlitespeed"))
	if err != nil {
		return errors.New("FailedToRestartWebServer")
	}

	return nil
}

func (r RuntimeCmdRepo) UpdatePhpSettings(
	hostname valueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	return nil
}

func (r RuntimeCmdRepo) UpdatePhpModules(
	hostname valueObject.Fqdn,
	modules []entity.PhpModule,
) error {
	return nil
}
