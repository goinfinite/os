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
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

var phpWebServerServiceName, phpWebServerServiceNameError = valueObject.NewServiceName(
	"php-webserver",
)

type RuntimeCmdRepo struct {
	persistentDbSvc  *internalDbInfra.PersistentDatabaseService
	runtimeQueryRepo *RuntimeQueryRepo
	fileClerk        tkInfra.FileClerk
	vhostHelpers     *vhostInfra.VirtualHostHelpers
}

func NewRuntimeCmdRepo(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimeCmdRepo {
	return &RuntimeCmdRepo{
		persistentDbSvc:  persistentDbSvc,
		runtimeQueryRepo: NewRuntimeQueryRepo(),
		fileClerk:        tkInfra.FileClerk{},
		vhostHelpers:     vhostInfra.NewVirtualHostHelpers(),
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
	if !repo.fileClerk.FileExists(phpCli) {
		return runResponse, errors.New("PhpCliNotFound")
	}

	timeoutSecs := uint64(600)
	if runRequest.TimeoutSecs != nil {
		timeoutSecs = *runRequest.TimeoutSecs
	}
	workingDir := infraEnvs.PrimaryVirtualHostPublicDir
	if !repo.vhostHelpers.IsPrimaryVirtualHost(runRequest.Hostname) {
		workingDir += "/" + runRequest.Hostname.String()
	}
	if !repo.fileClerk.FileExists(workingDir) {
		workingDir = infraEnvs.PrimaryVirtualHostPublicDir
	}

	cmdOutput, cmdErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:              phpCli,
		Args:                 []string{runRequest.Command.String()},
		Username:             infraEnvs.PhpWebServerUsername,
		WorkingDirectory:     workingDir,
		ShouldUseSubShell:    true,
		ExecutionTimeoutSecs: timeoutSecs,
	}).Run()
	stdOutput, err := valueObject.NewUnixCommandOutput(cmdOutput)
	if err != nil {
		return runResponse, err
	}

	if errorMessage, assertOk := cmdErr.(*tkInfra.ShellError); assertOk {
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

func (repo *RuntimeCmdRepo) restartPhpWebServer() error {
	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)
	err := servicesCmdRepo.Restart(phpWebServerServiceName)
	if err != nil {
		return errors.New("RestartWebServerFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) regexReplaceInFile(
	searchPattern, replacement, filePath string,
) error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-i", "-E", "s#" + searchPattern + "#" + replacement + "#g", filePath,
		},
	}).Run()
	return err
}

