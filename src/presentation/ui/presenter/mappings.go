package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type MappingsPresenter struct {
	virtualHostService *service.VirtualHostService
}

func NewMappingsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MappingsPresenter {
	return &MappingsPresenter{
		virtualHostService: service.NewVirtualHostService(persistentDbSvc),
	}
}

func (presenter *MappingsPresenter) getVhostsHostnames(
	vhostsWithMappings []dto.VirtualHostWithMappings,
) []string {
	vhostsHostnames := []string{}
	for _, vhostWithMappings := range vhostsWithMappings {
		vhostsHostnames = append(vhostsHostnames, vhostWithMappings.Hostname.String())
	}

	return vhostsHostnames
}

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.virtualHostService.ReadWithMappings()
	if responseOutput.Status != service.Success {
		return nil
	}

	vhostsWithMappings, assertOk := responseOutput.Body.([]dto.VirtualHostWithMappings)
	if !assertOk {
		return nil
	}

	pageContent := page.MappingsIndex(
		vhostsWithMappings, presenter.getVhostsHostnames(vhostsWithMappings),
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
