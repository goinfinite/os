package uiPresenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type MarketplacePresenter struct {
	marketplaceLiaison *liaison.MarketplaceLiaison
	persistentDbSvc    *internalDbInfra.PersistentDatabaseService
	trailDbSvc         *internalDbInfra.TrailDatabaseService
}

func NewMarketplacePresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MarketplacePresenter {
	return &MarketplacePresenter{
		marketplaceLiaison: liaison.NewMarketplaceLiaison(persistentDbSvc, trailDbSvc),
		persistentDbSvc:    persistentDbSvc,
		trailDbSvc:         trailDbSvc,
	}
}

func (presenter *MarketplacePresenter) catalogItemsGroupedByTypeFactory(
	catalogItemsList []entity.MarketplaceCatalogItem,
) CatalogItemsGroupedByType {
	appCatalogItems := []entity.MarketplaceCatalogItem{}
	frameworkCatalogItems := []entity.MarketplaceCatalogItem{}
	stackCatalogItems := []entity.MarketplaceCatalogItem{}
	for _, item := range catalogItemsList {
		switch item.Type.String() {
		case "app":
			appCatalogItems = append(appCatalogItems, item)
		case "framework":
			frameworkCatalogItems = append(frameworkCatalogItems, item)
		case "stack":
			stackCatalogItems = append(stackCatalogItems, item)
		}
	}

	return CatalogItemsGroupedByType{
		Apps:       appCatalogItems,
		Frameworks: frameworkCatalogItems,
		Stacks:     stackCatalogItems,
	}
}

func (presenter *MarketplacePresenter) MarketplaceOverviewFactory(listType string) (
	overview MarketplaceOverview, err error,
) {
	installedItemsList := []entity.MarketplaceInstalledItem{}
	if listType == "installed" || listType == "all" {
		responseOutput := presenter.marketplaceLiaison.ReadInstalledItems(
			map[string]interface{}{},
		)
		if responseOutput.Status != liaison.Success {
			return overview, errors.New("FailedToReadInstalledItems")
		}

		typedOutputBody, assertOk := responseOutput.Body.(dto.ReadMarketplaceInstalledItemsResponse)
		if !assertOk {
			return overview, errors.New("FailedToReadInstalledItems")
		}
		installedItemsList = typedOutputBody.MarketplaceInstalledItems
	}

	catalogItemsList := []entity.MarketplaceCatalogItem{}
	if listType == "catalog" || listType == "all" {
		responseOutput := presenter.marketplaceLiaison.ReadCatalog(
			map[string]interface{}{},
		)
		if responseOutput.Status != liaison.Success {
			return overview, errors.New("FailedToReadCatalogItems")
		}

		typedOutputBody, assertOk := responseOutput.Body.(dto.ReadMarketplaceCatalogItemsResponse)
		if !assertOk {
			return overview, errors.New("FailedToReadCatalogItems")
		}
		catalogItemsList = typedOutputBody.MarketplaceCatalogItems
	}

	return MarketplaceOverview{
		ListType:           listType,
		InstalledItemsList: installedItemsList,
		CatalogItemsList:   presenter.catalogItemsGroupedByTypeFactory(catalogItemsList),
	}, nil
}

func (presenter *MarketplacePresenter) Handler(c echo.Context) error {
	listType := "installed"
	if c.QueryParam("listType") != "" {
		listType = c.QueryParam("listType")
		if listType != "installed" && listType != "catalog" {
			slog.Error("InvalidMarketplaceListType", slog.Any("listType", listType))
			return nil
		}
	}

	vhostsHostnames, err := presenterHelper.ReadVirtualHostHostnames(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	if err != nil {
		slog.Error("ReadVirtualHostsHostnames", slog.String("err", err.Error()))
		return nil
	}

	marketplaceOverview, err := presenter.MarketplaceOverviewFactory(listType)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := MarketplaceIndex(vhostsHostnames, marketplaceOverview)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
