package infra

import (
	"errors"

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
	phpVersion, err := RuntimeQueryRepo{}.GetPhpVersion(hostname)
	if err != nil {
		return err
	}

	if phpVersion.Value == version {
		return nil
	}

	vhconfFile := WsQueryRepo{}.GetVirtualHostConfFilePath(hostname)
	newLsapiLine := "lsapi:lsphp" + version.GetWithoutDots()
	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/lsapi:lsphp[0-9][0-9]/"+newLsapiLine+"/g",
		vhconfFile,
	)
	if err != nil {
		return errors.New("UpdatePhpVersionFailed")
	}

	err = ServicesCmdRepo{}.Restart(valueObject.NewServiceNamePanic("openlitespeed"))
	if err != nil {
		return errors.New("RestartWebServerFailed")
	}

	return nil
}

func (r RuntimeCmdRepo) UpdatePhpSettings(
	hostname valueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	vhconfFile := WsQueryRepo{}.GetVirtualHostConfFilePath(hostname)
	for _, setting := range settings {
		name := setting.Name.String()
		value := setting.Value.String()
		if setting.Value.GetType() == "string" {
			value = "\"" + value + "\""
		}

		_, err := infraHelper.RunCmd(
			"sed",
			"-i",
			"s/"+name+" .*/"+name+" "+value+"/g",
			vhconfFile,
		)
		if err != nil {
			continue
		}
	}

	err := ServicesCmdRepo{}.Restart(valueObject.NewServiceNamePanic("openlitespeed"))
	if err != nil {
		return errors.New("RestartWebServerFailed")
	}

	return nil
}

func (r RuntimeCmdRepo) UpdatePhpModules(
	hostname valueObject.Fqdn,
	modules []entity.PhpModule,
) error {
	return nil
}
