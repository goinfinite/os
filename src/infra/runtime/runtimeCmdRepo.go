package runtimeInfra

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	domainRepository "github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

const phpDirectiveKeywordPattern = `php_(?:value|flag)`

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
	phpVersionEntity, err := repo.runtimeQueryRepo.ReadPhpVersion(
		runRequest.Hostname,
	)
	if err != nil {
		return runResponse, err
	}
	phpVersionWithoutDots := phpVersionEntity.Value.RemoveDots()
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
		createErr := repo.fileClerk.CreateDir(workingDir)
		if createErr != nil {
			return runResponse, errors.New(
				"PhpWorkingDirCreateFailed: " + createErr.Error(),
			)
		}
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

	if cmdErr != nil {
		shellError, assertOk := cmdErr.(*tkInfra.ShellError)
		if !assertOk {
			return runResponse, errors.New("RunPhpCommandFailed: " + cmdErr.Error())
		}

		stdError, err := valueObject.NewUnixCommandOutput(shellError.StdErr)
		if err != nil {
			return runResponse, err
		}

		runResponse.StdOutput = &stdOutput
		runResponse.StdError = &stdError
		runResponse.ExitCode = &shellError.ExitCode
		return runResponse, nil
	}

	successExitCode := 0
	return dto.RunPhpCommandResponse{
		StdOutput: &stdOutput,
		StdError:  nil,
		ExitCode:  &successExitCode,
	}, nil
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

