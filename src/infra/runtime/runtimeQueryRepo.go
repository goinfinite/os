package runtimeInfra

import (
	"encoding/json"
	"errors"
	"log/slog"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/goinfinite/os/src/domain/entity"
	domainRepository "github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type RuntimeQueryRepo struct {
	fileClerk tkInfra.FileClerk
}

func NewRuntimeQueryRepo() *RuntimeQueryRepo {
	return &RuntimeQueryRepo{fileClerk: tkInfra.FileClerk{}}
}

func (repo RuntimeQueryRepo) ReadPhpVirtualHostConfFilePath(
	vhostHostname tkValueObject.Fqdn,
) (phpVirtualHostConfFilePath tkValueObject.UnixAbsoluteFilePath, err error) {
	nonWildcardHostname := strings.ReplaceAll(vhostHostname.String(), "*.", "")
	rawPhpVirtualHostConfFilePath := infraEnvs.PhpWebServerConfDir + "/" +
		nonWildcardHostname + ".conf"

	primaryVirtualHostHostname, err := vhostInfra.NewVirtualHostHelpers().
		ReadPrimaryVirtualHostHostnameFromWebServerConf()
	if err != nil {
		return phpVirtualHostConfFilePath, errors.New(
			"WebServerPrimaryVirtualHostNotFound: " + err.Error(),
		)
	}

	primaryVirtualHostPhpConfFilePathStr := infraEnvs.PhpWebServerConfDir +
		"/primary.conf"
	if vhostHostname == primaryVirtualHostHostname {
		rawPhpVirtualHostConfFilePath = primaryVirtualHostPhpConfFilePathStr
	}

	phpVirtualHostConfFilePath, err = tkValueObject.NewUnixAbsoluteFilePath(
		rawPhpVirtualHostConfFilePath, false,
	)
	if err != nil {
		return phpVirtualHostConfFilePath, errors.New(
			"InvalidPhpVirtualHostConfFilePath: " + err.Error(),
		)
	}

	if !repo.fileClerk.FileExists(phpVirtualHostConfFilePath.String()) {
		return phpVirtualHostConfFilePath, domainRepository.ErrPhpVirtualHostNotFound
	}

	return phpVirtualHostConfFilePath, nil
}

func (repo RuntimeQueryRepo) ReadPhpVersionsInstalled() (
	phpVersions []valueObject.PhpVersion, err error,
) {
	output, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: infraEnvs.AwkBinaryPath,
		Args: []string{
			"/extprocessor lsphp/{print $2}", infraEnvs.PhpWebServerMainConfFilePath,
		},
	}).Run()
	if err != nil {
		return phpVersions, errors.New(
			"ReadPhpVersionsInstalledFailed: " + err.Error(),
		)
	}

	for version := range strings.SplitSeq(output, "\n") {
		if version == "" {
			continue
		}

		version = strings.Replace(version, "lsphp", "", 1)
		phpVersion, err := valueObject.NewPhpVersion(version)
		if err != nil {
			slog.Debug(
				"SkippingInvalidPhpVersion",
				slog.String("phpVersion", version),
			)
			continue
		}

		phpVersions = append(phpVersions, phpVersion)
	}

	return phpVersions, nil
}

func (repo RuntimeQueryRepo) ReadPhpVersion(
	hostname tkValueObject.Fqdn,
) (phpVersionEntity entity.PhpVersion, err error) {
	phpVirtualHostConfFilePath, err := repo.ReadPhpVirtualHostConfFilePath(hostname)
	if err != nil {
		return phpVersionEntity, err
	}

	currentPhpVersionStr, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: infraEnvs.AwkBinaryPath,
		Args: []string{
			"/lsapi:lsphp/ {gsub(/[^0-9]/, \"\", $2); print $2}",
			phpVirtualHostConfFilePath.String(),
		},
	}).Run()
	if err != nil {
		return phpVersionEntity, errors.New(
			"ReadCurrentPhpVersionFromFileFailed: " + err.Error(),
		)
	}

	currentPhpVersion, err := valueObject.NewPhpVersion(currentPhpVersionStr)
	if err != nil {
		return phpVersionEntity, errors.New("PhpVersionUnknown: " + err.Error())
	}

	phpVersions, err := repo.ReadPhpVersionsInstalled()
	if err != nil {
		return phpVersionEntity, err
	}

	phpVersionEntity = entity.NewPhpVersion(currentPhpVersion, phpVersions)
	return phpVersionEntity, nil
}

func (repo RuntimeQueryRepo) readPhpTimezones() (timezones []string, err error) {
	timezonesRaw, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "php",
		Args: []string{
			"-r", "echo json_encode(DateTimeZone::listIdentifiers());",
		},
	}).Run()
	if err != nil {
		return timezones, errors.New("ReadPhpTimezonesFailed: " + err.Error())
	}

	err = json.Unmarshal([]byte(timezonesRaw), &timezones)
	if err != nil {
		return timezones, errors.New("ParsePhpTimezonesFailed: " + err.Error())
	}

	return timezones, nil
}

