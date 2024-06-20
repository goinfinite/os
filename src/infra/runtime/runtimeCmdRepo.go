package runtimeInfra

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	"github.com/speedianet/os/src/infra/infraData"
	servicesInfra "github.com/speedianet/os/src/infra/services"
)

type RuntimeCmdRepo struct {
	runtimeQueryRepo RuntimeQueryRepo
}

func NewRuntimeCmdRepo() *RuntimeCmdRepo {
	return &RuntimeCmdRepo{runtimeQueryRepo: RuntimeQueryRepo{}}
}

func (repo *RuntimeCmdRepo) restartPhp() error {
	phpSvcName, _ := valueObject.NewServiceName("php-webserver")
	servicesCmdRepo := servicesInfra.ServicesCmdRepo{}
	err := servicesCmdRepo.Restart(phpSvcName)
	if err != nil {
		return errors.New("RestartWebServerFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) UpdatePhpVersion(
	hostname valueObject.Fqdn,
	version valueObject.PhpVersion,
) error {
	phpVersion, err := repo.runtimeQueryRepo.ReadPhpVersion(hostname)
	if err != nil {
		return err
	}

	if phpVersion.Value == version {
		return nil
	}

	phpConfFilePath, err := repo.runtimeQueryRepo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return err
	}

	newLsapiLine := "lsapi:lsphp" + version.GetWithoutDots()
	_, err = infraHelper.RunCmdWithSubShell(
		"sed -i 's/lsapi:lsphp[0-9][0-9]/" + newLsapiLine + "/g' " + phpConfFilePath.String(),
	)
	if err != nil {
		return errors.New("UpdatePhpVersionFailed: " + err.Error())
	}

	return repo.restartPhp()
}

func (repo *RuntimeCmdRepo) UpdatePhpSettings(
	hostname valueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	phpConfFilePath, err := repo.runtimeQueryRepo.GetVirtualHostPhpConfFilePath(hostname)
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
			"sed", "-i", "s/"+name+" .*/"+name+" "+value+"/g", phpConfFilePath.String(),
		)
		if err != nil {
			log.Printf("(%s) UpdatePhpSettingFailed: %s", name, err.Error())
			continue
		}
	}

	return repo.restartPhp()
}

