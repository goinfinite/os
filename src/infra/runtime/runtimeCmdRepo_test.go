package runtimeInfra

import (
	"os"
	"path/filepath"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestRuntimeCmdRepo(test *testing.T) {
	test.Skip("SkipRuntimeCmdRepoTest")

	// The integration checks need the application environment when enabled.
	testHelpers.LoadEnvVars()
	persistentDbSvc, err := internalDbInfra.NewPersistentDatabaseService()
	if err != nil {
		test.Fatalf("PersistentDatabaseServiceCreationFailed: %v", err)
	}
	runtimeCmdRepo := NewRuntimeCmdRepo(persistentDbSvc)

	primaryVirtualHost, err := vhostInfra.NewVirtualHostHelpers().
		ReadPrimaryVirtualHostHostname()
	if err != nil {
		test.Fatalf("PrimaryVirtualHostNotFound: %v", err)
	}
	phpVersion, err := valueObject.NewPhpVersion("8.1")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}

	test.Run("UpdatePhpVersion", func(subtest *testing.T) {
		err := runtimeCmdRepo.UpdatePhpVersion(primaryVirtualHost, phpVersion)
		if err != nil {
			subtest.Fatalf("UpdatePhpVersionShouldSucceed: %v", err)
		}
	})

	test.Run("UpdatePhpSettings", func(subtest *testing.T) {
		phpSettingName, err := valueObject.NewPhpSettingName("display_errors")
		if err != nil {
			subtest.Fatalf("PhpSettingNameCreationFailed: %v", err)
		}
		phpSettingType, err := valueObject.NewPhpSettingType("select")
		if err != nil {
			subtest.Fatalf("PhpSettingTypeCreationFailed: %v", err)
		}
		phpSettingValue, err := valueObject.NewPhpSettingValue("Off")
		if err != nil {
			subtest.Fatalf("PhpSettingValueCreationFailed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpSettings(
			primaryVirtualHost,
			[]entity.PhpSetting{
				entity.NewPhpSetting(
					phpSettingName, phpSettingType, phpSettingValue, nil,
				),
			},
		)
		if err != nil {
			subtest.Fatalf("UpdatePhpSettingsShouldSucceed: %v", err)
		}
	})

	test.Run("UpdatePhpModules", func(subtest *testing.T) {
		phpModuleName, err := valueObject.NewPhpModuleName("ioncube")
		if err != nil {
			subtest.Fatalf("PhpModuleNameCreationFailed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			primaryVirtualHost,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, true)},
		)
		if err != nil {
			subtest.Fatalf("UpdatePhpModulesEnableShouldSucceed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpModules(
			primaryVirtualHost,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, false)},
		)
		if err != nil {
			subtest.Fatalf("UpdatePhpModulesDisableShouldSucceed: %v", err)
		}
	})

	test.Run("UpdatePhpVirtualHostHostname", func(subtest *testing.T) {
		newHostname, err := tkValueObject.NewFqdn(
			primaryVirtualHost.String() + ".renamed",
		)
		if err != nil {
			subtest.Fatalf("RenamedVirtualHostCreationFailed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpVirtualHostHostname(
			primaryVirtualHost, newHostname, []tkValueObject.Fqdn{},
		)
		if err != nil {
			subtest.Fatalf("UpdatePhpVirtualHostHostnameShouldSucceed: %v", err)
		}

		err = runtimeCmdRepo.UpdatePhpVirtualHostHostname(
			newHostname, primaryVirtualHost, []tkValueObject.Fqdn{},
		)
		if err != nil {
			subtest.Fatalf(
				"UpdatePhpVirtualHostHostnameReverseShouldSucceed: %v", err,
			)
		}
	})

	test.Run("UpdatePhpVirtualHostHostnameNoOp", func(subtest *testing.T) {
		err := runtimeCmdRepo.UpdatePhpVirtualHostHostname(
			primaryVirtualHost, primaryVirtualHost, []tkValueObject.Fqdn{},
		)
		if err != nil {
			subtest.Fatalf(
				"UpdatePhpVirtualHostHostnameNoOpShouldReturnNil: %v", err,
			)
		}
	})
}

func TestPhpExtensionPackageName(test *testing.T) {
	phpVersion, err := valueObject.NewPhpVersion("8.1")
	if err != nil {
		test.Fatalf("PhpVersionCreationFailed: %v", err)
	}
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)
	testCases := []struct {
		moduleName          string
		expectedPackageName string
	}{
		{moduleName: "curl", expectedPackageName: "lsphp81-curl"},
		{moduleName: "mysqli", expectedPackageName: "lsphp81-mysql"},
		{moduleName: "pdo_mysql", expectedPackageName: "lsphp81-mysql"},
		{moduleName: "pdo_sqlite", expectedPackageName: "lsphp81-sqlite3"},
		{moduleName: "sqlite3", expectedPackageName: "lsphp81-sqlite3"},
	}

	for _, testCase := range testCases {
		test.Run(testCase.moduleName, func(subtest *testing.T) {
			actualPackageName := runtimeCmdRepo.phpExtensionPackageNameResolver(
				phpVersion, testCase.moduleName,
			)
			if actualPackageName != testCase.expectedPackageName {
				subtest.Errorf(
					"PhpExtensionPackageNameMismatch: expected %q, got %q",
					testCase.expectedPackageName,
					actualPackageName,
				)
			}
		})
	}
}

func TestIsPhpModuleIniFilePresent(test *testing.T) {
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)
	tempDirectory := test.TempDir()

	regularFilePath := filepath.Join(tempDirectory, "module.ini")
	err := os.WriteFile(regularFilePath, []byte{}, 0644)
	if err != nil {
		test.Fatalf("CreateRegularIniFileFailed: %v", err)
	}

	// Module configuration can be linked to a shared file.
	symbolicLinkPath := filepath.Join(tempDirectory, "module-link.ini")
	err = os.Symlink(regularFilePath, symbolicLinkPath)
	if err != nil {
		test.Fatalf("CreateIniFileSymlinkFailed: %v", err)
	}

	directoryPath := filepath.Join(tempDirectory, "directory.ini")
	err = os.Mkdir(directoryPath, 0755)
	if err != nil {
		test.Fatalf("CreateIniFileDirectoryFailed: %v", err)
	}

	testCases := []struct {
		testName         string
		filePath         string
		expectedPresence bool
	}{
		{
			testName:         "RegularFile",
			filePath:         regularFilePath,
			expectedPresence: true,
		},
		{
			testName:         "SymbolicLink",
			filePath:         symbolicLinkPath,
			expectedPresence: true,
		},
		{
			testName:         "MissingPath",
			filePath:         filepath.Join(tempDirectory, "missing.ini"),
			expectedPresence: false,
		},
		{
			testName:         "Directory",
			filePath:         directoryPath,
			expectedPresence: false,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			actualPresence, err := runtimeCmdRepo.isPhpModuleIniFilePresent(
				testCase.filePath,
			)
			if err != nil {
				subtest.Fatalf("ReadIniFilePresenceFailed: %v", err)
			}
			if actualPresence != testCase.expectedPresence {
				subtest.Errorf(
					"IniFilePresenceMismatch: expected %v, got %v",
					testCase.expectedPresence,
					actualPresence,
				)
			}
		})
	}
}
