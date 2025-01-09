package presenter

import (
	"log/slog"
	"net/http"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	marketplacePresenter *MarketplacePresenter
}

func NewOverviewPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		marketplacePresenter: NewMarketplacePresenter(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *OverviewPresenter) Handler(c echo.Context) error {
	vhostsHostnames, err := presenter.marketplacePresenter.ReadVhostsHostnames()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	marketplaceOverview, err := presenter.marketplacePresenter.MarketplaceOverviewFactory("all")
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.OverviewIndex(vhostsHostnames, marketplaceOverview)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