func (repo RuntimeQueryRepo) phpSettingFactory(
	setting string,
) (phpSettingEntity entity.PhpSetting, err error) {
	if setting == "" {
		return phpSettingEntity, errors.New("InvalidPhpSetting")
	}

	separatorIndex := strings.IndexFunc(setting, unicode.IsSpace)
	if separatorIndex < 0 {
		return phpSettingEntity, errors.New("InvalidPhpSetting")
	}

	settingNameStr := strings.TrimSpace(setting[:separatorIndex])
	settingValueStr := strings.TrimSpace(setting[separatorIndex:])
	if settingNameStr == "" || settingValueStr == "" {
		return phpSettingEntity, errors.New("InvalidPhpSetting")
	}

	settingName, err := valueObject.NewPhpSettingName(settingNameStr)
	if err != nil {
		return phpSettingEntity, errors.New("InvalidPhpSettingName")
	}

	settingValue, err := valueObject.NewPhpSettingValue(settingValueStr)
	if err != nil {
		return phpSettingEntity, errors.New("InvalidPhpSettingValue")
	}

	settingOptions := []valueObject.PhpSettingOption{}
	settingOptionValues := []string{}

	switch settingValue.ReadType() {
	case valueObject.PhpSettingValueTypeBool:
		settingOptionValues = []string{"On", "Off"}
	case valueObject.PhpSettingValueTypeNumber:
		settingOptionValues = []string{
			"0", "30", "60", "120", "300", "600", "900", "1800", "3600", "7200",
		}
	case valueObject.PhpSettingValueTypeByteSize:
		lastChar := settingValue[len(settingValue)-1]
		switch lastChar {
		case 'K':
			settingOptionValues = []string{"4096K", "8192K", "16384K"}
		case 'M':
			settingOptionValues = []string{
				"16M", "32M", "64M", "128M", "256M", "512M", "1024M", "2048M",
			}
		case 'G':
			settingOptionValues = []string{"1G", "2G", "4G"}
		}
	}

	switch settingName {
	case "error_reporting":
		settingOptionValues = []string{
			"E_ALL",
			"~E_ALL",
			"E_ALL & ~E_DEPRECATED & ~E_STRICT",
			"E_ALL & ~E_DEPRECATED & ~E_STRICT & ~E_NOTICE & ~E_WARNING",
			"E_ERROR|E_CORE_ERROR|E_COMPILE_ERROR",
		}
	case "date.timezone":
		settingOptionValues, err = repo.readPhpTimezones()
		if err != nil {
			slog.Error(
				"ReadPhpTimezonesFailed",
				slog.String("err", err.Error()),
			)
			settingOptionValues = []string{}
		}
	}

	if len(settingOptionValues) > 0 {
		for _, optionValue := range settingOptionValues {
			settingOption, optionErr := valueObject.NewPhpSettingOption(
				optionValue,
			)
			if optionErr != nil {
				slog.Debug(
					"SkippingInvalidPhpSettingOption",
					slog.String("option", optionValue),
				)
				continue
			}

			settingOptions = append(settingOptions, settingOption)
		}
	}

	settingTypeStr := "text"
	if len(settingOptions) > 0 {
		settingTypeStr = "select"
	}
	settingType, _ := valueObject.NewPhpSettingType(settingTypeStr)

	return entity.NewPhpSetting(
		settingName, settingType, settingValue, settingOptions,
	), nil
}

func (repo RuntimeQueryRepo) readPhpDirectiveSettingLines(
	phpVirtualHostConfFilePath tkValueObject.UnixAbsoluteFilePath,
) (settingLines []string, err error) {
	directiveRegex := regexp.MustCompile(
		`(?m)^[ \t]*` + phpDirectiveKeywordPattern +
			`[ \t]+([^ \t]+)[ \t]+(.+)$`,
	)

	directiveFindings, err := repo.fileClerk.FileContentRegexSearch(
		phpVirtualHostConfFilePath, directiveRegex,
	)
	if err != nil {
		return nil, errors.New("ReadPhpSettingsFailed: " + err.Error())
	}

	for _, directiveFinding := range directiveFindings {
		settingName := directiveFinding.Groups[0]
		settingValue := strings.TrimSpace(directiveFinding.Groups[1])
		settingLines = append(settingLines, settingName+" "+settingValue)
	}

	return settingLines, nil
}

func (repo RuntimeQueryRepo) ReadPhpSettings(
	hostname tkValueObject.Fqdn,
) (phpSettingsEntities []entity.PhpSetting, err error) {
	phpVirtualHostConfFilePath, err := repo.ReadPhpVirtualHostConfFilePath(hostname)
	if err != nil {
		return phpSettingsEntities, err
	}

	settingLines, err := repo.readPhpDirectiveSettingLines(
		phpVirtualHostConfFilePath,
	)
	if err != nil {
		return phpSettingsEntities, err
	}

	for _, settingLine := range settingLines {
		phpSettingEntity, err := repo.phpSettingFactory(settingLine)
		if err != nil {
			slog.Debug(
				"SkippingInvalidPhpSettingLine",
				slog.String("settingLine", settingLine),
				slog.String("err", err.Error()),
			)
			continue
		}

		phpSettingsEntities = append(phpSettingsEntities, phpSettingEntity)
	}

	return phpSettingsEntities, nil
}

