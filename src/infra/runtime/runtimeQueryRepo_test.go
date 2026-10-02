package runtimeInfra

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
)

func TestRuntimeQueryRepo(test *testing.T) {
	test.Skip("SkipRuntimeQueryRepoTest")

	// The integration checks need the application environment when enabled.
	testHelpers.LoadEnvVars()

	runtimeQueryRepo := NewRuntimeQueryRepo()

	test.Run("ReadPhpVersionsInstalled", func(subtest *testing.T) {
		phpVersions, err := runtimeQueryRepo.ReadPhpVersionsInstalled()

		if err != nil {
			subtest.Fatalf("ReadPhpVersionsInstalledFailed: %v", err)
		}

		if len(phpVersions) == 0 {
			subtest.Fatal("ReadPhpVersionsInstalledReturnedNoVersions")
		}
	})

	test.Run("ReadPhpConfigs", func(subtest *testing.T) {
		primaryVirtualHost, err := vhostInfra.NewVirtualHostHelpers().
			ReadPrimaryVirtualHostHostname()
		if err != nil {
			subtest.Fatalf("PrimaryVirtualHostNotFound: %v", err)
		}

		phpConfigs, err := runtimeQueryRepo.ReadPhpConfigs(primaryVirtualHost)

		if err != nil {
			subtest.Fatalf("ReadPhpConfigsFailed: %v", err)
		}

		if len(phpConfigs.Modules) == 0 {
			subtest.Fatal("ReadPhpConfigsReturnedNoModules")
		}
	})
}

func TestNormalizePhpModuleName(test *testing.T) {
	runtimeQueryRepo := NewRuntimeQueryRepo()
	testCases := []struct {
		testName           string
		rawModuleName      string
		expectedModuleName string
	}{
		{
			testName:           "ZendOpcache",
			rawModuleName:      "Zend OPcache",
			expectedModuleName: "opcache",
		},
		{
			testName:           "IonCubeLoader",
			rawModuleName:      "ionCube Loader",
			expectedModuleName: "ioncube",
		},
		{
			testName:           "TrimmedModuleName",
			rawModuleName:      " pcntl ",
			expectedModuleName: "pcntl",
		},
		{
			testName:           "EmptyAfterPrefixRemoval",
			rawModuleName:      "Zend",
			expectedModuleName: "",
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			actualModuleName := runtimeQueryRepo.normalizePhpModuleName(
				testCase.rawModuleName,
			)
			if actualModuleName != testCase.expectedModuleName {
				subtest.Errorf(
					"NormalizedModuleNameMismatch: expected %q, got %q",
					testCase.expectedModuleName,
					actualModuleName,
				)
			}
		})
	}
}

func TestPhpToolModuleBinaryPath(test *testing.T) {
	runtimeQueryRepo := NewRuntimeQueryRepo()
	phpVersion, err := valueObject.NewPhpVersion("8.3")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	moduleName, err := valueObject.NewPhpModuleName("pear")
	if err != nil {
		test.Fatalf("PhpModuleNameCreationFailed: %v", err)
	}

	expectedPath := "/usr/local/lsws/lsphp83/bin/pear"
	actualPath := runtimeQueryRepo.phpToolModuleBinaryPath(phpVersion, moduleName)
	if actualPath != expectedPath {
		test.Errorf(
			"PhpToolModuleBinaryPathMismatch: expected %q, got %q",
			expectedPath, actualPath,
		)
	}
}

