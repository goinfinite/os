package presenter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation/service"
	uiHelper "github.com/speedianet/os/src/presentation/ui/helper"
	"github.com/speedianet/os/src/presentation/ui/page"
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

func (presenter *MappingsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.virtualHostService.ReadWithMappings()
	if responseOutput.Status != service.Success {
		return nil
	}

	vhostsWithMappings, assertOk := responseOutput.Body.([]dto.VirtualHostWithMappings)
	if !assertOk {
		return nil
	}

	pageContent := page.MappingsIndex(vhostsWithMappings)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
