package runtimeInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
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

func TestNormalizedPhpModuleName(test *testing.T) {
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
			actualModuleName := runtimeQueryRepo.normalizedPhpModuleName(
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
