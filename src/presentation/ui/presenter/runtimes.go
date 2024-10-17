package presenter

import (
	"errors"
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
	"github.com/labstack/echo/v4"
)

type RuntimesPresenter struct {
	runtimeService     *service.RuntimeService
	virtualHostService *service.VirtualHostService
}

func NewRuntimesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimesPresenter {
	return &RuntimesPresenter{
		runtimeService:     service.NewRuntimeService(persistentDbSvc),
		virtualHostService: service.NewVirtualHostService(persistentDbSvc),
	}
}

func (presenter *RuntimesPresenter) readVhostsHostnames() ([]string, error) {
	vhostsHostnames := []string{}

	responseOutput := presenter.virtualHostService.Read()
	if responseOutput.Status != service.Success {
		responseBodyErrorStr := responseOutput.Body.(string)
		return vhostsHostnames, errors.New(responseBodyErrorStr)
	}

	existentVhosts, assertOk := responseOutput.Body.([]entity.VirtualHost)
	if !assertOk {
		return vhostsHostnames, errors.New(
			"InvalidVirtualHostsHostnamesStructure",
		)
	}

	for _, existentVhost := range existentVhosts {
		vhostsHostnames = append(vhostsHostnames, existentVhost.Hostname.String())
	}

	return vhostsHostnames, nil
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
		slog.Error("InvalidRuntimeType", slog.Any("error", err))
		return nil
	}

	primaryVhostHostname, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		slog.Error("ReadPrimaryVirtualHost", slog.Any("error", err))
		return nil
	}
	selectedVhostHostname := primaryVhostHostname
	if c.QueryParam("vhostHostname") != "" {
		selectedVhostHostname, err = valueObject.NewFqdn(c.QueryParam("vhostHostname"))
		if err != nil {
			slog.Error("InvalidVhostHostname", slog.Any("error", err))
			return nil
		}
	}

	runtimeOverview, err := presenter.runtimeOverviewFactory(
		runtimeType, selectedVhostHostname,
	)
	if err != nil {
		slog.Error("RuntimeOverviewFactoryError", slog.Any("error", err))
		return nil
	}

	vhostsHostnames, err := presenter.readVhostsHostnames()
	if err != nil {
		slog.Error("ReadVirtualHostsHostnames", slog.Any("error", err))
		return nil
	}

	pageContent := page.RuntimesIndex(runtimeOverview, vhostsHostnames)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