func (repo RuntimeQueryRepo) normalizePhpModuleName(
	rawModuleName string,
) string {
	normalizedModuleName := strings.ReplaceAll(rawModuleName, "Zend", "")
	normalizedModuleName = strings.ReplaceAll(
		normalizedModuleName, "Loader", "",
	)

	return strings.ToLower(strings.TrimSpace(normalizedModuleName))
}

func (repo RuntimeQueryRepo) readSupportedPhpModuleNames(
	assetFilePath string,
	version valueObject.PhpVersion,
) ([]valueObject.PhpModuleName, error) {
	rawModulesAsset, err := tkInfra.FileDeserializer(assetFilePath)
	if err != nil {
		return nil, errors.New("ReadPhpModulesAssetFailed: " + err.Error())
	}

	modulesByVersion, assertOk := rawModulesAsset["modules"].(map[string]any)
	if !assertOk {
		return nil, errors.New("InvalidPhpModulesAssetVersionsStructure")
	}

	rawModuleNames, exists := modulesByVersion[version.String()]
	if !exists {
		return nil, errors.New("PhpVersionNotFoundInPhpModulesAsset")
	}

	moduleNames, assertOk := rawModuleNames.([]any)
	if !assertOk {
		return nil, errors.New("InvalidPhpModulesAssetModuleNamesStructure")
	}

	supportedModuleNames := []valueObject.PhpModuleName{}
	for _, rawModuleName := range moduleNames {
		moduleName, err := valueObject.NewPhpModuleName(rawModuleName)
		if err != nil {
			return nil, errors.New("InvalidPhpModuleInAsset: " + err.Error())
		}

		supportedModuleNames = append(supportedModuleNames, moduleName)
	}

	return supportedModuleNames, nil
}

func (repo RuntimeQueryRepo) phpModulesFactory(
	rawPhpModuleOutput string,
	supportedModuleNames []valueObject.PhpModuleName,
) []entity.PhpModule {
	activePhpModuleNames := []string{}
	for rawModuleName := range strings.SplitSeq(rawPhpModuleOutput, "\n") {
		if rawModuleName == "" {
			continue
		}
		isModuleSectionHeader := strings.HasPrefix(rawModuleName, "[") &&
			strings.HasSuffix(rawModuleName, "]")
		if isModuleSectionHeader {
			continue
		}

		normalizedModuleName := repo.normalizePhpModuleName(rawModuleName)
		if normalizedModuleName == "" {
			continue
		}

		activePhpModuleNames = append(activePhpModuleNames, normalizedModuleName)
	}

	phpModulesEntities := []entity.PhpModule{}
	for _, moduleName := range supportedModuleNames {
		isModuleInstalled := slices.Contains(
			activePhpModuleNames, moduleName.String(),
		)
		phpModulesEntities = append(
			phpModulesEntities, entity.NewPhpModule(moduleName, isModuleInstalled),
		)
	}

	return phpModulesEntities
}

func (repo RuntimeQueryRepo) ReadPhpModules(
	version valueObject.PhpVersion,
) (phpModulesEntities []entity.PhpModule, err error) {
	supportedModuleNames, err := repo.readSupportedPhpModuleNames(
		infraEnvs.PhpWebServerModulesAssetFilePath, version,
	)
	if err != nil {
		return phpModulesEntities, err
	}

	rawPhpModuleOutput, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "/usr/local/lsws/lsphp" + version.GetWithoutDots() + "/bin/php",
		Args:    []string{"-m"},
	}).Run()
	if err != nil {
		return phpModulesEntities, errors.New(
			"ReadPhpModulesFailed: " + err.Error(),
		)
	}

	return repo.phpModulesFactory(rawPhpModuleOutput, supportedModuleNames), nil
}

func (repo RuntimeQueryRepo) ReadPhpConfigs(
	hostname tkValueObject.Fqdn,
) (phpConfigsEntity entity.PhpConfigs, err error) {
	phpVersionEntity, err := repo.ReadPhpVersion(hostname)
	if err != nil {
		return phpConfigsEntity, err
	}

	phpSettingsEntities, err := repo.ReadPhpSettings(hostname)
	if err != nil {
		return phpConfigsEntity, err
	}

	phpModulesEntities, err := repo.ReadPhpModules(phpVersionEntity.Value)
	if err != nil {
		return phpConfigsEntity, err
	}

	return entity.NewPhpConfigs(
		hostname, phpVersionEntity, phpSettingsEntities, phpModulesEntities,
	), nil
}
