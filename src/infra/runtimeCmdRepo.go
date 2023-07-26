package infra

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
	"golang.org/x/exp/slices"
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
			log.Printf("UpdatePhpSettingFailed: %s", err.Error())
			continue
		}
	}

	err := ServicesCmdRepo{}.Restart(valueObject.NewServiceNamePanic("openlitespeed"))
	if err != nil {
		return errors.New("RestartWebServerFailed")
	}

	return nil
}

func (r RuntimeCmdRepo) EnablePhpModule(
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
		log.Printf("InstallPhpPearModuleFailed: %s", err.Error())
		return errors.New("InstallPhpPearModuleFailed")
	}

	os.Symlink("/bin/sed", "/usr/bin/sed")

	switch module.Name.String() {
	case "mcrypt":
		err = infraHelper.InstallPkgs([]string{
			"libmcrypt-dev", "libmcrypt4",
		})
		if err != nil {
			log.Printf("InstallLibmcryptFailed: %s", err.Error())
			return errors.New("InstallLibmcryptFailed")
		}
	case "ssh2":
		err = infraHelper.InstallPkgs([]string{
			"libssh2-1-dev", "libssh2-1",
		})
		if err != nil {
			log.Printf("InstallLibssh2Failed: %s", err.Error())
			return errors.New("InstallLibssh2Failed")
		}
	case "yaml":
		err = infraHelper.InstallPkgs([]string{
			"libyaml-dev",
		})
		if err != nil {
			log.Printf("InstallLibyamlFailed: %s", err.Error())
			return errors.New("InstallLibyamlFailed")
		}
	case "xdebug", "parallel", "swoole", "sqlsrv":
		if phpVersion == "7.4" {
			log.Printf("PhpVersionUnsupportedByModule: %s", phpVersion)
			return errors.New("PhpVersionUnsupportedByModule")
		}
	}

	_, err = infraHelper.RunCmd(
		"bash",
		"-c",
		"echo | "+lsphpDir+"/bin/pecl install "+module.Name.String(),
	)
	if err != nil {
		log.Printf("InstallPeclModuleFailed: %s", err.Error())
		return errors.New("InstallPeclModuleFailed")
	}

	err = infraHelper.UpdateFile(
		modsAvailableDir+"/"+module.Name.String()+".ini",
		"extension="+module.Name.String()+".so",
		true,
	)
	if err != nil {
		log.Printf("CreatePhpModuleIniFileFailed: %s", err.Error())
		return errors.New("CreatePhpModuleIniFileFailed")
	}

	return nil
}

func (r RuntimeCmdRepo) DisablePhpModule(
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
		log.Printf("PhpModuleIniFileNotFound: %s", err.Error())
		return errors.New("PhpModuleIniFileNotFound")
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
		log.Printf("DisablePhpModuleFailed: %s", err.Error())
		return errors.New("DisablePhpModuleFailed")
	}

	return nil
}

func (r RuntimeCmdRepo) UpdatePhpModules(
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
			err := r.EnablePhpModule(phpVersion.Value, module)
			if err != nil {
				continue
			}
			continue
		}

		err := r.DisablePhpModule(phpVersion.Value, module)
		if err != nil {
			continue
		}
	}

	err = ServicesCmdRepo{}.Restart(valueObject.NewServiceNamePanic("openlitespeed"))
	if err != nil {
		return errors.New("RestartWebServerFailed")
	}

	return nil
}
