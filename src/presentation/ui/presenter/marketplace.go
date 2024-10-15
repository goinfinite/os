package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type MarketplacePresenter struct {
	marketplaceService *service.MarketplaceService
}

func NewMarketplacePresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *MarketplacePresenter {
	return &MarketplacePresenter{
		marketplaceService: service.NewMarketplaceService(persistentDbSvc),
	}
}

func (presenter *MarketplacePresenter) Handler(c echo.Context) error {
	responseOutput := presenter.marketplaceService.ReadInstalledItems()
	if responseOutput.Status != service.Success {
		return nil
	}

	installedItems, assertOk := responseOutput.Body.([]entity.MarketplaceInstalledItem)
	if !assertOk {
		return nil
	}

	pageContent := page.MarketplaceIndex(installedItems)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
