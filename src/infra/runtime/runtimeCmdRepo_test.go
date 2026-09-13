package runtimeInfra

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraEnvs "github.com/goinfinite/os/src/infra/envs"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
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

		_, err = runtimeCmdRepo.UpdatePhpModules(
			primaryVirtualHost,
			phpVersion,
			[]entity.PhpModule{entity.NewPhpModule(phpModuleName, true)},
		)
		if err != nil {
			subtest.Fatalf("UpdatePhpModulesEnableShouldSucceed: %v", err)
		}

		_, err = runtimeCmdRepo.UpdatePhpModules(
			primaryVirtualHost,
			phpVersion,
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
		{moduleName: "pdo_dblib", expectedPackageName: "lsphp81-sybase"},
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

func TestPhpModuleUpdateReport(test *testing.T) {
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

	packageInstallFailure := errors.New(
		"InstallPhpModulePackageFailed: exit status 100",
	)

	testCases := []struct {
		testName             string
		requestedModules     []entity.PhpModule
		activeModuleNames    []string
		updateErrors         map[string]error
		expectedUpdatedCount int
		expectedFailures     []string
	}{
		{
			testName:             "ReportsModuleLeftInDesiredStateAsUpdated",
			requestedModules:     phpModuleFactory("curl:true"),
			activeModuleNames:    []string{"curl"},
			expectedUpdatedCount: 1,
		},
		{
			testName:             "ReportsMixedBatchAsOneUpdatedOneFailed",
			requestedModules:     phpModuleFactory("curl:true", "opcache:false"),
			activeModuleNames:    []string{"curl", "opcache"},
			expectedUpdatedCount: 1,
			expectedFailures:     []string{"PhpModuleStatusMismatch"},
		},
		{
			testName:          "ReportsModuleStillActiveAsFailed",
			requestedModules:  phpModuleFactory("opcache:false"),
			activeModuleNames: []string{"opcache"},
			expectedFailures:  []string{"PhpModuleStatusMismatch"},
		},
		{
			testName:          "ReportsWriteErrorWithoutRewrappingIt",
			requestedModules:  phpModuleFactory("curl:true"),
			activeModuleNames: []string{"curl"},
			updateErrors:      map[string]error{"curl": packageInstallFailure},
			expectedFailures:  []string{packageInstallFailure.Error()},
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			activeModuleNames := map[string]any{}
			for _, activeModuleName := range testCase.activeModuleNames {
				activeModuleNames[activeModuleName] = nil
			}

			runtimeCmdRepo := NewRuntimeCmdRepo(nil)
			failedModuleErrors := runtimeCmdRepo.phpModuleStatusDoubleChecker(
				testCase.requestedModules, activeModuleNames, testCase.updateErrors,
			)
			report := runtimeCmdRepo.phpModuleUpdateResponseFactory(
				testCase.requestedModules, failedModuleErrors,
			)

			actualUpdatedCount := len(report.ModulesSuccessfullyUpdated)
			if actualUpdatedCount != testCase.expectedUpdatedCount {
				subtest.Errorf(
					"SuccessfullyUpdatedCountMismatch: expected %d, got %d, modules: %v",
					testCase.expectedUpdatedCount,
					actualUpdatedCount,
					report.ModulesSuccessfullyUpdated,
				)
			}

			actualFailures := report.FailedModulesWithReason
			if len(actualFailures) != len(testCase.expectedFailures) {
				subtest.Fatalf(
					"FailedModulesCountMismatch: expected %d, got %d, failures: %v",
					len(testCase.expectedFailures),
					len(actualFailures),
					actualFailures,
				)
			}

			for index, expectedFailureReason := range testCase.expectedFailures {
				actualFailureReason := actualFailures[index].Reason.String()
				if actualFailureReason != expectedFailureReason {
					subtest.Errorf(
						"PhpModuleFailureReasonMismatch: expected %s, got %s",
						expectedFailureReason, actualFailureReason,
					)
				}
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

func TestBuildPhpSettingLineRegex(test *testing.T) {
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)

	testCases := []struct {
		testName       string
		settingName    string
		confLine       string
		expectsMatch   bool
		expectedPrefix string
		expectedValue  string
	}{
		{
			testName:       "CapturesDirectivePrefixAndValue",
			settingName:    "memory_limit",
			confLine:       "php_value memory_limit 128M",
			expectsMatch:   true,
			expectedPrefix: "php_value memory_limit ",
			expectedValue:  "128M",
		},
		{
			testName:       "CapturesFlagDirective",
			settingName:    "display_errors",
			confLine:       "php_flag display_errors  Off",
			expectsMatch:   true,
			expectedPrefix: "php_flag display_errors  ",
			expectedValue:  "Off",
		},
		{
			testName:    "RejectsLineWithoutDirectivePrefix",
			settingName: "memory_limit",
			confLine:    "memory_limit 128M",
		},
		{
			testName:    "RejectsCommentedDirective",
			settingName: "memory_limit",
			confLine:    "# php_value memory_limit 128M",
		},
		{
			testName:    "TreatsDotInSettingNameAsLiteralChar",
			settingName: "session.name",
			confLine:    "php_value sessionXname \"PHPSESSID\"",
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			phpSettingLineRegex, err := runtimeCmdRepo.phpSettingLineRegexFactory(
				testCase.settingName,
			)
			if err != nil {
				subtest.Fatalf("BuildPhpSettingLineRegexFailed: %v", err)
			}

			submatches := phpSettingLineRegex.FindStringSubmatch(testCase.confLine)
			if !testCase.expectsMatch {
				if submatches != nil {
					subtest.Errorf(
						"ExpectedNoMatch, got: %q, settingName: %s",
						submatches, testCase.settingName,
					)
				}
				return
			}

			if len(submatches) != 3 {
				subtest.Fatalf(
					"ExpectedTwoCapturedGroups, got: %q, settingName: %s",
					submatches, testCase.settingName,
				)
			}
			if submatches[1] != testCase.expectedPrefix {
				subtest.Errorf(
					"DirectivePrefixMismatch: expected %q, got %q",
					testCase.expectedPrefix, submatches[1],
				)
			}
			if submatches[2] != testCase.expectedValue {
				subtest.Errorf(
					"SettingValueMismatch: expected %q, got %q",
					testCase.expectedValue, submatches[2],
				)
			}
		})
	}
}

func TestBuildListenerMapLineRegex(test *testing.T) {
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)

	listenerMapLineFactory := func(domains string) string {
		return "  map                     goinfinite.app " + domains
	}

	testCases := []struct {
		testName     string
		vhostName    string
		confLine     string
		expectsMatch bool
	}{
		{
			testName:     "MatchesCatchAllListenerMap",
			vhostName:    "goinfinite.app",
			confLine:     listenerMapLineFactory("*"),
			expectsMatch: true,
		},
		{
			testName:     "MatchesPlainListenerMap",
			vhostName:    "goinfinite.app",
			confLine:     listenerMapLineFactory("goinfinite.app"),
			expectsMatch: true,
		},
		{
			testName:     "MatchesWildcardListenerMap",
			vhostName:    "goinfinite.app",
			confLine:     listenerMapLineFactory("goinfinite.app, *.goinfinite.app"),
			expectsMatch: true,
		},
		{
			testName:     "MatchesAliasListenerMap",
			vhostName:    "goinfinite.app",
			confLine:     listenerMapLineFactory("www.goinfinite.app, *.goinfinite.app"),
			expectsMatch: true,
		},
		{
			testName:     "RejectsOtherVhostListenerMap",
			vhostName:    "goinfinite.app",
			confLine:     "  map                     blog.goinfinite.app blog.goinfinite.app",
			expectsMatch: false,
		},
		{
			testName:     "RejectsCommentedListenerMap",
			vhostName:    "goinfinite.app",
			confLine:     "# map                     goinfinite.app *",
			expectsMatch: false,
		},
		{
			testName:     "TreatsDotsAsLiteralChars",
			vhostName:    "goinfinite.app",
			confLine:     "  map                     goinfiniteXapp goinfiniteXapp",
			expectsMatch: false,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			mapLineRegex, err := runtimeCmdRepo.listenerMapLineRegexFactory(
				testCase.vhostName,
			)
			if err != nil {
				subtest.Fatalf("BuildListenerMapLineRegexFailed: %v", err)
			}

			actualMatch := mapLineRegex.MatchString(testCase.confLine)
			if actualMatch != testCase.expectsMatch {
				subtest.Errorf(
					"PhpListenerMapLineMatchMismatch: expected %v, got %v, confLine: %q",
					testCase.expectsMatch, actualMatch, testCase.confLine,
				)
			}
		})
	}
}

const testListenerMapLinePrefix = "  map                     "

func writeTempConfFile(
	test *testing.T, fileName, content string,
) tkValueObject.UnixAbsoluteFilePath {
	test.Helper()

	confFilePath := filepath.Join(test.TempDir(), fileName)
	writeErr := os.WriteFile(confFilePath, []byte(content), 0644)
	if writeErr != nil {
		test.Fatalf("WriteConfFileFailed: %v", writeErr)
	}

	filePath, pathErr := tkValueObject.NewUnixAbsoluteFilePath(confFilePath, false)
	if pathErr != nil {
		test.Fatalf("ConfFilePathCreationFailed: %v", pathErr)
	}

	return filePath
}

func resolvePhpWebServerConfOwnerAccount(test *testing.T) *user.User {
	test.Helper()

	confOwnerAccount, lookupErr := user.Lookup(
		infraEnvs.PhpWebServerConfOwnerUsername,
	)
	if lookupErr == nil {
		return confOwnerAccount
	}

	_, createErr := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "useradd",
		Args: []string{
			"--system", "--no-create-home",
			infraEnvs.PhpWebServerConfOwnerUsername,
		},
	}).Run()
	if createErr != nil {
		test.Skipf("PhpWebServerConfOwnerCreationFailed: %v", createErr)
	}

	confOwnerAccount, lookupErr = user.Lookup(
		infraEnvs.PhpWebServerConfOwnerUsername,
	)
	if lookupErr != nil {
		test.Skipf("PhpWebServerConfOwnerLookupFailed: %v", lookupErr)
	}

	return confOwnerAccount
}

