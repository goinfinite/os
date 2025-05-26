package uiPresenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/goinfinite/os/src/presentation/ui/layout"
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
	sslPairsReadResponseServiceOutput := presenter.sslService.Read(map[string]interface{}{
		"itemsPerPage": 1000,
	})
	if sslPairsReadResponseServiceOutput.Status != service.Success {
		slog.Debug("SslPairsServiceBadOutput")
		return nil
	}

	sslPairsReadResponse, assertOk := sslPairsReadResponseServiceOutput.Body.(dto.ReadSslPairsResponse)
	if !assertOk {
		slog.Debug("ReadSslPairsResponseAssertionError")
		return nil
	}

	vhostHostnames, err := presenterHelper.ReadVirtualHostHostnames(
		presenter.persistentDbSvc, presenter.trailDbSvc,
	)
	if err != nil {
		slog.Debug("ReadVirtualHostHostnamesError", "error", err)
		return nil
	}

	pageContent := SslsIndex(sslPairsReadResponse.SslPairs, vhostHostnames)
	return layout.Renderer(layout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
