package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type SslsPresenter struct {
	sslService      *service.SslService
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewSslsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SslsPresenter {
	return &SslsPresenter{
		sslService:      service.NewSslService(persistentDbSvc, transientDbSvc, trailDbSvc),
		persistentDbSvc: persistentDbSvc,
		trailDbSvc:      trailDbSvc,
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

	vhostHostnames, err := presenterHelper.ReadVirtualHostHostnames(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	if err != nil {
		slog.Debug("ReadVirtualHostHostnamesError", "error", err)
		return nil
	}

	pageContent := page.SslsIndex(sslPairs, vhostHostnames)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
