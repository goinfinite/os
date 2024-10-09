package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type SslsPresenter struct {
	sslService *service.SslService
}

func NewSslsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
) *SslsPresenter {
	return &SslsPresenter{
		sslService: service.NewSslService(persistentDbSvc, transientDbSvc),
	}
}

func (presenter *SslsPresenter) getVhostsHostnames(sslPairs []entity.SslPair) []string {
	vhostHostnames := []string{}
	for _, sslPair := range sslPairs {
		for _, vhostHostname := range sslPair.VirtualHostsHostnames {
			vhostHostnames = append(vhostHostnames, vhostHostname.String())
		}
	}

	return vhostHostnames
}

func (presenter *SslsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.sslService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	sslPairs, assertOk := responseOutput.Body.([]entity.SslPair)
	if !assertOk {
		return nil
	}

	pageContent := page.SslsIndex(sslPairs, presenter.getVhostsHostnames(sslPairs))
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