func (repo *RuntimeCmdRepo) restartPhpWebServer() error {
	err := repo.validatePhpWebServerConfig()
	if err != nil {
		return err
	}

	servicesCmdRepo := servicesInfra.NewServicesCmdRepo(repo.persistentDbSvc)
	err = servicesCmdRepo.Restart(valueObject.ServiceNamePhpWebServer)
	if err != nil {
		return errors.New("RestartPhpWebServerFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) resolvePhpVirtualHostTrustedOwners() (
	trustedDirOwnerUsernames []tkValueObject.UnixUsername,
) {
	return []tkValueObject.UnixUsername{
		tkValueObject.UnixUsername(infraEnvs.PhpWebServerUsername),
		tkValueObject.UnixUsername(infraEnvs.PhpWebServerConfOwnerUsername),
	}
}

func (repo *RuntimeCmdRepo) resolvePhpWebServerTrustedOwners() (
	trustedDirOwnerUsernames []tkValueObject.UnixUsername,
) {
	return []tkValueObject.UnixUsername{
		tkValueObject.UnixUsername(infraEnvs.PhpWebServerConfOwnerUsername),
	}
}

func (repo *RuntimeCmdRepo) replaceFileContentByRegex(
	filePath tkValueObject.UnixAbsoluteFilePath,
	trustedDirOwnerUsernames []tkValueObject.UnixUsername,
	regexPattern *regexp.Regexp,
	replacement string,
) (replacementCount int, err error) {
	return repo.fileClerk.FileContentRegexReplace(
		tkInfra.FileRegexReplaceSettings{
			FilePath:                 filePath,
			SymlinkPolicy:            &tkInfra.FileClerkSymlinkPolicyResolve,
			TrustedDirOwnerUsernames: trustedDirOwnerUsernames,
		},
		regexPattern, replacement,
	)
}

func (repo *RuntimeCmdRepo) listenerMapLineRegexFactory(
	vhostName string,
) (*regexp.Regexp, error) {
	quotedVhostName := regexp.QuoteMeta(vhostName)
	mapLineRegex, err := regexp.Compile(
		"(?m)^([[:space:]]*map[[:space:]]+)" + quotedVhostName + "[[:space:]]+.*$",
	)
	if err != nil {
		return nil, errors.New(
			"PhpListenerMapLineRegexCompileFailed: " + err.Error(),
		)
	}

	return mapLineRegex, nil
}

// The wildcard form makes this vhost the catch-all for the host and its
// subdomains. A subdomain with its own mapping takes precedence.
func (repo *RuntimeCmdRepo) listenerMapDomainsFactory(vhostName string) string {
	return vhostName + " " + vhostName + ", *." + vhostName
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
		if errors.Is(err, domainRepository.ErrPhpVirtualHostNotFound) {
			slog.Debug(
				"SkippingUpdatePhpVirtualHost",
				slog.String("reason", "PhpVirtualHostNotFound"),
			)
			return nil
		}
		return errors.New("PhpConfFilePathResolutionFailed: " + err.Error())
	}
	quotedPreviousHostname := regexp.QuoteMeta(previousHostname.String())
	quotedPkiConfDir := regexp.QuoteMeta(infraEnvs.PkiConfDir)
	newHostnameStr := newHostname.String()

	hostnameSubstitutionRegex := regexp.MustCompile(
		"(?m)(^|[^[:alnum:].-])" + quotedPreviousHostname +
			"([^[:alnum:].-]|$)",
	)
	hostnameSubstitutionReplacement := "${1}" + newHostnameStr + "${2}"

	listenerMapSubstitutionRegex, err := repo.listenerMapLineRegexFactory(
		previousHostname.String(),
	)
	if err != nil {
		return err
	}
	listenerMapSubstitutionReplacement := "${1}" +
		repo.listenerMapDomainsFactory(newHostnameStr)

	sslFilePathSubstitutionRegex := regexp.MustCompile(
		"(?m)(keyFile|certFile)([[:space:]]+)" + quotedPkiConfDir + "/" +
			quotedPreviousHostname + `\.(key|crt)([[:space:]]|$)`,
	)
	sslFilePathSubstitutionReplacement := "${1}${2}" + infraEnvs.PkiConfDir +
		"/" + newHostnameStr + ".${3}${4}"

	_, err = repo.replaceFileContentByRegex(
		phpConfFilePath, repo.resolvePhpVirtualHostTrustedOwners(),
		hostnameSubstitutionRegex, hostnameSubstitutionReplacement,
	)
	if err != nil {
		return errors.New("PhpConfHostnameSubstitutionFailed: " + err.Error())
	}

	phpWebServerMainConfFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.PhpWebServerMainConfFilePath, false,
	)
	if err != nil {
		return errors.New("DefinePhpWebServerMainConfFilePathError: " + err.Error())
	}
	phpWebServerTrustedOwners := repo.resolvePhpWebServerTrustedOwners()

	_, err = repo.replaceFileContentByRegex(
		phpWebServerMainConfFilePath, phpWebServerTrustedOwners,
		listenerMapSubstitutionRegex, listenerMapSubstitutionReplacement,
	)
	if err != nil {
		return errors.New(
			"HttpdListenerMapHostnameSubstitutionFailed: " + err.Error(),
		)
	}

	_, err = repo.replaceFileContentByRegex(
		phpWebServerMainConfFilePath, phpWebServerTrustedOwners,
		hostnameSubstitutionRegex, hostnameSubstitutionReplacement,
	)
	if err != nil {
		return errors.New("HttpdConfigHostnameSubstitutionFailed: " + err.Error())
	}

	_, err = repo.replaceFileContentByRegex(
		phpWebServerMainConfFilePath, phpWebServerTrustedOwners,
		sslFilePathSubstitutionRegex, sslFilePathSubstitutionReplacement,
	)
	if err != nil {
		return errors.New("HttpdConfigSslFilePathSubstitutionFailed: " + err.Error())
	}

	return repo.restartPhpWebServer()
}