func TestCountListenerMapLines(test *testing.T) {
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)

	testCases := []struct {
		testName      string
		confContent   string
		vhostName     string
		expectedCount int
	}{
		{
			testName: "CountsPlainMapLine",
			confContent: testListenerMapLinePrefix +
				"goinfinite.app goinfinite.app\n",
			vhostName:     "goinfinite.app",
			expectedCount: 1,
		},
		{
			testName: "CountsWildcardMapLinesInEveryListener",
			confContent: testListenerMapLinePrefix +
				"goinfinite.app goinfinite.app, *.goinfinite.app\n" +
				testListenerMapLinePrefix +
				"goinfinite.app goinfinite.app, *.goinfinite.app\n",
			vhostName:     "goinfinite.app",
			expectedCount: 2,
		},
		{
			testName: "RejectsOtherVhostMapLine",
			confContent: testListenerMapLinePrefix +
				"blog.goinfinite.app blog.goinfinite.app\n",
			vhostName:     "goinfinite.app",
			expectedCount: 0,
		},
		{
			testName: "RejectsHostnameAsMappedDomain",
			confContent: testListenerMapLinePrefix +
				"other.test goinfinite.app\n",
			vhostName:     "goinfinite.app",
			expectedCount: 0,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			confFilePath := writeTempConfFile(
				subtest, "httpd_config.conf", testCase.confContent,
			)

			count, err := runtimeCmdRepo.countListenerMapLines(
				confFilePath, testCase.vhostName,
			)
			if err != nil {
				subtest.Fatalf("CountListenerMapLinesFailed: %v", err)
			}
			if count != testCase.expectedCount {
				subtest.Errorf(
					"ListenerMapLineCountMismatch: expected %d, got %d",
					testCase.expectedCount, count,
				)
			}
		})
	}
}

