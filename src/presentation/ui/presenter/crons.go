package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/goinfinite/os/src/presentation/ui/page"
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

	pageContent := page.CronsIndex(typedOutputBody.Crons)
	return layout.Renderer(layout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
