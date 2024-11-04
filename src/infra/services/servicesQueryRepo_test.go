package servicesInfra

import (
	"testing"

	testHelpers "github.com/goinfinite/os/src/devUtils"
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
)

func TestServicesQueryRepo(t *testing.T) {
	testHelpers.LoadEnvVars()
	persistentDbSvc, _ := internalDbInfra.NewPersistentDatabaseService()
	servicesQueryRepo := NewServicesQueryRepo(persistentDbSvc)

	t.Run("ReturnServicesList", func(t *testing.T) {
		name, _ := valueObject.NewServiceName("node")

		paginationDto := useCase.ServicesDefaultPagination
		sortBy, _ := valueObject.NewPaginationSortBy("id")
		sortDirection, _ := valueObject.NewPaginationSortDirection("desc")
		paginationDto.SortBy = &sortBy
		paginationDto.SortDirection = &sortDirection

		readDto := dto.ReadInstalledServicesItemsRequest{
			Pagination: paginationDto,
			Name:       &name,
		}

		services, err := servicesQueryRepo.Read(readDto)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(services.Items) == 0 {
			t.Errorf("Expected a list of services, got %v", services)
		}
	})
}