func TestIsVirtualHostBlockPresent(test *testing.T) {
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)

	confContent := "virtualhost goinfinite.app {\n" +
		"  vhRoot                  /app/html/\n" +
		"}\n"
	confFilePath := writeTempConfFile(
		test, "httpd_config.conf", confContent,
	)

	isPresent, err := runtimeCmdRepo.isVirtualHostBlockPresent(
		confFilePath, "goinfinite.app",
	)
	if err != nil {
		test.Fatalf("ReadVirtualHostBlockFailed: %v", err)
	}
	if !isPresent {
		test.Errorf("ExpectedPhpVirtualHostBlockFound")
	}

	isMissingPresent, err := runtimeCmdRepo.isVirtualHostBlockPresent(
		confFilePath, "blog.goinfinite.app",
	)
	if err != nil {
		test.Fatalf("ReadVirtualHostBlockFailed: %v", err)
	}
	if isMissingPresent {
		test.Errorf("ExpectedPhpVirtualHostBlockAbsent")
	}
}

func TestMapVirtualHostOnEveryListener(test *testing.T) {
	resolvePhpWebServerConfOwnerAccount(test)
	test.Setenv(infraEnvs.PrimaryVirtualHostEnvKey, "goinfinite.app")

	runtimeCmdRepo := NewRuntimeCmdRepo(nil)
	vhostHostname, err := tkValueObject.NewFqdn("blog.goinfinite.app")
	if err != nil {
		test.Fatalf("VirtualHostHostnameCreationFailed: %v", err)
	}

	primaryListenerConf := "listener HTTP {\n" +
		"  address                 *:8080\n" +
		testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
		"}\n" +
		"listener HTTPS {\n" +
		"  address                 *:8443\n" +
		testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
		"}\n"
	expectedMapLine := testListenerMapLinePrefix +
		"blog.goinfinite.app blog.goinfinite.app, *.blog.goinfinite.app"

	test.Run("InsertsMissingMapForEveryListener", func(subtest *testing.T) {
		confFilePath := writeTempConfFile(
			subtest, "httpd_config.conf", primaryListenerConf,
		)

		err := runtimeCmdRepo.mapVirtualHostOnEveryListener(
			vhostHostname, confFilePath,
		)
		if err != nil {
			subtest.Fatalf("MapVirtualHostOnEveryListenerFailed: %v", err)
		}

		confContent, readErr := os.ReadFile(confFilePath.String())
		if readErr != nil {
			subtest.Fatalf("ReadHttpdConfigFailed: %v", readErr)
		}

		actualCount := strings.Count(string(confContent), expectedMapLine)
		if actualCount != 2 {
			subtest.Errorf(
				"ExpectedInsertedMapLineTwice, got %d, content: %q",
				actualCount, confContent,
			)
		}
	})

	test.Run("RepairsMapMissingFromOneListener", func(subtest *testing.T) {
		confContent := "listener HTTP {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			testListenerMapLinePrefix +
			"blog.goinfinite.app blog.goinfinite.app\n" +
			"}\n" +
			"listener HTTPS {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			"}\n"
		confFilePath := writeTempConfFile(
			subtest, "httpd_config.conf", confContent,
		)

		err := runtimeCmdRepo.mapVirtualHostOnEveryListener(
			vhostHostname, confFilePath,
		)
		if err != nil {
			subtest.Fatalf("MapVirtualHostOnEveryListenerFailed: %v", err)
		}

		readContent, readErr := os.ReadFile(confFilePath.String())
		if readErr != nil {
			subtest.Fatalf("ReadHttpdConfigFailed: %v", readErr)
		}

		actualCount := strings.Count(string(readContent), expectedMapLine)
		if actualCount != 2 {
			subtest.Errorf(
				"ExpectedRepairedMapLineTwice, got %d, content: %q",
				actualCount, readContent,
			)
		}
	})

	test.Run("KeepsExistingMapLine", func(subtest *testing.T) {
		confContent := "listener HTTP {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			testListenerMapLinePrefix +
			"blog.goinfinite.app blog.goinfinite.app\n" +
			"}\n"
		confFilePath := writeTempConfFile(
			subtest, "httpd_config.conf", confContent,
		)

		err := runtimeCmdRepo.mapVirtualHostOnEveryListener(
			vhostHostname, confFilePath,
		)
		if err != nil {
			subtest.Fatalf("MapVirtualHostOnEveryListenerFailed: %v", err)
		}

		readContent, readErr := os.ReadFile(confFilePath.String())
		if readErr != nil {
			subtest.Fatalf("ReadHttpdConfigFailed: %v", readErr)
		}
		if string(readContent) != confContent {
			subtest.Errorf("ExpectedUnchangedConfig, got: %q", readContent)
		}
	})

	test.Run("FailsWhenPrimaryMapLineIsMissing", func(subtest *testing.T) {
		confFilePath := writeTempConfFile(
			subtest, "httpd_config.conf",
			"listener HTTP {\n  address *:8080\n}\n",
		)

		err := runtimeCmdRepo.mapVirtualHostOnEveryListener(
			vhostHostname, confFilePath,
		)
		if err == nil {
			subtest.Fatalf("ExpectedPrimaryListenerMapLineNotFoundError")
		}
		if !strings.Contains(err.Error(), "PrimaryListenerMapLineNotFound") {
			subtest.Errorf("UnexpectedError: %v", err)
		}
	})
}

