package presenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	presenterDto "github.com/goinfinite/os/src/presentation/ui/presenter/dto"
	presenterValueObject "github.com/goinfinite/os/src/presentation/ui/presenter/valueObject"
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

func (presenter *MarketplacePresenter) readMarketplaceOverviewByType(
	listType presenterValueObject.MarketplaceListType,
) (overview presenterDto.MarketplaceOverview, err error) {
	var assertOk bool

	installedItemsList := []entity.MarketplaceInstalledItem{}
	if listType.String() == "installed" {
		responseOutput := presenter.marketplaceService.ReadInstalledItems()
		if responseOutput.Status != service.Success {
			return overview, errors.New("FailedToReadInstalledItems")
		}

		installedItemsList, assertOk = responseOutput.Body.([]entity.MarketplaceInstalledItem)
		if !assertOk {
			return overview, errors.New("FailedToReadInstalledItems")
		}
	}

	catalogItemsList := []entity.MarketplaceCatalogItem{}
	if listType.String() == "catalog" {
		responseOutput := presenter.marketplaceService.ReadCatalog()
		if responseOutput.Status != service.Success {
			return overview, errors.New("FailedToReadCatalogItems")
		}

		catalogItemsList, assertOk = responseOutput.Body.([]entity.MarketplaceCatalogItem)
		if !assertOk {
			return overview, errors.New("FailedToReadCatalogItems")
		}
	}

	return presenterDto.NewMarketplaceOverview(
		listType, installedItemsList, catalogItemsList,
	), nil
}

func (presenter *MarketplacePresenter) Handler(c echo.Context) error {
	rawListType := "installed"
	if c.QueryParam("listType") != "" {
		rawListType = c.QueryParam("listType")
	}
	listType, err := presenterValueObject.NewMarketplaceListType(rawListType)
	if err != nil {
		slog.Error(err.Error(), slog.Any("rawListType", rawListType))
		return nil
	}

	marketplaceOverview, err := presenter.readMarketplaceOverviewByType(listType)
	if err != nil {
		return nil
	}

	pageContent := page.MarketplaceIndex(marketplaceOverview)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
