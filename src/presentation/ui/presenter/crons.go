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

type CronsPresenter struct {
	cronService *service.CronService
}

func NewCronsPresenter(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) (presenter *CronsPresenter, err error) {
	cronService, err := service.NewCronService(trailDbSvc)
	if err != nil {
		return presenter, err
	}

	return &CronsPresenter{
		cronService: cronService,
	}, nil
}

func (presenter *CronsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.cronService.Read()
	if responseOutput.Status != service.Success {
		return nil
	}

	crons, assertOk := responseOutput.Body.([]entity.Cron)
	if !assertOk {
		return nil
	}

	pageContent := page.CronsIndex(crons)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
