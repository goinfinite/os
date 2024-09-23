package presenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation/service"
	uiHelper "github.com/speedianet/os/src/presentation/ui/helper"
	"github.com/speedianet/os/src/presentation/ui/page"
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

func (presenter *RuntimesPresenter) getVhostsHostnames() ([]string, error) {
	vhostsHostnames := []string{}

	responseOutput := presenter.virtualHostService.Read()
	if responseOutput.Status != service.Success {
		responseBodyErrorStr := responseOutput.Body.(string)
		return vhostsHostnames, errors.New(responseBodyErrorStr)
	}

	existentVhosts, assertOk := responseOutput.Body.([]entity.VirtualHost)
	if !assertOk {
		return vhostsHostnames, errors.New(
			"UnableToAssertExistentVirtualHostsHostnamesStructure",
		)
	}

	for _, existentVhost := range existentVhosts {
		vhostsHostnames = append(vhostsHostnames, existentVhost.Hostname.String())
	}

	return vhostsHostnames, nil
}

func (presenter *RuntimesPresenter) Handler(c echo.Context) error {
	selectedVhostHostname, err := valueObject.NewFqdn(c.QueryParam("vhostHostname"))
	if err != nil {
		primaryVhostHostname, err := infraHelper.GetPrimaryVirtualHost()
		if err != nil {
			return nil
		}
		selectedVhostHostname = primaryVhostHostname
	}

	isPhpInstalled := false

	requestBody := map[string]interface{}{"hostname": selectedVhostHostname.String()}
	responseOutput := presenter.runtimeService.ReadPhpConfigs(requestBody)
	if responseOutput.Status == service.Success {
		isPhpInstalled = true
	}

	phpConfigs, assertOk := responseOutput.Body.(entity.PhpConfigs)
	if !assertOk {
		isPhpInstalled = false
	}

	existentVhostsHostnames, err := presenter.getVhostsHostnames()
	if err != nil {
		slog.Error("GetExistentVirtualHostsHostnames", slog.Any("err", err))
		return nil
	}

	pageContent := page.RuntimesIndex(
		isPhpInstalled, phpConfigs, existentVhostsHostnames,
	)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