func (repo *RuntimeCmdRepo) validatePhpWebServerConfig() error {
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           infraEnvs.PhpWebServerConfigValidationCmd,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("PhpWebServerConfigValidationFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) UpdatePhpVirtualHostHostname(
	previousHostname, newHostname tkValueObject.Fqdn,
	aliasesHostnames []tkValueObject.Fqdn,
) error {
	if previousHostname == newHostname {
		return nil
	}

	pkiConfDir, parseErr := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.PkiConfDir, false,
	)
	if parseErr != nil {
		return errors.New("InvalidPkiConfDir: " + parseErr.Error())
	}

	createCertErr := infraHelper.CreateSelfSignedSsl(
		pkiConfDir, newHostname, aliasesHostnames,
	)
	if createCertErr != nil {
		return errors.New("CreateSelfSignedSslFailed: " + createCertErr.Error())
	}

	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(
		previousHostname,
	)
	if err != nil {
		if errors.Is(err, ErrPhpVirtualHostNotFound) {
			slog.Debug(
				"SkippingUpdatePhpVirtualHost",
				slog.String("reason", "PhpVirtualHostNotFound"),
			)
			return nil
		}
		return errors.New("PhpConfFilePathResolutionFailed: " + err.Error())
	}
	phpConfFilePathStr := phpConfFilePath.String()

	escapedPreviousHostname := strings.ReplaceAll(previousHostname.String(), ".", `\.`)
	newHostnameStr := newHostname.String()

	hostnameSubstitutionPattern := "(^|[^[:alnum:].-])" +
		escapedPreviousHostname + "([^[:alnum:].-]|$)"
	hostnameSubstitutionReplacement := `\1` + newHostnameStr + `\2`

	listenerMapSubstitutionPattern := "(map[[:space:]]+)" +
		escapedPreviousHostname + "([[:space:]]+)" + escapedPreviousHostname
	listenerMapSubstitutionReplacement := `\1` + newHostnameStr + `\2` +
		newHostnameStr

	sslFilePathSubstitutionPattern := "(keyFile|certFile)[[:space:]]+" +
		infraEnvs.PkiConfDir + "/" + escapedPreviousHostname +
		`\.(key|crt)([[:space:]]|$)`
	sslFilePathSubstitutionReplacement := `\1 ` + infraEnvs.PkiConfDir + "/" +
		newHostnameStr + `.\2\3`

	err = repo.regexReplaceInFile(
		hostnameSubstitutionPattern, hostnameSubstitutionReplacement,
		phpConfFilePathStr,
	)
	if err != nil {
		return errors.New("PhpConfHostnameSubstitutionFailed: " + err.Error())
	}

	httpdConfigFilePath := infraEnvs.PhpWebServerMainConfFilePath
	err = repo.regexReplaceInFile(
		listenerMapSubstitutionPattern, listenerMapSubstitutionReplacement,
		httpdConfigFilePath,
	)
	if err != nil {
		return errors.New("HttpdListenerMapHostnameSubstitutionFailed: " + err.Error())
	}

	err = repo.regexReplaceInFile(
		hostnameSubstitutionPattern, hostnameSubstitutionReplacement,
		httpdConfigFilePath,
	)
	if err != nil {
		return errors.New("HttpdConfigHostnameSubstitutionFailed: " + err.Error())
	}

	err = repo.regexReplaceInFile(
		sslFilePathSubstitutionPattern, sslFilePathSubstitutionReplacement,
		httpdConfigFilePath,
	)
	if err != nil {
		return errors.New("HttpdConfigSslFilePathSubstitutionFailed: " + err.Error())
	}

	err = repo.validatePhpWebServerConfig()
	if err != nil {
		return err
	}

	if phpWebServerServiceNameError != nil {
		return errors.New(
			"PhpWebServerServiceNameResolutionFailed: " +
				phpWebServerServiceNameError.Error(),
		)
	}

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)
	err = servicesCmdRepo.Restart(phpWebServerServiceName)
	if err != nil {
		return errors.New("PhpWebServerRestartFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) UpdatePhpVersion(
	hostname tkValueObject.Fqdn,
	version valueObject.PhpVersion,
) error {
	phpVersion, err := repo.runtimeQueryRepo.ReadPhpVersion(hostname)
	if err != nil {
		return err
	}

	if phpVersion.Value == version {
		return nil
	}

	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(hostname)
	if err != nil {
		return err
	}

	newLsapiLine := "lsapi:lsphp" + version.GetWithoutDots()
	err = repo.regexReplaceInFile(
		"lsapi:lsphp[0-9][0-9]", newLsapiLine, phpConfFilePath.String(),
	)
	if err != nil {
		return errors.New("UpdatePhpVersionFailed: " + err.Error())
	}

	isPrimaryVirtualHost := repo.vhostHelpers.IsPrimaryVirtualHost(hostname)
	if isPrimaryVirtualHost {
		sourcePhpCliPath := "/usr/local/lsws/lsphp" + version.GetWithoutDots() + "/bin/php"
		updatePhpCliVersionCmd := "unlink /usr/bin/php; ln -s " + sourcePhpCliPath + " /usr/bin/php"
		_, err = tkInfra.NewShell(tkInfra.ShellSettings{
			Command:           updatePhpCliVersionCmd,
			ShouldUseSubShell: true,
		}).Run()
		if err != nil {
			return errors.New("UpdatePhpCliVersionError: " + err.Error())
		}
	}

	return repo.restartPhpWebServer()
}

func (repo *RuntimeCmdRepo) UpdatePhpSettings(
	hostname tkValueObject.Fqdn,
	settings []entity.PhpSetting,
) error {
	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(hostname)
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
			settingValue = strings.Replace(settingValue, "/", "\\/", -1)
			settingValue = strings.Replace(settingValue, "#", "\\#", -1)
		}

		phpSettingLinePattern := settingName + " .*"
		phpSettingLineReplacement := settingName + " " + settingValue

		err := repo.regexReplaceInFile(
			phpSettingLinePattern, phpSettingLineReplacement, phpConfigFilePathStr,
		)
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

	return repo.restartPhpWebServer()
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

	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           "echo | " + lsphpDir + "/bin/pecl install " + moduleNameStr,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("InstallPeclModuleFailed: " + err.Error())
	}

	moduleConfigFilePath := modsAvailableDir + "/" + moduleNameStr + ".ini"
	moduleConfigFileContent := "extension=" + moduleNameStr + ".so"
	err = repo.fileClerk.UpdateFileContent(
		moduleConfigFilePath, moduleConfigFileContent, true,
	)
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
	hostname tkValueObject.Fqdn,
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

	return repo.restartPhpWebServer()
}

func (repo *RuntimeCmdRepo) CreatePhpVirtualHost(hostname tkValueObject.Fqdn) error {
	vhostExists := true

	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(hostname)
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
	err = repo.fileClerk.CopyFile(templatePhpVhostConfFilePath, phpConfFilePathStr)
	if err != nil {
		return errors.New("CopyPhpConfTemplateError: " + err.Error())
	}

	hostnameStr := hostname.String()
	err = repo.regexReplaceInFile(
		infraEnvs.DefaultPrimaryVhost, hostnameStr, phpConfFilePathStr,
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
	err = repo.fileClerk.UpdateFileContent(
		infraEnvs.PhpWebServerMainConfFilePath, phpVhostHttpdConf, shouldOverwrite,
	)
	if err != nil {
		return errors.New("AddVirtualHostAtHttpdConfFileError: " + err.Error())
	}

	listenerMapRegex := `^[[:space:]]*map[[:space:]]\+[[:alnum:].-]\+[[:space:]]\+\*`
	newListenerMapLine := "\\ \\ map                     " + hostnameStr + " " + hostnameStr
	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "sed",
		Args: []string{
			"-ie", "/" + listenerMapRegex + "/a" + newListenerMapLine,
			infraEnvs.PhpWebServerMainConfFilePath,
		},
	}).Run()
	if err != nil {
		return errors.New("UpdateListenerMapLineError: " + err.Error())
	}

	return nil
}