func TestReadSupportedPhpModuleNames(test *testing.T) {
	assetFilePath := filepath.Join(test.TempDir(), "modules.yaml")
	assetContent := []byte(`modules:
  "5.6":
    - curl
    - pear
  "8.5":
    - curl
`)
	err := os.WriteFile(assetFilePath, assetContent, 0644)
	if err != nil {
		test.Fatalf("WritePhpModulesAssetFailed: %v", err)
	}

	phpModuleNameFactory := func(
		rawModuleNames ...string,
	) []valueObject.PhpModuleName {
		moduleNames := []valueObject.PhpModuleName{}
		for _, rawModuleName := range rawModuleNames {
			moduleName, err := valueObject.NewPhpModuleName(rawModuleName)
			if err != nil {
				test.Fatalf(
					"PhpModuleNameCreationFailed: %v, name: %s", err, rawModuleName,
				)
			}

			moduleNames = append(moduleNames, moduleName)
		}

		return moduleNames
	}

	testCases := []struct {
		testName               string
		phpVersionValue        string
		expectedModuleNames    []valueObject.PhpModuleName
		expectAbsentPhpVersion bool
	}{
		{
			testName:            "ReadsModulesOfTheRequestedVersion",
			phpVersionValue:     "5.6",
			expectedModuleNames: phpModuleNameFactory("curl", "pear"),
		},
		{
			testName:            "ReadsOnlyTheVersionAskedFor",
			phpVersionValue:     "8.5",
			expectedModuleNames: phpModuleNameFactory("curl"),
		},
		{
			testName:               "FailsWhenTheVersionIsAbsentFromTheAsset",
			phpVersionValue:        "8.0",
			expectAbsentPhpVersion: true,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			phpVersion, err := valueObject.NewPhpVersion(testCase.phpVersionValue)
			if err != nil {
				subtest.Fatalf("PhpVersionCreationFailed: %v", err)
			}

			runtimeQueryRepo := NewRuntimeQueryRepo()
			actualModuleNames, err := runtimeQueryRepo.readSupportedPhpModuleNames(
				assetFilePath, phpVersion,
			)

			if testCase.expectAbsentPhpVersion {
				if err == nil {
					subtest.Errorf(
						"ExpectedAbsentPhpVersionError, version: %s",
						testCase.phpVersionValue,
					)
				}
				return
			}

			if err != nil {
				subtest.Fatalf("ReadSupportedPhpModuleNamesFailed: %v", err)
			}

			if !slices.Equal(actualModuleNames, testCase.expectedModuleNames) {
				subtest.Errorf(
					"SupportedPhpModuleNamesMismatch: expected %v, got %v, version: %s",
					testCase.expectedModuleNames, actualModuleNames,
					testCase.phpVersionValue,
				)
			}
		})
	}
}

func TestPhpModulesFactory(test *testing.T) {
	phpModuleNameFactory := func(
		rawModuleNames ...string,
	) []valueObject.PhpModuleName {
		moduleNames := []valueObject.PhpModuleName{}
		for _, rawModuleName := range rawModuleNames {
			moduleName, err := valueObject.NewPhpModuleName(rawModuleName)
			if err != nil {
				test.Fatalf(
					"PhpModuleNameCreationFailed: %v, name: %s", err, rawModuleName,
				)
			}

			moduleNames = append(moduleNames, moduleName)
		}

		return moduleNames
	}

	phpModuleFactory := func(moduleSpecs ...string) []entity.PhpModule {
		phpModules := []entity.PhpModule{}
		for _, moduleSpec := range moduleSpecs {
			phpModule, err := entity.NewPhpModuleFromString(moduleSpec)
			if err != nil {
				test.Fatalf("PhpModuleCreationFailed: %v, spec: %s", err, moduleSpec)
			}

			phpModules = append(phpModules, phpModule)
		}

		return phpModules
	}

	testCases := []struct {
		testName             string
		rawPhpModuleOutput   string
		supportedModuleNames []valueObject.PhpModuleName
		toolModuleStatuses   map[string]bool
		expectedModules      []entity.PhpModule
	}{
		{
			testName:             "MarksInstalledSupportedModulesAsActive",
			rawPhpModuleOutput:   "[PHP Modules]\ncurl\nZend OPcache\npcntl\n",
			supportedModuleNames: phpModuleNameFactory("curl", "opcache", "ssh2"),
			expectedModules: phpModuleFactory(
				"curl:true", "opcache:true", "ssh2:false",
			),
		},
		{
			testName:             "SkipsSectionHeaders",
			rawPhpModuleOutput:   "[PHP Modules]\ncurl\n[Zend Modules]\nZend OPcache\n",
			supportedModuleNames: phpModuleNameFactory("curl", "opcache"),
			expectedModules:      phpModuleFactory("curl:true", "opcache:true"),
		},
		{
			testName:             "DropsInstalledModulesAbsentFromTheAsset",
			rawPhpModuleOutput:   "[PHP Modules]\ncurl\npcntl\n",
			supportedModuleNames: phpModuleNameFactory("curl"),
			expectedModules:      phpModuleFactory("curl:true"),
		},
		{
			testName:             "MarksToolModuleAsActiveWhenItsBinaryIsPresent",
			rawPhpModuleOutput:   "[PHP Modules]\ncurl\n",
			supportedModuleNames: phpModuleNameFactory("curl", "pear"),
			toolModuleStatuses:   map[string]bool{"pear": true},
			expectedModules:      phpModuleFactory("curl:true", "pear:true"),
		},
		{
			testName:             "MarksToolModuleAsInactiveWhenItsBinaryIsAbsent",
			rawPhpModuleOutput:   "[PHP Modules]\ncurl\n",
			supportedModuleNames: phpModuleNameFactory("curl", "pear"),
			toolModuleStatuses:   map[string]bool{"pear": false},
			expectedModules:      phpModuleFactory("curl:true", "pear:false"),
		},
		{
			testName:             "IgnoresExtensionListingForToolModules",
			rawPhpModuleOutput:   "[PHP Modules]\ncurl\npear\n",
			supportedModuleNames: phpModuleNameFactory("curl", "pear"),
			toolModuleStatuses:   map[string]bool{"pear": false},
			expectedModules:      phpModuleFactory("curl:true", "pear:false"),
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			runtimeQueryRepo := NewRuntimeQueryRepo()
			actualModules := runtimeQueryRepo.phpModulesFactory(
				testCase.rawPhpModuleOutput, testCase.supportedModuleNames,
				testCase.toolModuleStatuses,
			)

			if !slices.Equal(actualModules, testCase.expectedModules) {
				subtest.Errorf(
					"PhpModulesMismatch: expected %v, got %v, output: %s",
					testCase.expectedModules, actualModules,
					testCase.rawPhpModuleOutput,
				)
			}
		})
	}
}

