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
) *MarketplacePresenter {
	return &MarketplacePresenter{
		marketplaceService: service.NewMarketplaceService(persistentDbSvc),
		virtualHostService: service.NewVirtualHostService(persistentDbSvc),
	}
}

func (presenter *MarketplacePresenter) readVhostsHostnames() ([]string, error) {
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

func (presenter *MarketplacePresenter) catalogItemsGroupedByTypeFactory(
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

func (presenter *MarketplacePresenter) marketplaceOverviewFactory(listType string) (
	overview page.MarketplaceOverview, err error,
) {
	var assertOk bool

	installedItemsList := []entity.MarketplaceInstalledItem{}
	if listType == "installed" {
		responseOutput := presenter.marketplaceService.ReadInstalledItems(
			map[string]interface{}{},
		)
		if responseOutput.Status != service.Success {
			return overview, errors.New("FailedToReadInstalledItems")
		}

		installedItemsList, assertOk = responseOutput.Body.([]entity.MarketplaceInstalledItem)
		if !assertOk {
			return overview, errors.New("FailedToReadInstalledItems")
		}
	}

	catalogItemsList := []entity.MarketplaceCatalogItem{}
	if listType == "catalog" {
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
		catalogItemsList = typedOutputBody.Items
	}

	return page.MarketplaceOverview{
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

	vhostsHostnames, err := presenter.readVhostsHostnames()
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	marketplaceOverview, err := presenter.marketplaceOverviewFactory(listType)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.MarketplaceIndex(vhostsHostnames, marketplaceOverview)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
