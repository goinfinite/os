package presenter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/entity"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation/service"
	uiHelper "github.com/speedianet/os/src/presentation/ui/helper"
	"github.com/speedianet/os/src/presentation/ui/page"
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

func (presenter *SslsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.sslService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	sslPairs, assertOk := responseOutput.Body.([]entity.SslPair)
	if !assertOk {
		return nil
	}

	pageContent := page.SslsIndex(sslPairs, uiHelper.GetVhostHostnames(sslPairs))
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
