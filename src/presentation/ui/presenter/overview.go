package presenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	transientDbSvc       *internalDbInfra.TransientDatabaseService
	marketplacePresenter *MarketplacePresenter
	servicesService      *service.ServicesService
}

func NewOverviewPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		transientDbSvc:       transientDbSvc,
		marketplacePresenter: NewMarketplacePresenter(persistentDbSvc, trailDbSvc),
		servicesService:      service.NewServicesService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *OverviewPresenter) installableServicesGroupedByTypeFactory(
	installableServicesList []entity.InstallableService,
) page.InstallableServicesGroupedByType {
	installableServicesGroupedByType := page.InstallableServicesGroupedByType{
		Runtime:   []entity.InstallableService{},
		Database:  []entity.InstallableService{},
		Webserver: []entity.InstallableService{},
		Other:     []entity.InstallableService{},
	}

	for _, item := range installableServicesList {
		switch item.Type {
		case valueObject.ServiceTypeRuntime:
			installableServicesGroupedByType.Runtime = append(
				installableServicesGroupedByType.Runtime, item,
			)
		case valueObject.ServiceTypeDatabase:
			installableServicesGroupedByType.Database = append(
				installableServicesGroupedByType.Database, item,
			)
		case valueObject.ServiceTypeWebServer:
			installableServicesGroupedByType.Webserver = append(
				installableServicesGroupedByType.Webserver, item,
			)
		case valueObject.ServiceTypeOther:
			installableServicesGroupedByType.Other = append(
				installableServicesGroupedByType.Other, item,
			)
		}
	}

	return installableServicesGroupedByType
}

func (presenter *OverviewPresenter) readInstalledServices(c echo.Context) (
	responseDto dto.ReadInstalledServicesItemsResponse, err error,
) {
	pageNumber := uint16(0)
	pageNumberQueryParam := c.QueryParam("pageNumber")
	if pageNumberQueryParam != "" {
		pageNumber, _ = voHelper.InterfaceToUint16(pageNumberQueryParam)
	}

	itemsPerPage := uint16(5)
	itemsPerPageQueryParam := c.QueryParam("itemsPerPage")
	if itemsPerPageQueryParam != "" {
		itemsPerPage, _ = voHelper.InterfaceToUint16(itemsPerPageQueryParam)
	}

	readInstalledServicesRequestBody := map[string]interface{}{
		"pageNumber":           pageNumber,
		"itemsPerPage":         itemsPerPage,
		"shouldIncludeMetrics": true,
	}

	nameQueryParam := c.QueryParam("name")
	if nameQueryParam != "" {
		readInstalledServicesRequestBody["name"] = nameQueryParam
	}

	natureQueryParam := c.QueryParam("nature")
	if natureQueryParam != "" {
		readInstalledServicesRequestBody["nature"] = natureQueryParam
	}

	typeQueryParam := c.QueryParam("type")
	if typeQueryParam != "" {
		readInstalledServicesRequestBody["type"] = typeQueryParam
	}

	statusQueryParam := c.QueryParam("status")
	if statusQueryParam != "" {
		readInstalledServicesRequestBody["status"] = statusQueryParam
	}

	installedItemsResponseOutput := presenter.servicesService.ReadInstalledItems(
		readInstalledServicesRequestBody,
	)
	if installedItemsResponseOutput.Status != service.Success {
		return responseDto, errors.New("FailedToReadInstalledServices")
	}

	installedItemsTypedOutputBody, assertOk := installedItemsResponseOutput.Body.(dto.ReadInstalledServicesItemsResponse)
	if !assertOk {
		return responseDto, errors.New("FailedToReadInstalledServices")
	}

	return installedItemsTypedOutputBody, nil
}

func (presenter *OverviewPresenter) servicesOverviewFactory(c echo.Context) (
	overview page.ServicesOverview, err error,
) {
	installedItemsResponseDto, err := presenter.readInstalledServices(c)
	if err != nil {
		return overview, err
	}

	installableItemsResponseOutput := presenter.servicesService.ReadInstallableItems(
		map[string]interface{}{},
	)
	if installableItemsResponseOutput.Status != service.Success {
		return overview, errors.New("FailedToReadInstallableServices")
	}

	installableItemsTypedOutputBody, assertOk := installableItemsResponseOutput.Body.(dto.ReadInstallableServicesItemsResponse)
	if !assertOk {
		return overview, errors.New("FailedToReadInstallableServices")
	}
	installableServicesGroupedByType := presenter.installableServicesGroupedByTypeFactory(
		installableItemsTypedOutputBody.InstallableServices,
	)

	return page.ServicesOverview{
		InstalledServicesResponseDto: installedItemsResponseDto,
		InstallableServices:          installableServicesGroupedByType,
	}, nil
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
	o11yOverview, err := useCase.ReadO11yOverview(o11yQueryRepo, false)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	servicesOverview, err := presenter.servicesOverviewFactory(c)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.OverviewIndex(
		vhostsHostnames, marketplaceOverview, o11yOverview, servicesOverview,
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
