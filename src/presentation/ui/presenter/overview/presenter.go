package uiPresenter

import (
	"errors"
	"log/slog"
	"net/http"

	tkPresentation "github.com/goinfinite/tk/src/presentation"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	o11yInfra "github.com/goinfinite/os/src/infra/o11y"
	"github.com/goinfinite/os/src/presentation/liaison"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	presenterMarketplace "github.com/goinfinite/os/src/presentation/ui/presenter/marketplace"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfraDb "github.com/goinfinite/tk/src/infra/db"
	"github.com/labstack/echo/v4"
)

type OverviewPresenter struct {
	persistentDbSvc        *internalDbInfra.PersistentDatabaseService
	transientDbSvc         *tkInfraDb.TransientDatabaseService
	trailDbSvc             *internalDbInfra.TrailDatabaseService
	marketplacePresenter   *presenterMarketplace.MarketplacePresenter
	servicesLiaison        *liaison.ServicesLiaison
	terminalSessionLiaison *liaison.TerminalSessionLiaison
}

func NewOverviewPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *tkInfraDb.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *OverviewPresenter {
	return &OverviewPresenter{
		persistentDbSvc:        persistentDbSvc,
		transientDbSvc:         transientDbSvc,
		trailDbSvc:             trailDbSvc,
		marketplacePresenter:   presenterMarketplace.NewMarketplacePresenter(persistentDbSvc, trailDbSvc),
		servicesLiaison:        liaison.NewServicesLiaison(persistentDbSvc, trailDbSvc),
		terminalSessionLiaison: liaison.NewTerminalSessionLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *OverviewPresenter) installableServicesGroupedByTypeFactory(
	installableServicesList []entity.InstallableService,
) InstallableServicesGroupedByType {
	installableServicesGroupedByType := InstallableServicesGroupedByType{
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
	readInstalledServicesRequestBody := presenterHelper.ReadPaginationRequestParams(c)
	readInstalledServicesRequestBody["shouldIncludeMetrics"] = true

	natureQueryParam := c.QueryParam("nature")
	if natureQueryParam != "" {
		readInstalledServicesRequestBody["nature"] = natureQueryParam
	}

	nameQueryParam := c.QueryParam("name")
	if nameQueryParam != "" {
		readInstalledServicesRequestBody["name"] = nameQueryParam
	}

	typeQueryParam := c.QueryParam("type")
	if typeQueryParam != "" {
		readInstalledServicesRequestBody["type"] = typeQueryParam
	}

	statusQueryParam := c.QueryParam("status")
	if statusQueryParam != "" {
		readInstalledServicesRequestBody["status"] = statusQueryParam
	}

	installedItemsResponseOutput := presenter.servicesLiaison.ReadInstalledItems(
		readInstalledServicesRequestBody,
	)
	if installedItemsResponseOutput.Status != tkPresentation.LiaisonResponseStatusSuccess {
		return responseDto, errors.New("ReadInstalledServicesLiaisonBadResponse")
	}

	installedItemsTypedOutputBody, assertOk := installedItemsResponseOutput.Body.(dto.ReadInstalledServicesItemsResponse)
	if !assertOk {
		return responseDto, errors.New("AssertReadInstalledServicesResponseFailed")
	}

	return installedItemsTypedOutputBody, nil
}

func (presenter *OverviewPresenter) servicesOverviewFactory(c echo.Context) (
	overview ServicesOverview, err error,
) {
	installedItemsResponseDto, err := presenter.readInstalledServices(c)
	if err != nil {
		return overview, err
	}

	installableItemsResponseOutput := presenter.servicesLiaison.ReadInstallableItems(
		map[string]any{
			"itemsPerPage": 50,
		},
	)
	if installableItemsResponseOutput.Status != tkPresentation.LiaisonResponseStatusSuccess {
		return overview, errors.New("ReadInstallableServicesLiaisonBadResponse")
	}

	installableItemsTypedOutputBody, assertOk := installableItemsResponseOutput.Body.(dto.ReadInstallableServicesItemsResponse)
	if !assertOk {
		return overview, errors.New("AssertReadInstallableServicesResponseFailed")
	}
	installableServicesGroupedByType := presenter.installableServicesGroupedByTypeFactory(
		installableItemsTypedOutputBody.InstallableServices,
	)

	return ServicesOverview{
		InstalledServicesResponseDto: installedItemsResponseDto,
		InstallableServices:          installableServicesGroupedByType,
	}, nil
}

func (presenter *OverviewPresenter) readTerminalSessions(c echo.Context) (
	responseDto dto.ReadTerminalSessionsResponse, err error,
) {
	operatorAccountId, assertOk := c.Get("operatorAccountId").(tkValueObject.AccountId)
	if !assertOk {
		return responseDto, errors.New("OperatorAccountIdNotFound")
	}

	readRequestBody := presenterHelper.ReadPaginationRequestParams(c)
	readRequestBody["operatorAccountId"] = operatorAccountId

	responseOutput := presenter.terminalSessionLiaison.Read(readRequestBody)
	if responseOutput.Status != tkPresentation.LiaisonResponseStatusSuccess {
		return responseDto, errors.New("ReadTerminalSessionsLiaisonBadResponse")
	}

	typedOutputBody, assertOk := responseOutput.Body.(dto.ReadTerminalSessionsResponse)
	if !assertOk {
		return responseDto, errors.New("AssertReadTerminalSessionsResponseFailed")
	}

	return typedOutputBody, nil
}

func (presenter *OverviewPresenter) TerminalSessionsTableHandler(c echo.Context) error {
	responseDto, err := presenter.readTerminalSessions(c)
	if err != nil {
		slog.Error("ReadTerminalSessionsError", slog.String("err", err.Error()))
		return c.NoContent(http.StatusInternalServerError)
	}

	return TerminalSessionsDataTable(responseDto).Render(c.Request().Context(), c.Response())
}

func (presenter *OverviewPresenter) ServicesTableHandler(c echo.Context) error {
	responseDto, err := presenter.readInstalledServices(c)
	if err != nil {
		slog.Error("ReadInstalledServicesError", slog.String("err", err.Error()))
		return c.NoContent(http.StatusInternalServerError)
	}

	return InstalledServicesDataTable(responseDto).Render(c.Request().Context(), c.Response())
}

func (presenter *OverviewPresenter) MarketplaceTableHandler(c echo.Context) error {
	marketplaceOverview, err := presenter.marketplacePresenter.MarketplaceOverviewFactory(
		"installed", presenterHelper.ReadPaginationRequestParams(c),
	)
	if err != nil {
		slog.Error("ReadInstalledMarketplaceItemsError", slog.String("err", err.Error()))
		return c.NoContent(http.StatusInternalServerError)
	}

	return InstalledMarketplaceItemsDataTable(marketplaceOverview).Render(
		c.Request().Context(), c.Response(),
	)
}

func (presenter *OverviewPresenter) Handler(echoContext echo.Context) error {
	vhostsHostnames, err := presenterHelper.ReadVirtualHostHostnames(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	if err != nil {
		slog.Error("ReadVirtualHostsHostnames", slog.String("err", err.Error()))
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	marketplaceOverview, err := presenter.marketplacePresenter.MarketplaceOverviewFactory(
		"all", presenterHelper.ReadPaginationRequestParams(echoContext),
	)
	if err != nil {
		slog.Error("MarketplaceOverviewFactoryError", slog.String("err", err.Error()))
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(presenter.transientDbSvc)
	o11yOverview, err := useCase.ReadO11yOverview(o11yQueryRepo, false)
	if err != nil {
		slog.Error("ReadO11yOverviewError", slog.String("err", err.Error()))
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	servicesOverview, err := presenter.servicesOverviewFactory(echoContext)
	if err != nil {
		slog.Error("ServicesOverviewFactoryError", slog.String("err", err.Error()))
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	terminalSessionsResponseDto, err := presenter.readTerminalSessions(echoContext)
	if err != nil {
		slog.Error("ReadTerminalSessionsError", slog.String("err", err.Error()))
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	pageContent := OverviewIndex(
		vhostsHostnames, marketplaceOverview, o11yOverview, servicesOverview,
		terminalSessionsResponseDto,
	)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  echoContext,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
