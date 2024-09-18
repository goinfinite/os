package presenter

import (
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
	runtimeService *service.RuntimeService
}

func NewRuntimesPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimesPresenter {
	return &RuntimesPresenter{
		runtimeService: service.NewRuntimeService(persistentDbSvc),
	}
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

	pageContent := page.RuntimesIndex(selectedVhostHostname, isPhpInstalled, phpConfigs)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
