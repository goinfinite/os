package presenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type MarketplacePresenter struct {
	marketplaceService *service.MarketplaceService
	virtualHostService *service.VirtualHostService
}

func NewMarketplacePresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *MarketplacePresenter {
	return &MarketplacePresenter{
		marketplaceService: service.NewMarketplaceService(persistentDbSvc, trailDbSvc),
		virtualHostService: service.NewVirtualHostService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *MarketplacePresenter) ReadVhostsHostnames() ([]string, error) {
	vhostHostnames := []string{}

	responseOutput := presenter.virtualHostService.Read()
	if responseOutput.Status != service.Success {
		return vhostHostnames, errors.New("FailedToReadVirtualHosts")
	}

	vhosts, assertOk := responseOutput.Body.([]entity.VirtualHost)
	if !assertOk {
		return vhostHostnames, errors.New("FailedToReadVirtualHosts")
	}

	for _, vhost := range vhosts {
		vhostHostnames = append(vhostHostnames, vhost.Hostname.String())
	}

	return vhostHostnames, nil
}

func (presenter *MarketplacePresenter) CatalogItemsGroupedByTypeFactory(
	catalogItemsList []entity.MarketplaceCatalogItem,
) page.CatalogItemsGroupedByType {
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

	return page.CatalogItemsGroupedByType{
		Apps:       appCatalogItems,
		Frameworks: frameworkCatalogItems,
		Stacks:     stackCatalogItems,
	}
}

func (presenter *MarketplacePresenter) MarketplaceOverviewFactory(listType string) (
	overview page.MarketplaceOverview, err error,
) {
	installedItemsList := []entity.MarketplaceInstalledItem{}
	if listType == "installed" || listType == "all" {
		responseOutput := presenter.marketplaceService.ReadInstalledItems(
			map[string]interface{}{},
		)
		if responseOutput.Status != service.Success {
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
		responseOutput := presenter.marketplaceService.ReadCatalog(
			map[string]interface{}{},
		)
		if responseOutput.Status != service.Success {
			return overview, errors.New("FailedToReadCatalogItems")
		}

		typedOutputBody, assertOk := responseOutput.Body.(dto.ReadMarketplaceCatalogItemsResponse)
		if !assertOk {
			return overview, errors.New("FailedToReadCatalogItems")
		}
		catalogItemsList = typedOutputBody.MarketplaceCatalogItems
	}

	return page.MarketplaceOverview{
		ListType:           listType,
		InstalledItemsList: installedItemsList,
		CatalogItemsList:   presenter.CatalogItemsGroupedByTypeFactory(catalogItemsList),
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

	vhostsHostnames, err := presenter.ReadVhostsHostnames()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	marketplaceOverview, err := presenter.MarketplaceOverviewFactory(listType)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.MarketplaceIndex(vhostsHostnames, marketplaceOverview)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
