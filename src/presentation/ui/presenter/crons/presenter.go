package uiPresenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type CronsPresenter struct {
	cronService *service.CronService
}

func NewCronsPresenter(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *CronsPresenter {
	return &CronsPresenter{
		cronService: service.NewCronService(trailDbSvc),
	}
}

func (presenter *CronsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.cronService.Read(map[string]interface{}{})
	if responseOutput.Status != service.Success {
		return nil
	}

	typedOutputBody, assertOk := responseOutput.Body.(dto.ReadCronsResponse)
	if !assertOk {
		return nil
	}

	pageContent := CronsIndex(typedOutputBody.Crons)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