func (repo *RuntimeCmdRepo) UpdatePhpVersion(
	hostname tkValueObject.Fqdn,
	version valueObject.PhpVersion,
) error {
	phpVersionEntity, err := repo.runtimeQueryRepo.ReadPhpVersion(hostname)
	if err != nil {
		return err
	}

	if phpVersionEntity.Value == version {
		return nil
	}

	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(
		hostname,
	)
	if err != nil {
		return err
	}

	newLsapiLine := "lsapi:lsphp" + version.RemoveDots()
	lsapiLineRegex := regexp.MustCompile(`lsapi:lsphp[0-9][0-9]\b`)
	_, err = repo.replaceFileContentByRegex(
		phpConfFilePath, repo.resolvePhpVirtualHostTrustedOwners(),
		lsapiLineRegex, newLsapiLine,
	)
	if err != nil {
		return errors.New("UpdatePhpVersionFailed: " + err.Error())
	}

	isPrimaryVirtualHost := repo.vhostHelpers.IsPrimaryVirtualHost(hostname)
	if isPrimaryVirtualHost {
		sourcePhpCliPath := "/usr/local/lsws/lsphp" +
			version.RemoveDots() + "/bin/php"
		updatePhpCliVersionCmd := "unlink /usr/bin/php; ln -s " +
			sourcePhpCliPath + " /usr/bin/php"
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

func (repo *RuntimeCmdRepo) phpSettingLineRegexFactory(
	settingName string,
) (*regexp.Regexp, error) {
	searchPattern := `(?m)^([ \t]*` + phpDirectiveKeywordPattern + `[ \t]+` +
		regexp.QuoteMeta(settingName) + `[ \t]+)(.*)$`

	return regexp.Compile(searchPattern)
}

func (repo *RuntimeCmdRepo) phpSettingLinesMatchDesiredValue(
	settingLineRegex *regexp.Regexp,
	desiredValue, confContent string,
) bool {
	matchedSettingLines := settingLineRegex.FindAllStringSubmatch(
		confContent, -1,
	)
	if len(matchedSettingLines) == 0 {
		return false
	}

	for _, matchedSettingLine := range matchedSettingLines {
		currentValue := matchedSettingLine[2]
		if currentValue != desiredValue {
			return false
		}
	}

	return true
}

func (repo *RuntimeCmdRepo) phpSettingsConfFileUpdater(
	phpConfFilePath tkValueObject.UnixAbsoluteFilePath,
	settingsEntities []entity.PhpSetting,
) (hasUpdatedSetting bool, err error) {
	confFilePathStr := phpConfFilePath.String()

	confContent, err := repo.fileClerk.ReadFileContent(confFilePathStr, nil)
	if err != nil {
		return false, err
	}

	hasWriteFailure := false
	trustedOwners := repo.resolvePhpVirtualHostTrustedOwners()
	for _, settingEntity := range settingsEntities {
		settingName := settingEntity.Name.String()
		desiredValue := settingEntity.Value.String()
		if settingEntity.Value.ReadType() == valueObject.PhpSettingValueTypeString {
			desiredValue = "\"" + desiredValue + "\""
		}

		settingLineRegex, regexErr := repo.phpSettingLineRegexFactory(settingName)
		if regexErr != nil {
			slog.Error(
				"PhpSettingLineRegexCompileFailed",
				slog.String("settingName", settingName),
				slog.String("err", regexErr.Error()),
			)
			hasWriteFailure = true
			continue
		}

		if repo.phpSettingLinesMatchDesiredValue(
			settingLineRegex, desiredValue, confContent,
		) {
			continue
		}

		replacement := "${1}" + strings.ReplaceAll(desiredValue, "$", "$$")
		replaceCount, replaceErr := repo.replaceFileContentByRegex(
			phpConfFilePath, trustedOwners, settingLineRegex, replacement,
		)
		if replaceErr != nil {
			slog.Error(
				"UpdatePhpSettingFailed",
				slog.String("settingName", settingName),
				slog.String("settingValue", desiredValue),
				slog.String("err", replaceErr.Error()),
			)
			hasWriteFailure = true
			continue
		}

		if replaceCount == 0 {
			slog.Error(
				"PhpSettingDirectiveNotFound",
				slog.String("settingName", settingName),
				slog.String("confFilePath", confFilePathStr),
			)
			hasWriteFailure = true
			continue
		}

		refreshedContent, refreshErr := repo.fileClerk.ReadFileContent(
			confFilePathStr, nil,
		)
		if refreshErr != nil {
			slog.Error(
				"RefreshPhpConfContentFailed",
				slog.String("confFilePath", confFilePathStr),
				slog.String("err", refreshErr.Error()),
			)
		}
		if refreshErr == nil {
			confContent = refreshedContent
		}

		hasUpdatedSetting = true
	}

	if hasWriteFailure && !hasUpdatedSetting {
		return false, errors.New("PhpSettingsUpdateFailed: NoSettingWritten")
	}

	return hasUpdatedSetting, nil
}

func (repo *RuntimeCmdRepo) UpdatePhpSettings(
	hostname tkValueObject.Fqdn,
	settingsEntities []entity.PhpSetting,
) error {
	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(
		hostname,
	)
	if err != nil {
		return err
	}

	hasUpdatedSetting, err := repo.phpSettingsConfFileUpdater(
		phpConfFilePath, settingsEntities,
	)
	if err != nil {
		return err
	}
	if !hasUpdatedSetting {
		return nil
	}

	return repo.restartPhpWebServer()
}

func (repo *RuntimeCmdRepo) phpExtensionPackageNameResolver(
	phpVersion valueObject.PhpVersion,
	moduleName string,
) string {
	lsphpPackagePrefix := "lsphp" + phpVersion.RemoveDots() + "-"
	switch moduleName {
	case "mysqli", "pdo_mysql":
		return lsphpPackagePrefix + "mysql"
	case "pdo_sqlite", "sqlite3":
		return lsphpPackagePrefix + "sqlite3"
	case "pdo_dblib":
		return lsphpPackagePrefix + "sybase"
	default:
		return lsphpPackagePrefix + moduleName
	}
}

func (repo *RuntimeCmdRepo) isPhpModuleIniFilePresent(
	filePath string,
) (bool, error) {
	fileInfo, err := os.Lstat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return !fileInfo.IsDir(), nil
}

func (repo *RuntimeCmdRepo) enablePhpModule(
	phpVersion valueObject.PhpVersion,
	moduleEntity entity.PhpModule,
) error {
	moduleNameStr := moduleEntity.Name.String()
	lsphpDir := "/usr/local/lsws/lsphp" + phpVersion.RemoveDots()
	iniRootDir := lsphpDir + "/etc/php/" + phpVersion.String()
	modsAvailableDir := iniRootDir + "/mods-available"
	modsDisabledDir := iniRootDir + "/mods-disabled"

	disabledIniFile := filepath.Join(modsDisabledDir, moduleNameStr+".ini")
	isDisabledIniFilePresent, err := repo.isPhpModuleIniFilePresent(
		disabledIniFile,
	)
	if err != nil {
		return errors.New("ReadDisabledPhpModuleIniFileFailed: " + err.Error())
	}
	if isDisabledIniFilePresent {
		enabledIniFile := filepath.Join(modsAvailableDir, moduleNameStr+".ini")
		err = os.Rename(disabledIniFile, enabledIniFile)
		if err != nil {
			return errors.New("EnablePhpModuleFailed: " + err.Error())
		}

		return nil
	}

	err = infraHelper.InstallPkgs([]string{
		repo.phpExtensionPackageNameResolver(phpVersion, moduleNameStr),
	})
	if err != nil {
		return errors.New("InstallPhpModulePackageFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) disablePhpModule(
	phpVersion valueObject.PhpVersion,
	moduleEntity entity.PhpModule,
) error {
	moduleNameStr := moduleEntity.Name.String()
	if repo.runtimeQueryRepo.isPhpToolModule(moduleNameStr) {
		err := infraHelper.RemovePkgs([]string{
			repo.phpExtensionPackageNameResolver(phpVersion, moduleNameStr),
		})
		if err != nil {
			return errors.New("RemovePhpModulePackageFailed: " + err.Error())
		}

		return nil
	}

	iniRootDir := "/usr/local/lsws/lsphp" +
		phpVersion.RemoveDots() + "/etc/php/" + phpVersion.String()
	modsAvailableDir := iniRootDir + "/mods-available"
	modsDisabledDir := iniRootDir + "/mods-disabled"

	enabledIniFile := filepath.Join(modsAvailableDir, moduleNameStr+".ini")
	isEnabledIniFilePresent, err := repo.isPhpModuleIniFilePresent(
		enabledIniFile,
	)
	if err != nil {
		return errors.New("ReadEnabledPhpModuleIniFileFailed: " + err.Error())
	}
	if !isEnabledIniFilePresent {
		return errors.New("PhpModuleIniFileNotFound: " + enabledIniFile)
	}

	disabledIniFile := filepath.Join(modsDisabledDir, moduleNameStr+".ini")
	err = os.MkdirAll(modsDisabledDir, 0755)
	if err != nil {
		return errors.New("CreatePhpModulesDisabledDirFailed: " + err.Error())
	}

	err = os.Rename(enabledIniFile, disabledIniFile)
	if err != nil {
		return errors.New("DisablePhpModuleFailed: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) updatePhpModule(
	phpVersion valueObject.PhpVersion,
	moduleEntity entity.PhpModule,
) error {
	if moduleEntity.Status {
		return repo.enablePhpModule(phpVersion, moduleEntity)
	}

	return repo.disablePhpModule(phpVersion, moduleEntity)
}

func (repo *RuntimeCmdRepo) applyPhpModuleStatusUpdates(
	phpVersion valueObject.PhpVersion,
	modulesEntities []entity.PhpModule,
) (anyPhpModuleUpdated bool, moduleErrors map[string]error) {
	moduleErrors = map[string]error{}

	for _, moduleEntity := range modulesEntities {
		moduleError := repo.updatePhpModule(phpVersion, moduleEntity)
		if moduleError != nil {
			slog.Error(
				"UpdatePhpModuleFailed",
				slog.String("moduleName", moduleEntity.Name.String()),
				slog.Bool("desiredStatus", moduleEntity.Status),
				slog.String("err", moduleError.Error()),
			)
			moduleErrors[moduleEntity.Name.String()] = moduleError
			continue
		}

		anyPhpModuleUpdated = true
	}

	return anyPhpModuleUpdated, moduleErrors
}

func (repo *RuntimeCmdRepo) readActivePhpModuleNames(
	phpVersion valueObject.PhpVersion,
) (map[string]any, error) {
	allModulesEntities, err := repo.runtimeQueryRepo.ReadPhpModules(phpVersion)
	if err != nil {
		return nil, err
	}

	activeModuleNames := map[string]any{}
	for _, moduleEntity := range allModulesEntities {
		if !moduleEntity.Status {
			continue
		}

		activeModuleNames[moduleEntity.Name.String()] = nil
	}

	return activeModuleNames, nil
}

func (repo *RuntimeCmdRepo) phpModuleStatusDoubleChecker(
	modulesEntities []entity.PhpModule,
	activeModuleNames map[string]any,
	updateErrors map[string]error,
) map[string]error {
	failedModulesWithReason := map[string]error{}

	for _, moduleEntity := range modulesEntities {
		moduleName := moduleEntity.Name.String()
		updateError, hasUpdateError := updateErrors[moduleName]
		if hasUpdateError {
			failedModulesWithReason[moduleName] = updateError
			continue
		}

		_, isModuleActive := activeModuleNames[moduleName]
		if isModuleActive == moduleEntity.Status {
			continue
		}

		slog.Error(
			"PhpModuleStatusMismatch",
			slog.String("moduleName", moduleName),
			slog.Bool("expectedStatus", moduleEntity.Status),
			slog.Bool("actualStatus", isModuleActive),
		)
		failedModulesWithReason[moduleName] = errors.New("PhpModuleStatusMismatch")
	}

	return failedModulesWithReason
}

func (repo *RuntimeCmdRepo) phpModuleUpdateResponseFactory(
	modulesEntities []entity.PhpModule,
	moduleErrors map[string]error,
) dto.UpdatePhpModulesResponse {
	modulesSuccessfullyUpdated := []dto.PhpModuleUpdate{}
	failedModulesWithReason := []dto.PhpModuleUpdateFailure{}

	for _, moduleEntity := range modulesEntities {
		moduleError, hasModuleError := moduleErrors[moduleEntity.Name.String()]
		if hasModuleError {
			failedModulesWithReason = append(
				failedModulesWithReason,
				dto.NewPhpModuleUpdateFailure(
					moduleEntity.Name, moduleEntity.Status,
					valueObject.NewFailureReason(moduleError.Error()),
				),
			)
			continue
		}

		modulesSuccessfullyUpdated = append(
			modulesSuccessfullyUpdated,
			dto.NewPhpModuleUpdate(moduleEntity.Name, moduleEntity.Status),
		)
	}

	return dto.NewUpdatePhpModulesResponse(
		modulesSuccessfullyUpdated, failedModulesWithReason,
	)
}

func (repo *RuntimeCmdRepo) UpdatePhpModules(
	hostname tkValueObject.Fqdn,
	expectedPhpVersion valueObject.PhpVersion,
	modulesToUpdateEntities []entity.PhpModule,
) (responseDto dto.UpdatePhpModulesResponse, err error) {
	phpVersionEntity, err := repo.runtimeQueryRepo.ReadPhpVersion(hostname)
	if err != nil {
		return responseDto, err
	}
	if phpVersionEntity.Value != expectedPhpVersion {
		return responseDto, domainRepository.ErrPhpVersionChanged
	}

	anyPhpModuleUpdated, moduleErrors := repo.applyPhpModuleStatusUpdates(
		phpVersionEntity.Value, modulesToUpdateEntities,
	)

	if anyPhpModuleUpdated {
		err = repo.restartPhpWebServer()
		if err != nil {
			return responseDto, err
		}
	}

	activeModuleNames, err := repo.readActivePhpModuleNames(
		phpVersionEntity.Value,
	)
	if err != nil {
		return responseDto, errors.New(
			"ReadActivePhpModulesFailed: " + err.Error(),
		)
	}

	failedModuleErrors := repo.phpModuleStatusDoubleChecker(
		modulesToUpdateEntities, activeModuleNames, moduleErrors,
	)

	return repo.phpModuleUpdateResponseFactory(
		modulesToUpdateEntities, failedModuleErrors,
	), nil
}

func (repo *RuntimeCmdRepo) removeListenerMapLines(
	hostname tkValueObject.Fqdn,
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
) error {
	hostnameStr := strings.ReplaceAll(hostname.String(), "*.", "")
	mapLineRegex, err := repo.listenerMapLineRegexFactory(hostnameStr)
	if err != nil {
		return err
	}
	listenerMapLineRemovalRegex, err := regexp.Compile(
		mapLineRegex.String() + `\n?`,
	)
	if err != nil {
		return errors.New(
			"PhpListenerMapLineRemovalRegexCompileFailed: " + err.Error(),
		)
	}
	_, err = repo.replaceFileContentByRegex(
		phpWebServerMainConfFilePath, repo.resolvePhpWebServerTrustedOwners(),
		listenerMapLineRemovalRegex, "",
	)
	if err != nil {
		return errors.New("RemoveListenerMapLineError: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) countListenerMapLines(
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
	vhostName string,
) (int, error) {
	mapLineRegex, err := repo.listenerMapLineRegexFactory(vhostName)
	if err != nil {
		return 0, err
	}

	findings, err := repo.fileClerk.FileContentRegexSearch(
		phpWebServerMainConfFilePath, mapLineRegex,
	)
	if err != nil {
		return 0, errors.New("ReadPhpListenerMapLinesError: " + err.Error())
	}

	return len(findings), nil
}

func (repo *RuntimeCmdRepo) readListenerMapLineCounts(
	hostname tkValueObject.Fqdn,
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
) (
	primaryHostname tkValueObject.Fqdn,
	primaryMapLineCount, hostnameMapLineCount int,
	err error,
) {
	primaryHostname, err = repo.vhostHelpers.ReadPrimaryVirtualHostHostname()
	if err != nil {
		return primaryHostname, 0, 0, errors.New(
			"ReadPrimaryVirtualHostHostnameError: " + err.Error(),
		)
	}

	primaryMapLineCount, err = repo.countListenerMapLines(
		phpWebServerMainConfFilePath, primaryHostname.String(),
	)
	if err != nil {
		return primaryHostname, 0, 0, err
	}

	hostnameStr := strings.ReplaceAll(hostname.String(), "*.", "")
	hostnameMapLineCount, err = repo.countListenerMapLines(
		phpWebServerMainConfFilePath, hostnameStr,
	)
	if err != nil {
		return primaryHostname, 0, 0, err
	}

	return primaryHostname, primaryMapLineCount, hostnameMapLineCount, nil
}

func (repo *RuntimeCmdRepo) mapVirtualHostOnEveryListener(
	hostname tkValueObject.Fqdn,
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
) error {
	primaryHostname, primaryMapLineCount, hostnameMapLineCount, err := repo.
		readListenerMapLineCounts(hostname, phpWebServerMainConfFilePath)
	if err != nil {
		return err
	}
	if primaryMapLineCount == 0 {
		return errors.New("PrimaryListenerMapLineNotFound")
	}
	if hostnameMapLineCount == primaryMapLineCount {
		return nil
	}

	err = repo.removeListenerMapLines(
		hostname, phpWebServerMainConfFilePath,
	)
	if err != nil {
		return err
	}

	primaryMapLineRegex, err := repo.listenerMapLineRegexFactory(
		primaryHostname.String(),
	)
	if err != nil {
		return err
	}

	hostnameStr := strings.ReplaceAll(hostname.String(), "*.", "")
	newListenerMapLine := "  map                     " +
		repo.listenerMapDomainsFactory(hostnameStr)
	_, err = repo.replaceFileContentByRegex(
		phpWebServerMainConfFilePath, repo.resolvePhpWebServerTrustedOwners(),
		primaryMapLineRegex, "${0}\n"+newListenerMapLine,
	)
	if err != nil {
		return errors.New("UpdateListenerMapLineError: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) listenerMapLinesDoubleChecker(
	hostname tkValueObject.Fqdn,
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
) error {
	_, primaryMapLineCount, hostnameMapLineCount, err := repo.
		readListenerMapLineCounts(hostname, phpWebServerMainConfFilePath)
	if err != nil {
		return err
	}
	if hostnameMapLineCount != primaryMapLineCount {
		return errors.New("PhpVirtualHostListenerMapNotFound")
	}

	return nil
}

func (repo *RuntimeCmdRepo) virtualHostBlockRegexFactory(
	hostname string,
) (*regexp.Regexp, error) {
	baseHostname := strings.ReplaceAll(hostname, "*.", "")
	quotedHostname := regexp.QuoteMeta(baseHostname)
	vhostBlockRegex, err := regexp.Compile(
		"(?m)^[[:space:]]*virtualhost[[:space:]]+" + quotedHostname +
			"[[:space:]]*\\{",
	)
	if err != nil {
		return nil, errors.New(
			"PhpVirtualHostBlockRegexCompileFailed: " + err.Error(),
		)
	}

	return vhostBlockRegex, nil
}

func (repo *RuntimeCmdRepo) isVirtualHostBlockPresent(
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
	hostname string,
) (bool, error) {
	vhostBlockRegex, err := repo.virtualHostBlockRegexFactory(hostname)
	if err != nil {
		return false, err
	}

	findings, err := repo.fileClerk.FileContentRegexSearch(
		phpWebServerMainConfFilePath, vhostBlockRegex,
	)
	if err != nil {
		return false, errors.New("ReadPhpVirtualHostBlockError: " + err.Error())
	}

	return len(findings) > 0, nil
}

func (repo *RuntimeCmdRepo) insertVirtualHostBlock(
	hostname tkValueObject.Fqdn,
	phpConfFilePath tkValueObject.UnixAbsoluteFilePath,
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
) error {
	vhostBlockIsPresent, err := repo.isVirtualHostBlockPresent(
		phpWebServerMainConfFilePath, hostname.String(),
	)
	if err != nil {
		return err
	}
	if vhostBlockIsPresent {
		return nil
	}

	hostnameStr := strings.ReplaceAll(hostname.String(), "*.", "")
	phpVhostHttpdConf := `
virtualhost ` + hostnameStr + ` {
  vhRoot                  /app/html/` + hostnameStr + `/
  configFile              ` + phpConfFilePath.String() + `
  allowSymbolLink         1
  enableScript            1
  restrained              0
  setUIDMode              0
}
`
	err = repo.fileClerk.AppendFileContent(tkInfra.FileAppendSettings{
		FilePath:                 phpWebServerMainConfFilePath,
		SymlinkPolicy:            &tkInfra.FileClerkSymlinkPolicyResolve,
		TrustedDirOwnerUsernames: repo.resolvePhpWebServerTrustedOwners(),
	}, phpVhostHttpdConf)
	if err != nil {
		return errors.New("AddVirtualHostAtHttpdConfFileError: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) removeVirtualHostBlock(
	hostname tkValueObject.Fqdn,
	phpWebServerMainConfFilePath tkValueObject.UnixAbsoluteFilePath,
) error {
	vhostBlockRegex, err := repo.virtualHostBlockRegexFactory(hostname.String())
	if err != nil {
		return err
	}
	vhostBlockRemovalRegex, err := regexp.Compile(
		vhostBlockRegex.String() + "[^}]*\\}\n?",
	)
	if err != nil {
		return errors.New(
			"PhpVirtualHostBlockRemovalRegexCompileFailed: " + err.Error(),
		)
	}
	_, err = repo.replaceFileContentByRegex(
		phpWebServerMainConfFilePath, repo.resolvePhpWebServerTrustedOwners(),
		vhostBlockRemovalRegex, "",
	)
	if err != nil {
		return errors.New("RemoveVirtualHostBlockError: " + err.Error())
	}

	blockStillPresent, err := repo.isVirtualHostBlockPresent(
		phpWebServerMainConfFilePath, hostname.String(),
	)
	if err != nil {
		return err
	}
	if blockStillPresent {
		return errors.New("RemovePhpVirtualHostBlockFailed")
	}

	return nil
}

func (repo *RuntimeCmdRepo) DeletePhpVirtualHost(hostname tkValueObject.Fqdn) error {
	phpWebServerMainConfFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.PhpWebServerMainConfFilePath, false,
	)
	if err != nil {
		return errors.New("DefinePhpWebServerMainConfFilePathError: " + err.Error())
	}

	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(
		hostname,
	)
	if err != nil && !errors.Is(err, domainRepository.ErrPhpVirtualHostNotFound) {
		return errors.New("ReadPhpVirtualHostConfFilePathError: " + err.Error())
	}

	err = repo.removeListenerMapLines(hostname, phpWebServerMainConfFilePath)
	if err != nil {
		return err
	}

	err = repo.removeVirtualHostBlock(hostname, phpWebServerMainConfFilePath)
	if err != nil {
		return err
	}

	err = repo.fileClerk.DeleteFile(phpConfFilePath.String())
	if err != nil {
		return errors.New("RemovePhpConfFileError: " + err.Error())
	}

	return repo.restartPhpWebServer()
}

func (repo *RuntimeCmdRepo) createVirtualHostConfFile(
	hostname tkValueObject.Fqdn,
	phpConfFilePath tkValueObject.UnixAbsoluteFilePath,
) error {
	hostnameStr := strings.ReplaceAll(hostname.String(), "*.", "")
	templatePhpVhostConfFilePath := infraEnvs.PhpWebServerConfDir + "/template"
	err := repo.fileClerk.CopyFile(
		templatePhpVhostConfFilePath, phpConfFilePath.String(),
	)
	if err != nil {
		return errors.New("CopyPhpConfTemplateError: " + err.Error())
	}

	primaryVirtualHostPlaceholderRegex := regexp.MustCompile(
		regexp.QuoteMeta(infraEnvs.PrimaryVirtualHostPlaceholderHostname),
	)
	_, err = repo.replaceFileContentByRegex(
		phpConfFilePath, repo.resolvePhpVirtualHostTrustedOwners(),
		primaryVirtualHostPlaceholderRegex, hostnameStr,
	)
	if err != nil {
		return errors.New("UpdatePhpVirtualHostConfFileError: " + err.Error())
	}

	return nil
}

func (repo *RuntimeCmdRepo) CreatePhpVirtualHost(hostname tkValueObject.Fqdn) error {
	phpWebServerMainConfFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
		infraEnvs.PhpWebServerMainConfFilePath, false,
	)
	if err != nil {
		return errors.New("DefinePhpWebServerMainConfFilePathError: " + err.Error())
	}

	phpConfFilePath, err := repo.runtimeQueryRepo.ReadPhpVirtualHostConfFilePath(
		hostname,
	)
	phpConfExists := err == nil
	if err != nil && !errors.Is(err, domainRepository.ErrPhpVirtualHostNotFound) {
		return errors.New("ReadPhpVirtualHostConfFilePathError: " + err.Error())
	}

	if !phpConfExists {
		err = repo.createVirtualHostConfFile(hostname, phpConfFilePath)
		if err != nil {
			return err
		}
	}

	err = repo.insertVirtualHostBlock(
		hostname, phpConfFilePath, phpWebServerMainConfFilePath,
	)
	if err != nil {
		return err
	}

	err = repo.mapVirtualHostOnEveryListener(
		hostname, phpWebServerMainConfFilePath,
	)
	if err != nil {
		return err
	}

	err = repo.listenerMapLinesDoubleChecker(
		hostname, phpWebServerMainConfFilePath,
	)
	if err != nil {
		return err
	}

	return repo.restartPhpWebServer()
}
