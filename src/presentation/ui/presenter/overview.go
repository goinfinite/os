package presenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	marketplaceService *service.MarketplaceService
}

func NewOverviewPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		marketplaceService: service.NewMarketplaceService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *OverviewPresenter) readMarketplaceInstalledItems() (
	responseDto dto.ReadMarketplaceInstalledItemsResponse, err error,
) {
	responseOutput := presenter.marketplaceService.ReadInstalledItems(
		map[string]interface{}{},
	)
	if responseOutput.Status != service.Success {
		return responseDto, errors.New("FailedToReadInstalledItems")
	}

	typedOutputBody, assertOk := responseOutput.Body.(dto.ReadMarketplaceInstalledItemsResponse)
	if !assertOk {
		return responseDto, errors.New("FailedToReadInstalledItems")
	}
	return typedOutputBody, nil
}

func (presenter *OverviewPresenter) Handler(c echo.Context) error {
	marketplaceInstalledItemsResponseDto, err := presenter.readMarketplaceInstalledItems()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.OverviewIndex(marketplaceInstalledItemsResponseDto)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
