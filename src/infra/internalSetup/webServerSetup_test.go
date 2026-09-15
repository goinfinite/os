package internalSetupInfra

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

const phpWebServerConfFixture string = "httpdworkers 2                \n" +
	"\n" +
	"tuning  {\n" +
	"  maxConnections          2000\n" +
	"}\n" +
	"\n" +
	"extprocessor lsphp85 {\n" +
	"  type                    lsapi\n" +
	"  maxConns                5\n" +
	"  env                     PHP_LSAPI_CHILDREN=5\n" +
	"}\n" +
	"\n" +
	"extprocessor lsphp81 {\n" +
	"  type                    lsapi\n" +
	"  maxConns                5\n" +
	"  env                     PHP_LSAPI_CHILDREN=5\n" +
	"}\n" +
	"\n" +
	"railsDefaults  {\n" +
	"  maxConns                1\n" +
	"  env                     PHP_LSAPI_CHILDREN=7\n" +
	"}\n"

func TestPhpWebServerCapacityCalculator(test *testing.T) {
	webServerSetup := NewWebServerSetup(nil, nil)

	testCases := []struct {
		testName             string
		memoryTotalBytes     uint64
		cpuCores             float64
		expectedCapacity     uint64
		expectedWorkersCount uint64
	}{
		{
			testName:             "SetsFiveProcessesPerGiB",
			memoryTotalBytes:     4 * 1024 * 1024 * 1024,
			cpuCores:             2,
			expectedCapacity:     20,
			expectedWorkersCount: 1,
		},
		{
			testName:             "SetsTwoWorkersOnEightCores",
			memoryTotalBytes:     8 * 1024 * 1024 * 1024,
			cpuCores:             8,
			expectedCapacity:     40,
			expectedWorkersCount: 2,
		},
		{
			testName:             "SetsEightWorkersOnThirtyTwoCores",
			memoryTotalBytes:     32 * 1024 * 1024 * 1024,
			cpuCores:             32,
			expectedCapacity:     160,
			expectedWorkersCount: 8,
		},
		{
			testName:             "CapsCapacityAtThreeHundred",
			memoryTotalBytes:     100 * 1024 * 1024 * 1024,
			cpuCores:             64,
			expectedCapacity:     300,
			expectedWorkersCount: 16,
		},
		{
			testName:             "SetsMinimumWorkersOnFractionalCpu",
			memoryTotalBytes:     1 * 1024 * 1024 * 1024,
			cpuCores:             0.5,
			expectedCapacity:     5,
			expectedWorkersCount: 1,
		},
		{
			testName:             "SetsMinimumCapacityOnSubGiBMemory",
			memoryTotalBytes:     256 * 1024 * 1024,
			cpuCores:             0.5,
			expectedCapacity:     5,
			expectedWorkersCount: 1,
		},
	}

	for _, testCase := range testCases {
		test.Run(testCase.testName, func(subtest *testing.T) {
			memoryTotal, err := tkValueObject.NewByte(testCase.memoryTotalBytes)
			if err != nil {
				subtest.Fatalf("MemoryTotalCreationFailed: %v", err)
			}

			lsapiCapacity, workersCount := webServerSetup.phpWebServerCapacityCalculator(
				memoryTotal, testCase.cpuCores,
			)
			if lsapiCapacity != testCase.expectedCapacity {
				subtest.Errorf(
					"ExpectedLsapiCapacity: got %d, want %d",
					lsapiCapacity, testCase.expectedCapacity,
				)
			}
			if workersCount != testCase.expectedWorkersCount {
				subtest.Errorf(
					"ExpectedWorkersCount: got %d, want %d",
					workersCount, testCase.expectedWorkersCount,
				)
			}
		})
	}
}

func TestPhpWebServerCapacityFileUpdater(test *testing.T) {
	webServerSetup := NewWebServerSetup(nil, nil)

	confFilePath := filepath.Join(test.TempDir(), "httpd_config.conf")
	writeErr := os.WriteFile(
		confFilePath, []byte(phpWebServerConfFixture), 0644,
	)
	if writeErr != nil {
		test.Fatalf("WritePhpWebServerConfFixtureFailed: %v", writeErr)
	}

	lsapiCapacity := uint64(20)
	workersCount := uint64(4)
	err := webServerSetup.phpWebServerCapacityFileUpdater(
		confFilePath, lsapiCapacity, workersCount,
	)
	if err != nil {
		test.Fatalf("PhpWebServerCapacityFileUpdaterFailed: %v", err)
	}

	updatedContentBytes, readErr := os.ReadFile(confFilePath)
	if readErr != nil {
		test.Fatalf("ReadUpdatedConfFailed: %v", readErr)
	}
	updatedContent := string(updatedContentBytes)

	childrenSettingCount := strings.Count(
		updatedContent, "PHP_LSAPI_CHILDREN=20",
	)
	if childrenSettingCount != 2 {
		test.Errorf(
			"ExpectedChildrenSettingInBothLsphpBlocks: got %d, content: %q",
			childrenSettingCount, updatedContent,
		)
	}

	nonLsphpChildrenSettingCount := strings.Count(
		updatedContent, "PHP_LSAPI_CHILDREN=7",
	)
	if nonLsphpChildrenSettingCount != 1 {
		test.Errorf(
			"ExpectedNonLsphpChildrenSettingUnchanged: got %d, content: %q",
			nonLsphpChildrenSettingCount, updatedContent,
		)
	}

	connectionsSettingCount := strings.Count(updatedContent, "maxConns 20")
	if connectionsSettingCount != 2 {
		test.Errorf(
			"ExpectedConnectionsSettingInBothLsphpBlocks: got %d, content: %q",
			connectionsSettingCount, updatedContent,
		)
	}

	workersSettingCount := strings.Count(updatedContent, "httpdworkers 4")
	if workersSettingCount != 1 {
		test.Errorf(
			"ExpectedWorkersSettingOnce: got %d, content: %q",
			workersSettingCount, updatedContent,
		)
	}

	autoCalculatedCommentCount := strings.Count(
		updatedContent, "# AUTO CALCULATED. DO NOT EDIT.",
	)
	if autoCalculatedCommentCount != 5 {
		test.Errorf(
			"ExpectedAutoCalculatedCommentOnEverySetting: got %d, content: %q",
			autoCalculatedCommentCount, updatedContent,
		)
	}

	if !strings.Contains(updatedContent, "maxConnections          2000") {
		test.Errorf(
			"ExpectedTuningMaxConnectionsUnchanged: %q", updatedContent,
		)
	}
	if !strings.Contains(updatedContent, "maxConns                1") {
		test.Errorf(
			"ExpectedDefaultsMaxConnsUnchanged: %q", updatedContent,
		)
	}
}