func TestListenerMapLinesDoubleChecker(test *testing.T) {
	resolvePhpWebServerConfOwnerAccount(test)
	test.Setenv(infraEnvs.PrimaryVirtualHostEnvKey, "goinfinite.app")

	runtimeCmdRepo := NewRuntimeCmdRepo(nil)
	vhostHostname, err := tkValueObject.NewFqdn("blog.goinfinite.app")
	if err != nil {
		test.Fatalf("VirtualHostHostnameCreationFailed: %v", err)
	}

	test.Run("AcceptsMapPresentInEveryListener", func(subtest *testing.T) {
		confContent := "listener HTTP {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			testListenerMapLinePrefix +
			"blog.goinfinite.app blog.goinfinite.app\n" +
			"}\n" +
			"listener HTTPS {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			testListenerMapLinePrefix +
			"blog.goinfinite.app blog.goinfinite.app\n" +
			"}\n"
		confFilePath := writeTempConfFile(
			subtest, "httpd_config.conf", confContent,
		)

		err := runtimeCmdRepo.listenerMapLinesDoubleChecker(
			vhostHostname, confFilePath,
		)
		if err != nil {
			subtest.Fatalf("ListenerMapDoubleCheckFailed: %v", err)
		}
	})

	test.Run("FailsWhenMapMissingFromOneListener", func(subtest *testing.T) {
		confContent := "listener HTTP {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			testListenerMapLinePrefix +
			"blog.goinfinite.app blog.goinfinite.app\n" +
			"}\n" +
			"listener HTTPS {\n" +
			testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
			"}\n"
		confFilePath := writeTempConfFile(
			subtest, "httpd_config.conf", confContent,
		)

		err := runtimeCmdRepo.listenerMapLinesDoubleChecker(
			vhostHostname, confFilePath,
		)
		if err == nil {
			subtest.Fatalf("ExpectedPhpVirtualHostListenerMapNotFoundError")
		}
		if !strings.Contains(err.Error(), "PhpVirtualHostListenerMapNotFound") {
			subtest.Errorf("UnexpectedError: %v", err)
		}
	})
}

