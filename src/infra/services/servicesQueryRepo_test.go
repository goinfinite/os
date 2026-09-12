package servicesInfra

import (
	"os"
	"path/filepath"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestServicesQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	servicesQueryRepo := NewServicesQueryRepo(persistentDbSvc)

	t.Run("ReadInstalledItems", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("nginx")

		readInstalledItemsRequestDto := dto.ReadInstalledServicesItemsRequest{
			Pagination:  useCase.ServicesDefaultPagination,
			ServiceName: &name,
		}

		services, err := servicesQueryRepo.ReadInstalledItems(
			readInstalledItemsRequestDto,
		)
		if err != nil {
			t.Errorf("ReadInstalledItemsShouldSucceed: %v", err)
		}

		if len(services.InstalledServices) == 0 {
			t.Error("NoInstalledItemsFound")
		}
	})

	t.Run("ReadFirstInstalledItem", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("nginx")

		readFirstInstalledRequestDto := dto.ReadFirstInstalledServiceItemsRequest{
			ServiceName: &name,
		}

		_, err := servicesQueryRepo.ReadFirstInstalledItem(
			readFirstInstalledRequestDto,
		)
		if err != nil {
			t.Errorf("ReadFirstInstalledItemShouldSucceed: %v", err)
		}
	})

	t.Run("IsInstalled", func(t *testing.T) {
		installedName, _ := valueObject.NewServiceName("nginx")
		isInstalled := servicesQueryRepo.IsInstalled(installedName)
		if !isInstalled {
			t.Error("InstalledServiceShouldReturnTrue")
		}

		missingName, _ := valueObject.NewServiceName("nonexistent-svc-xyz")
		isMissingInstalled := servicesQueryRepo.IsInstalled(missingName)
		if isMissingInstalled {
			t.Error("MissingServiceShouldReturnFalse")
		}
	})
}

func TestServicesQueryRepoInstallableServiceFactory(t *testing.T) {
	servicesQueryRepo := &ServicesQueryRepo{}

	t.Run("ParsesManifestFields", func(t *testing.T) {
		manifestFilePathStr := filepath.Join(t.TempDir(), "manifest.json")
		manifestContent := `{
	"manifestVersion": "v1",
	"name": "node",
	"nature": "solo",
	"type": "runtime",
	"startCmd": "/usr/bin/node",
	"description": "Node.js runtime",
	"versions": ["22.13.1", "20.18.1"],
	"portBindings": ["3000/tcp"],
	"installCmdSteps": ["echo installNode"],
	"installTimeoutSecs": 900
}`
		manifestWriteErr := os.WriteFile(
			manifestFilePathStr, []byte(manifestContent), 0644,
		)
		if manifestWriteErr != nil {
			t.Fatalf("ManifestWriteFailed: %v", manifestWriteErr)
		}

		manifestFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
			manifestFilePathStr, false,
		)
		installableService, err := servicesQueryRepo.installableServiceFactory(
			manifestFilePath,
		)
		if err != nil {
			t.Fatalf("ManifestParseFailed: %v", err)
		}

		if installableService.Name.String() != "node" {
			t.Errorf("ExpectedNodeNameButGot: %s", installableService.Name.String())
		}
		if installableService.Nature.String() != "solo" {
			t.Errorf("ExpectedSoloNatureButGot: %s", installableService.Nature.String())
		}
		if installableService.Type.String() != "runtime" {
			t.Errorf("ExpectedRuntimeTypeButGot: %s", installableService.Type.String())
		}
		if installableService.ManifestVersion.String() != "v1" {
			t.Errorf(
				"ExpectedV1ManifestVersionButGot: %s",
				installableService.ManifestVersion.String(),
			)
		}
		if installableService.StartCmd.String() != "/usr/bin/node" {
			t.Errorf(
				"ExpectedNodeStartCmdButGot: %s",
				installableService.StartCmd.String(),
			)
		}
		if len(installableService.Versions) != 2 {
			t.Errorf("ExpectedTwoVersionsButGot: %d", len(installableService.Versions))
		}
		if len(installableService.PortBindings) != 1 {
			t.Errorf(
				"ExpectedOnePortBindingButGot: %d",
				len(installableService.PortBindings),
			)
		}
		if len(installableService.InstallCmdSteps) != 1 {
			t.Errorf(
				"ExpectedOneInstallCmdStepButGot: %d",
				len(installableService.InstallCmdSteps),
			)
		}
		if installableService.InstallTimeoutSecs != tkValueObject.UnixTime(900) {
			t.Errorf(
				"ExpectedInstallTimeout900ButGot: %d",
				installableService.InstallTimeoutSecs.Int64(),
			)
		}
		if installableService.StopTimeoutSecs != tkValueObject.UnixTime(600) {
			t.Errorf(
				"ExpectedDefaultStopTimeout600ButGot: %d",
				installableService.StopTimeoutSecs.Int64(),
			)
		}
	})

	t.Run("RejectsInvalidManifests", func(t *testing.T) {
		testCases := []struct {
			testName        string
			manifestContent string
			expectedError   string
		}{
			{
				testName:        "MissingRequiredParam",
				manifestContent: `{"name": "node"}`,
				expectedError:   "MissingParam: nature",
			},
			{
				testName: "InvalidVersionsStructure",
				manifestContent: `{"manifestVersion":"v1","name":"node",` +
					`"nature":"solo","type":"runtime","startCmd":"/usr/bin/node",` +
					`"description":"Node.js runtime",` +
					`"installCmdSteps":["echo installNode"],` +
					`"versions":"22.13.1"}`,
				expectedError: "InvalidServiceVersionsStructure",
			},
			{
				testName: "InvalidServiceName",
				manifestContent: `{"manifestVersion":"v1","name":"Node JS",` +
					`"nature":"solo","type":"runtime","startCmd":"/usr/bin/node",` +
					`"description":"Node.js runtime",` +
					`"installCmdSteps":["echo installNode"]}`,
				expectedError: "InvalidServiceName",
			},
		}

		for _, testCase := range testCases {
			t.Run(testCase.testName, func(t *testing.T) {
				manifestFilePathStr := filepath.Join(t.TempDir(), "manifest.json")
				manifestWriteErr := os.WriteFile(
					manifestFilePathStr, []byte(testCase.manifestContent), 0644,
				)
				if manifestWriteErr != nil {
					t.Fatalf("ManifestWriteFailed: %v", manifestWriteErr)
				}

				manifestFilePath, _ := tkValueObject.NewUnixAbsoluteFilePath(
					manifestFilePathStr, false,
				)
				_, err := servicesQueryRepo.installableServiceFactory(
					manifestFilePath,
				)
				if err == nil {
					t.Fatalf("ExpectedErrorButGotNil")
				}
				if err.Error() != testCase.expectedError {
					t.Errorf(
						"ExpectedErrorMismatch: expected %s, got %s",
						testCase.expectedError, err.Error(),
					)
				}
			})
		}
	})
}

