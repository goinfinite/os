package servicesInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
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