func TestRemoveListenerMapLines(test *testing.T) {
	resolvePhpWebServerConfOwnerAccount(test)
	test.Setenv(infraEnvs.PrimaryVirtualHostEnvKey, "goinfinite.app")

	runtimeCmdRepo := NewRuntimeCmdRepo(nil)
	vhostHostname, err := tkValueObject.NewFqdn("blog.goinfinite.app")
	if err != nil {
		test.Fatalf("VirtualHostHostnameCreationFailed: %v", err)
	}

	confContent := "listener HTTP {\n" +
		testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
		testListenerMapLinePrefix +
		"blog.goinfinite.app blog.goinfinite.app\n" +
		"}\n" +
		"\n" +
		"virtualhost blog.goinfinite.app {\n" +
		"  vhRoot                  /app/html/blog.goinfinite.app/\n" +
		"  configFile              /app/conf/php-webserver/blog.conf\n" +
		"}\n"
	mainConfFilePath := writeTempConfFile(
		test, "httpd_config.conf", confContent,
	)

	err = runtimeCmdRepo.removeListenerMapLines(
		vhostHostname, mainConfFilePath,
	)
	if err != nil {
		test.Fatalf("RemoveListenerMapLinesFailed: %v", err)
	}

	readContent, readErr := os.ReadFile(mainConfFilePath.String())
	if readErr != nil {
		test.Fatalf("ReadHttpdConfigFailed: %v", readErr)
	}
	contentStr := string(readContent)
	if strings.Contains(contentStr, testListenerMapLinePrefix+"blog.goinfinite.app") {
		test.Errorf("ExpectedPhpVirtualHostMapLineRemoved, got: %q", contentStr)
	}
	if !strings.Contains(contentStr, testListenerMapLinePrefix+"goinfinite.app") {
		test.Errorf("ExpectedPrimaryMapLineKept, got: %q", contentStr)
	}
	if !strings.Contains(contentStr, "virtualhost blog.goinfinite.app") {
		test.Errorf("ExpectedPhpVirtualHostBlockKept, got: %q", contentStr)
	}
}

