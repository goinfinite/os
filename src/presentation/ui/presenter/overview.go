package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/useCase"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	transientDbSvc       *internalDbInfra.TransientDatabaseService
	marketplacePresenter *MarketplacePresenter
}

func NewOverviewPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		transientDbSvc:       transientDbSvc,
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

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(presenter.transientDbSvc)
	o11yOverview, err := useCase.ReadO11yOverview(o11yQueryRepo)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.OverviewIndex(vhostsHostnames, marketplaceOverview, o11yOverview)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
