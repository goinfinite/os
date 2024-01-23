package infra

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	"golang.org/x/exp/slices"
)

type RuntimeCmdRepo struct {
}

func (repo RuntimeCmdRepo) restartPhp() error {
	phpSvcName, _ := valueObject.NewServiceName("php")
	err := servicesInfra.ServicesCmdRepo{}.Restart(phpSvcName)
	if err != nil {
		return errors.New("RestartWebServerFailed: " + err.Error())
	}

	return nil
}

func (repo RuntimeCmdRepo) UpdatePhpVersion(
	hostname valueObject.Fqdn,
	version valueObject.PhpVersion,
) error {
	queryRepo := RuntimeQueryRepo{}

	phpVersion, err := queryRepo.GetPhpVersion(hostname)
	if err != nil {
		return err
	}

	if phpVersion.Value == version {
		return nil
	}

	vhconfFile, err := queryRepo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return err
	}

	newLsapiLine := "lsapi:lsphp" + version.GetWithoutDots()
	_, err = infraHelper.RunCmd(
		"sed",
		"-i",
		"s/lsapi:lsphp[0-9][0-9]/"+newLsapiLine+"/g",
		vhconfFile,
	)
	if err != nil {
		return errors.New("UpdatePhpVersionFailed: " + err.Error())
	}

	return repo.restartPhp()
}

func (repo RuntimeCmdRepo) UpdatePhpSettings(
	hostname valueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	vhconfFile, err := RuntimeQueryRepo{}.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return err
	}

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
			log.Printf("UpdatePhpSettingFailed: %s", err.Error())
			continue
		}
	}

	return repo.restartPhp()
}

func (repo RuntimeCmdRepo) EnablePhpModule(
	phpVersion valueObject.PhpVersion,
	module entity.PhpModule,
) error {
	lsphpDir := "/usr/local/lsws/lsphp" + phpVersion.GetWithoutDots()
	iniRootDir := lsphpDir + "/etc/php/" + phpVersion.String()
	modsAvailableDir := iniRootDir + "/mods-available"
	modsDisabledDir := iniRootDir + "/mods-disabled"

	disabledInitFile, err := infraHelper.GetFilePathWithMatch(
		modsDisabledDir,
		module.Name.String()+".ini",
	)
	if err == nil {
		enabledIniFile := strings.Replace(
			disabledInitFile,
			modsDisabledDir,
			modsAvailableDir,
			1,
		)

		os.Rename(
			disabledInitFile,
			enabledIniFile,
		)

		return nil
	}

	lsphpPkgPrefix := "lsphp" + phpVersion.GetWithoutDots() + "-"
	err = infraHelper.InstallPkgs([]string{
		lsphpPkgPrefix + module.Name.String(),
	})
	if err == nil {
		return nil
	}

	err = infraHelper.InstallPkgs([]string{
		lsphpPkgPrefix + "pear",
	})
	if err != nil {
		return errors.New("InstallPhpPearModuleFailed: " + err.Error())
	}

	os.Symlink("/bin/sed", "/usr/bin/sed")

	switch module.Name.String() {
	case "mcrypt":
		err = infraHelper.InstallPkgs([]string{
			"libmcrypt-dev", "libmcrypt4",
		})
		if err != nil {
			return errors.New("InstallLibmcryptFailed: " + err.Error())
		}
	case "ssh2":
		err = infraHelper.InstallPkgs([]string{
			"libssh2-1-dev", "libssh2-1",
		})
		if err != nil {
			return errors.New("InstallLibssh2Failed: " + err.Error())
		}
	case "yaml":
		err = infraHelper.InstallPkgs([]string{
			"libyaml-dev",
		})
		if err != nil {
			return errors.New("InstallLibyamlFailed: " + err.Error())
		}
	case "xdebug", "parallel", "swoole", "sqlsrv":
		if phpVersion == "7.4" {
			return errors.New("PhpVersionUnsupportedByModule: " + phpVersion.String())
		}
	}

	_, err = infraHelper.RunCmd(
		"bash",
		"-c",
		"echo | "+lsphpDir+"/bin/pecl install "+module.Name.String(),
	)
	if err != nil {
		return errors.New("InstallPeclModuleFailed: " + err.Error())
	}

	err = infraHelper.UpdateFile(
		modsAvailableDir+"/"+module.Name.String()+".ini",
		"extension="+module.Name.String()+".so",
		true,
	)
	if err != nil {
		return errors.New("CreatePhpModuleIniFileFailed: " + err.Error())
	}

	return nil
}

func (repo RuntimeCmdRepo) DisablePhpModule(
	phpVersion valueObject.PhpVersion,
	module entity.PhpModule,
) error {
	iniRootDir := "/usr/local/lsws/lsphp" +
		phpVersion.GetWithoutDots() + "/etc/php/" + phpVersion.String()
	modsAvailableDir := iniRootDir + "/mods-available"
	modsDisabledDir := iniRootDir + "/mods-disabled"

	enabledIniFile, err := infraHelper.GetFilePathWithMatch(
		modsAvailableDir,
		module.Name.String()+".ini",
	)
	if err != nil {
		return errors.New("PhpModuleIniFileNotFound: " + err.Error())
	}
	disabledIniFile := strings.Replace(
		enabledIniFile,
		modsAvailableDir,
		modsDisabledDir,
		1,
	)

	os.Mkdir(modsDisabledDir, 0755)
	err = os.Rename(
		enabledIniFile,
		disabledIniFile,
	)
	if err != nil {
		return errors.New("DisablePhpModuleFailed: " + err.Error())
	}

	return nil
}

func (repo RuntimeCmdRepo) UpdatePhpModules(
	hostname valueObject.Fqdn,
	modules []entity.PhpModule,
) error {
	phpVersion, err := RuntimeQueryRepo{}.GetPhpVersion(hostname)
	if err != nil {
		return err
	}

	allModules, err := RuntimeQueryRepo{}.GetPhpModules(phpVersion.Value)
	if err != nil {
		return err
	}

	var activeModules []string
	for _, module := range allModules {
		if module.Status {
			activeModules = append(activeModules, module.Name.String())
		}
	}

	for _, module := range modules {
		isModuleEnabled := slices.Contains(activeModules, module.Name.String())
		if isModuleEnabled && module.Status {
			continue
		}

		if module.Status {
			err := repo.EnablePhpModule(phpVersion.Value, module)
			if err != nil {
				continue
			}
			continue
		}

		err := repo.DisablePhpModule(phpVersion.Value, module)
		if err != nil {
			continue
		}
	}

	return repo.restartPhp()
}