func TestRemoveVirtualHostBlock(test *testing.T) {
	resolvePhpWebServerConfOwnerAccount(test)
	test.Setenv(infraEnvs.PrimaryVirtualHostEnvKey, "goinfinite.app")

	runtimeCmdRepo := NewRuntimeCmdRepo(nil)
	vhostHostname, err := tkValueObject.NewFqdn("blog.goinfinite.app")
	if err != nil {
		test.Fatalf("VirtualHostHostnameCreationFailed: %v", err)
	}

	confContent := "listener HTTP {\n" +
		testListenerMapLinePrefix + "goinfinite.app goinfinite.app\n" +
		testListenerMapLinePrefix +
		"blog.goinfinite.app blog.goinfinite.app\n" +
		"}\n" +
		"\n" +
		"virtualhost blog.goinfinite.app {\n" +
		"  vhRoot                  /app/html/blog.goinfinite.app/\n" +
		"  configFile              /app/conf/php-webserver/blog.conf\n" +
		"}\n"
	mainConfFilePath := writeTempConfFile(
		test, "httpd_config.conf", confContent,
	)

	err = runtimeCmdRepo.removeVirtualHostBlock(
		vhostHostname, mainConfFilePath,
	)
	if err != nil {
		test.Fatalf("RemoveVirtualHostBlockFailed: %v", err)
	}

	readContent, readErr := os.ReadFile(mainConfFilePath.String())
	if readErr != nil {
		test.Fatalf("ReadHttpdConfigFailed: %v", readErr)
	}
	contentStr := string(readContent)
	if strings.Contains(contentStr, "virtualhost blog.goinfinite.app") {
		test.Errorf("ExpectedPhpVirtualHostBlockRemoved, got: %q", contentStr)
	}
	if !strings.Contains(contentStr, testListenerMapLinePrefix+"blog.goinfinite.app") {
		test.Errorf("ExpectedPhpVirtualHostMapLineKept, got: %q", contentStr)
	}
	if !strings.Contains(contentStr, testListenerMapLinePrefix+"goinfinite.app") {
		test.Errorf("ExpectedPrimaryMapLineKept, got: %q", contentStr)
	}
}

func TestArePhpSettingLinesInPlace(test *testing.T) {
	runtimeCmdRepo := NewRuntimeCmdRepo(nil)

	const confContent = "php_value max_execution_time 30\n" +
		"php_value memory_limit 128M\n" +
		"php_flag display_errors  Off\n" +
		"php_value session.name \"PHPSESSID\"\n" +
		"php_value sessionXname \"NOPE\"\n" +
		"php_flag soap.wsdl_cache_enabled 1\n" +
		"# php_flag soap.wsdl_cache_enabled 0\n" +
		"# php_value sendmail_path \"/usr/sbin/sendmail -t -i\"\n"

	testCases := []struct {
		testName          string
		settingName       string
		desiredValue      string
		expectedIsInPlace bool
	}{
		{
			testName:          "AlreadyInPlace",
			settingName:       "memory_limit",
			desiredValue:      "128M",
			expectedIsInPlace: true,
		},
		{
			testName:          "DifferentValueNeedsWrite",
			settingName:       "memory_limit",
			desiredValue:      "256M",
			expectedIsInPlace: false,
		},
		{
			testName:          "ExtraWhitespaceIsTolerated",
			settingName:       "display_errors",
			desiredValue:      "Off",
			expectedIsInPlace: true,
		},
		{
			testName:          "QuotedStringNeedsExactMatch",
			settingName:       "session.name",
			desiredValue:      "\"PHPSESSID\"",
			expectedIsInPlace: true,
		},
		{
			testName:          "SimilarNameIsNotTheSetting",
			settingName:       "session.name",
			desiredValue:      "\"NOPE\"",
			expectedIsInPlace: false,
		},
		{
			testName:          "AbsentDirectiveNeedsWrite",
			settingName:       "upload_max_filesize",
			desiredValue:      "64M",
			expectedIsInPlace: false,
		},
		{
			testName:          "CommentIsNotAnActiveDirective",
			settingName:       "sendmail_path",
			desiredValue:      "\"/usr/sbin/sendmail -t -i\"",
			expectedIsInPlace: false,
		},
		{
			testName:          "CommentCannotMaskActiveDirective",
			settingName:       "soap.wsdl_cache_enabled",
			desiredValue:      "1",
			expectedIsInPlace: true,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			phpSettingLineRegex, err := runtimeCmdRepo.phpSettingLineRegexFactory(
				testCase.settingName,
			)
			if err != nil {
				subtest.Fatalf("BuildPhpSettingLineRegexFailed: %v", err)
			}

			actualIsInPlace := runtimeCmdRepo.phpSettingLinesMatchDesiredValue(
				phpSettingLineRegex, testCase.desiredValue, confContent,
			)
			if actualIsInPlace != testCase.expectedIsInPlace {
				subtest.Errorf(
					"PhpSettingInPlaceMismatch: expected %v, got %v, settingName: %s",
					testCase.expectedIsInPlace,
					actualIsInPlace,
					testCase.settingName,
				)
			}
		})
	}
}

