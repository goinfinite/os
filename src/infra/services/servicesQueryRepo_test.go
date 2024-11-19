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
			t.Errorf("Expected no error, got %v", err)
		}

		if len(services.InstallableServices) == 0 {
			t.Error("NoInstallableItemsFound")
		}
	})

	t.Run("ReadOneInstallableItem", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("node")

		readInstallableItemsRequestDto := dto.ReadInstallableServicesItemsRequest{
			ServiceName: &name,
		}

		_, err := servicesQueryRepo.ReadOneInstallableItem(
			readInstallableItemsRequestDto,
		)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
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
			t.Errorf("Expected no error, got %v", err)
		}

		if len(services.InstalledServices) == 0 {
			t.Error("NoInstalledItemsFound")
		}
	})

	t.Run("ReadOneInstalledItem", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("node")

		readInstalledItemsRequestDto := dto.ReadInstalledServicesItemsRequest{
			ServiceName: &name,
		}

		_, err := servicesQueryRepo.ReadOneInstalledItem(
			readInstalledItemsRequestDto,
		)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
