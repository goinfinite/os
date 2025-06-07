package runtimeInfra

import (
	"errors"
	"log/slog"
	"os"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
)

type RuntimeCmdRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	runtimeQueryRepo RuntimeQueryRepo
}

func NewRuntimeCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimeCmdRepo {
	return &RuntimeCmdRepo{
		persistentDbSvc:  persistentDbSvc,
		runtimeQueryRepo: RuntimeQueryRepo{},
	}
}

func (repo *RuntimeCmdRepo) RunPhpCommand(
	runRequest dto.RunPhpCommandRequest,
) (runResponse dto.RunPhpCommandResponse, err error) {
	phpVersionEntity, err := repo.runtimeQueryRepo.ReadPhpVersion(runRequest.Hostname)
	if err != nil {
		return runResponse, err
	}
	phpVersionWithoutDots := phpVersionEntity.Value.GetWithoutDots()
	if phpVersionWithoutDots == "" {
		return runResponse, errors.New("PhpVersionNotFound")
	}

	phpCli := "/usr/local/lsws/lsphp" + phpVersionWithoutDots + "/bin/php"
	if !infraHelper.FileExists(phpCli) {
		return runResponse, errors.New("PhpCliNotFound")
	}

	timeoutSecs := uint64(600)
	if runRequest.TimeoutSecs != nil {
		timeoutSecs = *runRequest.TimeoutSecs
	}
	workingDir := infraEnvs.PrimaryPublicDir
	if !infraHelper.IsPrimaryVirtualHost(runRequest.Hostname) {
		workingDir += "/" + runRequest.Hostname.String()
	}
	if !infraHelper.FileExists(workingDir) {
		workingDir = infraEnvs.PrimaryPublicDir
	}

	cmdOutput, cmdErr := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               phpCli,
		Args:                  []string{runRequest.Command.String()},
		Username:              infraEnvs.PhpWebserverUsername,
		WorkingDirectory:      workingDir,
		ShouldRunWithSubShell: true,
		ExecutionTimeoutSecs:  timeoutSecs,
	})
	stdOutput, err := valueObject.NewUnixCommandOutput(cmdOutput)
	if err != nil {
		return runResponse, err
	}

	if errorMessage, assertOk := cmdErr.(*infraHelper.CmdError); assertOk {
		stdError, err := valueObject.NewUnixCommandOutput(errorMessage.StdErr)
		if err != nil {
			return runResponse, err
		}

		runResponse.StdOutput = &stdOutput
		runResponse.StdError = &stdError
		runResponse.ExitCode = &errorMessage.ExitCode
		return runResponse, nil
	}

	successExitCode := 0
	return dto.RunPhpCommandResponse{
		StdOutput: &stdOutput,
		StdError:  nil,
		ExitCode:  &successExitCode,
	}, nil
}

func (repo *RuntimeCmdRepo) restartPhpWebserver() error {
	phpSvcName, _ := valueObject.NewServiceName("php-webserver")
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)
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
	updatePhpVersionCmd := "sed -i 's/lsapi:lsphp[0-9][0-9]/" + newLsapiLine +
		"/g' " + phpConfFilePath.String()
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               updatePhpVersionCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("UpdatePhpVersionFailed: " + err.Error())
	}

	isPrimaryVirtualHost := infraHelper.IsPrimaryVirtualHost(hostname)
	if isPrimaryVirtualHost {
		sourcePhpCliPath := "/usr/local/lsws/lsphp" + version.GetWithoutDots() + "/bin/php"
		updatePhpCliVersionCmd := "unlink /usr/bin/php; ln -s " + sourcePhpCliPath + " /usr/bin/php"
		_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command:               updatePhpCliVersionCmd,
			ShouldRunWithSubShell: true,
		})
		if err != nil {
			return errors.New("UpdatePhpCliVersionError: " + err.Error())
		}
	}

	return repo.restartPhpWebserver()
}

func (repo *RuntimeCmdRepo) UpdatePhpSettings(
	hostname valueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	phpConfFilePath, err := repo.runtimeQueryRepo.GetVirtualHostPhpConfFilePath(hostname)
	if err != nil {
		return err
	}
	phpConfigFilePathStr := phpConfFilePath.String()

	for _, setting := range settings {
		settingName := setting.Name.String()
		settingValue := setting.Value.String()
		if setting.Value.GetType() == "string" {
			settingValue = "\"" + settingValue + "\""
			settingValue = strings.Replace(settingValue, "|", "\\|", -1)
		}

		_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command: "sed",
			Args: []string{
				"-i", "s|" + settingName + " .*|" + settingName + " " + settingValue + "|g",
				phpConfigFilePathStr,
			},
		})
		if err != nil {
			slog.Debug(
				"UpdatePhpSettingFailed",
				slog.String("settingName", settingName),
				slog.String("settingValue", settingValue),
				slog.String("err", err.Error()),
			)
			continue
		}
	}

	return repo.restartPhpWebserver()
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

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               "echo | " + lsphpDir + "/bin/pecl install " + moduleNameStr,
		ShouldRunWithSubShell: true,
	})
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

	return repo.restartPhpWebserver()
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
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "sed",
		Args: []string{
			"-ie", "s/" + infraEnvs.DefaultPrimaryVhost + "/" + hostnameStr + "/g", phpConfFilePathStr,
		},
	})
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
		infraEnvs.PhpWebserverMainConfFilePath, phpVhostHttpdConf, shouldOverwrite,
	)
	if err != nil {
		return errors.New("AddVirtualHostAtHttpdConfFileError: " + err.Error())
	}

	listenerMapRegex := `^[[:space:]]*map[[:space:]]\+[[:alnum:].-]\+[[:space:]]\+\*`
	newListenerMapLine := "\\ \\ map                     " + hostnameStr + " " + hostnameStr
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "sed",
		Args: []string{
			"-ie", "/" + listenerMapRegex + "/a" + newListenerMapLine,
			infraEnvs.PhpWebserverMainConfFilePath,
		},
	})
	if err != nil {
		return errors.New("UpdateListenerMapLineError: " + err.Error())
	}

	return nil
}