func TestApplyPhpSettingsToConfFile(test *testing.T) {
	confOwnerAccount := resolvePhpWebServerConfOwnerAccount(test)

	webServerAccount, lookupErr := user.Lookup(infraEnvs.PhpWebServerUsername)
	if lookupErr != nil {
		test.Skipf("PhpWebServerAccountLookupFailed: %v", lookupErr)
	}

	confOwnerUid, parseErr := strconv.Atoi(confOwnerAccount.Uid)
	if parseErr != nil {
		test.Fatalf("PhpWebServerConfOwnerIdParseFailed: %v", parseErr)
	}
	confOwnerGid, parseErr := strconv.Atoi(confOwnerAccount.Gid)
	if parseErr != nil {
		test.Fatalf("PhpWebServerConfOwnerGroupParseFailed: %v", parseErr)
	}
	webServerUid, parseErr := strconv.Atoi(webServerAccount.Uid)
	if parseErr != nil {
		test.Fatalf("PhpWebServerAccountIdParseFailed: %v", parseErr)
	}
	webServerGid, parseErr := strconv.Atoi(webServerAccount.Gid)
	if parseErr != nil {
		test.Fatalf("PhpWebServerAccountGroupParseFailed: %v", parseErr)
	}

	appDir := filepath.Join(test.TempDir(), "app")
	confDir := filepath.Join(appDir, "conf")
	phpWebServerConfDir := filepath.Join(confDir, "php-webserver")
	mkdirErr := os.MkdirAll(phpWebServerConfDir, 0755)
	if mkdirErr != nil {
		test.Fatalf("PhpWebServerConfDirCreationFailed: %v", mkdirErr)
	}

	chownErr := os.Chown(appDir, webServerUid, webServerGid)
	if chownErr != nil {
		test.Fatalf("AppDirOwnershipChangeFailed: %v", chownErr)
	}
	chownErr = os.Chown(confDir, webServerUid, webServerGid)
	if chownErr != nil {
		test.Fatalf("ConfDirOwnershipChangeFailed: %v", chownErr)
	}
	chownErr = os.Chown(phpWebServerConfDir, confOwnerUid, confOwnerGid)
	if chownErr != nil {
		test.Fatalf("PhpWebServerConfDirOwnershipChangeFailed: %v", chownErr)
	}

	confFileFactory := func(
		subtest *testing.T, fileName string, settingLines []string,
	) tkValueObject.UnixAbsoluteFilePath {
		subtest.Helper()

		confContent := "phpIniOverride  {\n"
		for _, settingLine := range settingLines {
			confContent += settingLine + "\n"
		}
		confContent += "}\n"

		confFilePath := filepath.Join(phpWebServerConfDir, fileName)
		writeErr := os.WriteFile(confFilePath, []byte(confContent), 0644)
		if writeErr != nil {
			subtest.Fatalf("WritePhpConfFileFailed: %v", writeErr)
		}

		filePath, pathErr := tkValueObject.NewUnixAbsoluteFilePath(
			confFilePath, false,
		)
		if pathErr != nil {
			subtest.Fatalf("PhpConfFilePathCreationFailed: %v", pathErr)
		}

		return filePath
	}

	phpSettingFactory := func(
		subtest *testing.T, name, value string,
	) entity.PhpSetting {
		subtest.Helper()

		settingName, nameErr := valueObject.NewPhpSettingName(name)
		if nameErr != nil {
			subtest.Fatalf(
				"PhpSettingNameCreationFailed: %v, name: %s", nameErr, name,
			)
		}
		settingValue, valueErr := valueObject.NewPhpSettingValue(value)
		if valueErr != nil {
			subtest.Fatalf(
				"PhpSettingValueCreationFailed: %v, value: %s", valueErr, value,
			)
		}
		settingType, _ := valueObject.NewPhpSettingType("select")

		return entity.NewPhpSetting(settingName, settingType, settingValue, nil)
	}

	runtimeCmdRepo := NewRuntimeCmdRepo(nil)

	test.Run("WritesSettingsThroughMixedOwnerChain", func(subtest *testing.T) {
		confFilePath := confFileFactory(subtest, "mixedOwner.conf", []string{
			"php_value memory_limit 128M",
			"php_flag display_errors On",
			"php_value session.save_path \"/tmp\"",
		})

		hasUpdated, err := runtimeCmdRepo.phpSettingsConfFileUpdater(
			confFilePath,
			[]entity.PhpSetting{
				phpSettingFactory(subtest, "memory_limit", "256M"),
				phpSettingFactory(subtest, "display_errors", "Off"),
				phpSettingFactory(subtest, "session.save_path", "/tmp/$USER"),
			},
		)
		if err != nil {
			subtest.Fatalf("ApplyPhpSettingsFailed: %v", err)
		}
		if !hasUpdated {
			subtest.Errorf("ExpectedSettingsUpdated")
		}

		confContent, readErr := os.ReadFile(confFilePath.String())
		if readErr != nil {
			subtest.Fatalf("ReadPhpConfFileFailed: %v", readErr)
		}
		confContentStr := string(confContent)

		expectedLines := []string{
			"php_value memory_limit 256M",
			"php_flag display_errors Off",
			"php_value session.save_path \"/tmp/$USER\"",
		}
		for _, expectedLine := range expectedLines {
			if !strings.Contains(confContentStr, expectedLine) {
				subtest.Errorf(
					"ExpectedConfLineMissing: %q, content: %q",
					expectedLine, confContentStr,
				)
			}
		}
	})

	test.Run("SkipsSettingsAlreadyInPlace", func(subtest *testing.T) {
		confFilePath := confFileFactory(subtest, "inPlace.conf", []string{
			"php_value memory_limit 256M",
		})

		hasUpdated, err := runtimeCmdRepo.phpSettingsConfFileUpdater(
			confFilePath,
			[]entity.PhpSetting{phpSettingFactory(subtest, "memory_limit", "256M")},
		)
		if err != nil {
			subtest.Fatalf("ApplyPhpSettingsFailed: %v", err)
		}
		if hasUpdated {
			subtest.Errorf("ExpectedNoSettingsUpdated")
		}
	})

	test.Run("FailsWhenNoDirectiveMatches", func(subtest *testing.T) {
		confFilePath := confFileFactory(subtest, "missingDirective.conf", []string{
			"php_value memory_limit 128M",
		})

		hasUpdated, err := runtimeCmdRepo.phpSettingsConfFileUpdater(
			confFilePath,
			[]entity.PhpSetting{phpSettingFactory(subtest, "upload_max_filesize", "64M")},
		)
		if err == nil {
			subtest.Fatalf("ExpectedNoSettingWrittenError")
		}
		if hasUpdated {
			subtest.Errorf("ExpectedNoSettingsUpdated")
		}
		if !strings.Contains(err.Error(), "NoSettingWritten") {
			subtest.Errorf("ExpectedNoSettingWrittenError, got: %v", err)
		}
	})

	test.Run("WritesRemainingSettingsWhenOneDirectiveIsMissing", func(
		subtest *testing.T,
	) {
		confFilePath := confFileFactory(subtest, "partial.conf", []string{
			"php_value memory_limit 128M",
		})

		hasUpdated, err := runtimeCmdRepo.phpSettingsConfFileUpdater(
			confFilePath,
			[]entity.PhpSetting{
				phpSettingFactory(subtest, "memory_limit", "256M"),
				phpSettingFactory(subtest, "upload_max_filesize", "64M"),
			},
		)
		if err != nil {
			subtest.Fatalf("ApplyPhpSettingsFailed: %v", err)
		}
		if !hasUpdated {
			subtest.Errorf("ExpectedSettingsUpdated")
		}

		confContent, readErr := os.ReadFile(confFilePath.String())
		if readErr != nil {
			subtest.Fatalf("ReadPhpConfFileFailed: %v", readErr)
		}
		if !strings.Contains(string(confContent), "php_value memory_limit 256M") {
			subtest.Errorf(
				"ExpectedUpdatedDirective, content: %q", confContent,
			)
		}
	})
}
