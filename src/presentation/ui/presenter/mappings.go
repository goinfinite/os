package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type MappingsPresenter struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewMappingsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MappingsPresenter {
	return &MappingsPresenter{
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *MappingsPresenter) extractVirtualHostHostnames(
	vhostsWithMappings []dto.VirtualHostWithMappings,
) []string {
	vhostsHostnames := []string{}
	for _, vhostWithMappings := range vhostsWithMappings {
		vhostsHostnames = append(vhostsHostnames, vhostWithMappings.Hostname.String())
	}

	return vhostsHostnames
}

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	virtualHostService := service.NewVirtualHostService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	readVirtualHostsServiceOutput := virtualHostService.ReadWithMappings(map[string]interface{}{
		"itemsPerPage": 1000,
	})
	if readVirtualHostsServiceOutput.Status != service.Success {
		slog.Debug("ReadMappingsServiceOutputBadStatus")
		return nil
	}

	readVirtualHostsResponse, assertOk := readVirtualHostsServiceOutput.Body.(dto.ReadVirtualHostsResponse)
	if !assertOk {
		slog.Debug("ReadMappingsServiceOutputBodyAssertionFailed")
		return nil
	}

	servicesService := service.NewServicesService(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	installedServicesResponseOutput := servicesService.ReadInstalledItems(
		map[string]interface{}{"itemsPerPage": 1000},
	)
	if installedServicesResponseOutput.Status != service.Success {
		slog.Debug("ReadInstalledItemsFailed", slog.Any("output", installedServicesResponseOutput))
		return nil
	}

	installedServicesResponseDto, assertOk := installedServicesResponseOutput.Body.(dto.ReadInstalledServicesItemsResponse)
	if !assertOk {
		return nil
	}

	pageContent := page.MappingsIndex(
		readVirtualHostsResponse.VirtualHostWithMappings,
		presenter.extractVirtualHostHostnames(readVirtualHostsResponse.VirtualHostWithMappings),
		installedServicesResponseDto.InstalledServices,
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