func (repo *RuntimeCmdRepo) EnablePhpModule(
	phpVersion valueObject.PhpVersion,
	module entity.PhpModule,
) error {
	lsphpDir := "/usr/local/lsws/lsphp" + phpVersion.GetWithoutDots()
	iniRootDir := lsphpDir + "/etc/php/" + phpVersion.String()
	modsAvailableDir := iniRootDir + "/mods-available"
	modsDisabledDir := iniRootDir + "/mods-disabled"

	moduleNameStr := module.Name.String()
	disabledInitFile, err := infraHelper.GetFilePathWithMatch(
		modsDisabledDir, moduleNameStr+".ini",
	)
	if err == nil {
		enabledIniFile := strings.Replace(
			disabledInitFile, modsDisabledDir, modsAvailableDir, 1,
		)

		os.Rename(disabledInitFile, enabledIniFile)
		return nil
	}

	lsphpPkgPrefix := "lsphp" + phpVersion.GetWithoutDots() + "-"
	err = infraHelper.InstallPkgs([]string{lsphpPkgPrefix + moduleNameStr})
	if err == nil {
		return nil
	}

	err = infraHelper.InstallPkgs([]string{lsphpPkgPrefix + "pear"})
	if err != nil {
		return errors.New("InstallPhpPearModuleFailed: " + err.Error())
	}

	_ = os.Symlink("/bin/sed", "/usr/bin/sed")

	dependenciesToInstall := []string{}
	// cSpell:disable
	switch moduleNameStr {
	case "mcrypt":
		dependenciesToInstall = []string{"libmcrypt-dev", "libmcrypt4"}
	case "ssh2":
		dependenciesToInstall = []string{"libssh2-1-dev", "libssh2-1"}
	case "yaml":
		dependenciesToInstall = []string{"libyaml-dev"}
	case "xdebug", "parallel", "swoole", "sqlsrv":
		if phpVersion == "7.4" {
			return errors.New("PhpVersionUnsupportedByModule: " + phpVersion.String())
		}
	}
	// cSpell:enable
	err = infraHelper.InstallPkgs(dependenciesToInstall)
	if err != nil {
		return errors.New("InstallModuleFailed: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell(
		"echo | " + lsphpDir + "/bin/pecl install " + moduleNameStr,
	)
	if err != nil {
		return errors.New("InstallPeclModuleFailed: " + err.Error())
	}

	moduleConfigFilePath := modsAvailableDir + "/" + moduleNameStr + ".ini"
	moduleConfigFileContent := "extension=" + moduleNameStr + ".so"
	err = infraHelper.UpdateFile(moduleConfigFilePath, moduleConfigFileContent, true)
	if err != nil {
		return errors.New("CreatePhpModuleIniFileFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) DisablePhpModule(
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
		enabledIniFile, modsAvailableDir, modsDisabledDir, 1,
	)

	os.Mkdir(modsDisabledDir, 0755)
	err = os.Rename(enabledIniFile, disabledIniFile)
	if err != nil {
		return errors.New("DisablePhpModuleFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) UpdatePhpModules(
	hostname valueObject.Fqdn,
	modules []entity.PhpModule,
) error {
	phpVersion, err := repo.runtimeQueryRepo.ReadPhpVersion(hostname)
	if err != nil {
		return err
	}

	allModules, err := repo.runtimeQueryRepo.ReadPhpModules(phpVersion.Value)
	if err != nil {
		return err
	}

	activeModuleNames := map[string]interface{}{}
	for _, module := range allModules {
		if !module.Status {
			continue
		}

		activeModuleNames[module.Name.String()] = nil
	}

	for _, module := range modules {
		shouldEnable := module.Status
		_, isModuleCurrentlyEnabled := activeModuleNames[module.Name.String()]

		if shouldEnable {
			if isModuleCurrentlyEnabled {
				continue
			}

			err := repo.EnablePhpModule(phpVersion.Value, module)
			if err != nil {
				continue
			}

			continue
		}

		if !isModuleCurrentlyEnabled {
			continue
		}

		err := repo.DisablePhpModule(phpVersion.Value, module)
		if err != nil {
			continue
		}
	}

	return repo.restartPhp()
}

func (repo *RuntimeCmdRepo) CreatePhpVirtualHost(hostname valueObject.Fqdn) error {
	vhostExists := true

	phpConfFilePath, err := repo.runtimeQueryRepo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		if err.Error() != "VirtualHostNotFound" {
			return err
		}
		vhostExists = false
	}

	if vhostExists {
		return nil
	}

	phpConfFilePathStr := phpConfFilePath.String()
	templatePhpVhostConfFilePath := "/app/conf/php-webserver/template"
	err = infraHelper.CopyFile(templatePhpVhostConfFilePath, phpConfFilePathStr)
	if err != nil {
		return errors.New("CopyPhpConfTemplateError: " + err.Error())
	}

	hostnameStr := hostname.String()
	_, err = infraHelper.RunCmd(
		"sed", "-ie", "s/speedia.net/"+hostnameStr+"/g", phpConfFilePathStr,
	)
	if err != nil {
		return errors.New("UpdatePhpVirtualHostConfFileError: " + err.Error())
	}

	phpVhostHttpdConf := `
virtualhost ` + hostname.String() + ` {
  vhRoot                  /app/html/` + hostnameStr + `/
  configFile              ` + phpConfFilePathStr + `
  allowSymbolLink         1
  enableScript            1
  restrained              0
  setUIDMode              0
}
`
	shouldOverwrite := false
	err = infraHelper.UpdateFile(
		infraData.GlobalConfigs.OlsHttpdConfFilePath, phpVhostHttpdConf, shouldOverwrite,
	)
	if err != nil {
		return errors.New("AddVirtualHostAtHttpdConfFileError: " + err.Error())
	}

	listenerMapRegex := `^[[:space:]]*map[[:space:]]\+[[:alnum:].-]\+[[:space:]]\+\*`
	newListenerMapLine := "\\ \\ map                     " + hostnameStr + " " + hostnameStr
	_, err = infraHelper.RunCmd(
		"sed", "-ie", "/"+listenerMapRegex+"/a"+newListenerMapLine,
		infraData.GlobalConfigs.OlsHttpdConfFilePath,
	)
	if err != nil {
		return errors.New("UpdateListenerMapLineError: " + err.Error())
	}

	return nil
}