func TestPhpSettingFactory(test *testing.T) {
	runtimeQueryRepo := NewRuntimeQueryRepo()

	testCases := []struct {
		testName      string
		confLine      string
		expectedName  string
		expectedType  string
		expectedValue string
		expectFailure bool
	}{
		{
			testName:      "ReadsByteSizeValue",
			confLine:      "memory_limit 128M",
			expectedName:  "memory_limit",
			expectedType:  "select",
			expectedValue: "128M",
		},
		{
			testName:      "ReadsBoolValue",
			confLine:      "display_errors On",
			expectedName:  "display_errors",
			expectedType:  "select",
			expectedValue: "On",
		},
		{
			testName:      "ReadsMultiWordQuotedValue",
			confLine:      "error_reporting \"E_ALL & ~E_DEPRECATED & ~E_STRICT\"",
			expectedName:  "error_reporting",
			expectedType:  "select",
			expectedValue: "E_ALL & ~E_DEPRECATED & ~E_STRICT",
		},
		{
			testName:      "ReadsPathWithSpaces",
			confLine:      "sendmail_path \"/usr/sbin/sendmail -t -i\"",
			expectedName:  "sendmail_path",
			expectedType:  "text",
			expectedValue: "/usr/sbin/sendmail -t -i",
		},
		{
			testName:      "RejectsLineWithoutValue",
			confLine:      "memory_limit",
			expectFailure: true,
		},
		{
			testName:      "RejectsEmptyLine",
			confLine:      "",
			expectFailure: true,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			phpSetting, err := runtimeQueryRepo.phpSettingFactory(testCase.confLine)
			if testCase.expectFailure {
				if err == nil {
					subtest.Errorf(
						"ExpectedInvalidPhpSettingError, confLine: %q",
						testCase.confLine,
					)
				}
				return
			}

			if err != nil {
				subtest.Fatalf("UnexpectedError: %v, confLine: %q", err, testCase.confLine)
			}
			if phpSetting.Name.String() != testCase.expectedName {
				subtest.Errorf(
					"PhpSettingNameMismatch: expected %q, got %q",
					testCase.expectedName, phpSetting.Name.String(),
				)
			}
			if phpSetting.Type.String() != testCase.expectedType {
				subtest.Errorf(
					"PhpSettingTypeMismatch: expected %q, got %q",
					testCase.expectedType, phpSetting.Type.String(),
				)
			}
			if phpSetting.Value.String() != testCase.expectedValue {
				subtest.Errorf(
					"PhpSettingValueMismatch: expected %q, got %q",
					testCase.expectedValue, phpSetting.Value.String(),
				)
			}
		})
	}
}
