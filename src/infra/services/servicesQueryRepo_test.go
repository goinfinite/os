package servicesInfra

import (
	"os"
	"path/filepath"
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TestServicesQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	servicesQueryRepo := NewServicesQueryRepo(persistentDbSvc)

	t.Run("ReadInstallableItems", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("node")

		readInstallableItemsRequestDto := dto.ReadInstallableServicesItemsRequest{
			ServiceName: &name,
		}

		services, err := servicesQueryRepo.ReadInstallableItems(
			readInstallableItemsRequestDto,
		)
		if err != nil {
			t.Errorf("ReadInstallableItemsShouldSucceed: %v", err)
		}

		if len(services.InstallableServices) == 0 {
			t.Error("NoInstallableItemsFound")
		}
	})

	t.Run("ReadFirstInstallableItem", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("node")

		readInstallableItemsRequestDto := dto.ReadInstallableServicesItemsRequest{
			ServiceName: &name,
		}

		_, err := servicesQueryRepo.ReadFirstInstallableItem(
			readInstallableItemsRequestDto,
		)
		if err != nil {
			t.Errorf("ReadFirstInstallableItemShouldSucceed: %v", err)
		}
	})

	t.Run("ReadInstalledItems", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("node")

		readInstalledItemsRequestDto := dto.ReadInstalledServicesItemsRequest{
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
		name, _ := valueObject.NewServiceName("node")

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
		installedName, _ := valueObject.NewServiceName("node")
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

func TestServicesQueryRepoInstallableServiceFactoryExecGroup(t *testing.T) {
	testCases := []struct {
		name        string
		execGroup   string
		expectGroup string
		expectError bool
	}{
		{
			name:        "ValidGroup",
			execGroup:   "WWW-DATA",
			expectGroup: "www-data",
		},
		{
			name: "MissingGroup",
		},
		{
			name:        "InvalidGroup",
			execGroup:   "1invalid",
			expectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			manifestContent := `name: test-service
nature: solo
type: runtime
description: Test service
startCmd: /usr/bin/test-service
installCmdSteps:
  - true
`
			if testCase.execGroup != "" {
				manifestContent += "execGroup: " + testCase.execGroup + "\n"
			}

			manifestFilePath := filepath.Join(t.TempDir(), "manifest.yml")
			err := os.WriteFile(manifestFilePath, []byte(manifestContent), 0600)
			if err != nil {
				t.Fatalf("WriteManifestShouldSucceed: %v", err)
			}

			serviceFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
				manifestFilePath, false,
			)
			if err != nil {
				t.Fatalf("NewManifestPathShouldSucceed: %v", err)
			}

			service, err := (&ServicesQueryRepo{}).installableServiceFactory(
				serviceFilePath,
			)
			if testCase.expectError {
				if err == nil {
					t.Fatal("InvalidExecGroupShouldFail")
				}
				return
			}
			if err != nil {
				t.Fatalf("ReadManifestShouldSucceed: %v", err)
			}

			if testCase.expectGroup == "" {
				if service.ExecGroup != nil {
					t.Fatal("MissingExecGroupShouldRemainUnset")
				}
				return
			}

			if service.ExecGroup == nil {
				t.Fatal("ExecGroupShouldBeParsed")
			}
			if service.ExecGroup.String() != testCase.expectGroup {
				t.Fatalf(
					"UnexpectedExecGroup: %q; expected %q",
					service.ExecGroup.String(), testCase.expectGroup,
				)
			}
		})
	}
}