func TestServicesQueryRepoParseManifestCmdStepsReplacesInstallPackages(
	t *testing.T,
) {
	servicesQueryRepo := &ServicesQueryRepo{}
	rawCmdSteps := []any{"install_packages -qqy jq"}

	cmdSteps, err := servicesQueryRepo.parseManifestCmdSteps(
		serviceCmdStepTypeInstall, rawCmdSteps,
	)
	if err != nil {
		t.Fatalf("ParseManifestCmdStepsShouldSucceed: %v", err)
	}

	if len(cmdSteps) != 1 {
		t.Fatalf("UnexpectedCmdStepsCount: %d", len(cmdSteps))
	}

	expectedCommand := "DEBIAN_FRONTEND=noninteractive apt-get install -y -qqy jq"
	if cmdSteps[0].String() != expectedCommand {
		t.Errorf(
			"UnexpectedCommand: %q; expected %q",
			cmdSteps[0].String(), expectedCommand,
		)
	}
}

func TestServicesQueryRepoParseManifestCmdStepsPreservesNonInstallCommands(
	t *testing.T,
) {
	servicesQueryRepo := &ServicesQueryRepo{}
	testCases := []struct {
		stepsType       serviceCmdStepType
		rawCommand      string
		expectedCommand string
	}{
		{
			stepsType:       serviceCmdStepTypeUninstall,
			rawCommand:      "install_packages -qqy jq",
			expectedCommand: "install_packages -qqy jq",
		},
		{
			stepsType:       serviceCmdStepTypeStop,
			rawCommand:      "uninstall_packages -qqy jq",
			expectedCommand: "uninstall_packages -qqy jq",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.rawCommand, func(t *testing.T) {
			rawCmdSteps := []any{testCase.rawCommand}

			cmdSteps, err := servicesQueryRepo.parseManifestCmdSteps(
				testCase.stepsType, rawCmdSteps,
			)
			if err != nil {
				t.Fatalf("ParseManifestCmdStepsShouldSucceed: %v", err)
			}

			if len(cmdSteps) != 1 {
				t.Fatalf("UnexpectedCmdStepsCount: %d", len(cmdSteps))
			}

			if cmdSteps[0].String() != testCase.expectedCommand {
				t.Errorf(
					"UnexpectedCommand: %q; expected %q",
					cmdSteps[0].String(), testCase.expectedCommand,
				)
			}
		})
	}
}
