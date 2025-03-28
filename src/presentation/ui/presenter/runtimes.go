package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	presenterDto "github.com/goinfinite/os/src/presentation/ui/presenter/dto"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type RuntimesPresenter struct {
	runtimeService  *service.RuntimeService
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewRuntimesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *RuntimesPresenter {
	return &RuntimesPresenter{
		runtimeService:  service.NewRuntimeService(persistentDbSvc, trailDbSvc),
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *RuntimesPresenter) runtimeOverviewFactory(
	runtimeType valueObject.RuntimeType,
	selectedVhostHostname valueObject.Fqdn,
) (runtimeOverview presenterDto.RuntimeOverview, err error) {
	isInstalled := false
	isVirtualHostUsingRuntime := false

	var phpConfigsPtr *entity.PhpConfigs
	if runtimeType.String() == "php" {
		requestBody := map[string]interface{}{"hostname": selectedVhostHostname.String()}
		responseOutput := presenter.runtimeService.ReadPhpConfigs(requestBody)

		isInstalled = true
		isVirtualHostUsingRuntime = true
		if responseOutput.Status != service.Success {
			isVirtualHostUsingRuntime = false
			responseOutputBodyStr, assertOk := responseOutput.Body.(string)
			if assertOk {
				isInstalled = responseOutputBodyStr != "ServiceUnavailable"
			}
		}

		if isInstalled {
			phpConfigs, assertOk := responseOutput.Body.(entity.PhpConfigs)
			if assertOk {
				phpConfigsPtr = &phpConfigs
			}
		}
	}

	return presenterDto.NewRuntimeOverview(
		selectedVhostHostname, runtimeType, isInstalled,
		isVirtualHostUsingRuntime, phpConfigsPtr,
	), nil
}

func (presenter *RuntimesPresenter) Handler(c echo.Context) error {
	rawRuntimeType := "php"
	if c.QueryParam("runtimeType") != "" {
		rawRuntimeType = c.QueryParam("runtimeType")
	}
	runtimeType, err := valueObject.NewRuntimeType(rawRuntimeType)
	if err != nil {
		slog.Error("InvalidRuntimeType", slog.String("err", err.Error()))
		return nil
	}

	primaryVhostHostname, err := infraHelper.ReadPrimaryVirtualHostHostname()
	if err != nil {
		slog.Error("ReadPrimaryVirtualHost", slog.String("err", err.Error()))
		return nil
	}
	selectedVhostHostname := primaryVhostHostname
	if c.QueryParam("vhostHostname") != "" {
		selectedVhostHostname, err = valueObject.NewFqdn(c.QueryParam("vhostHostname"))
		if err != nil {
			slog.Error("InvalidVhostHostname", slog.String("err", err.Error()))
			return nil
		}
	}

	runtimeOverview, err := presenter.runtimeOverviewFactory(
		runtimeType, selectedVhostHostname,
	)
	if err != nil {
		slog.Error("RuntimeOverviewFactoryError", slog.String("err", err.Error()))
		return nil
	}

	vhostsHostnames, err := presenterHelper.ReadVirtualHostHostnames(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	if err != nil {
		slog.Error("ReadVirtualHostsHostnames", slog.String("err", err.Error()))
		return nil
	}

	pageContent := page.RuntimesIndex(runtimeOverview, vhostsHostnames)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
