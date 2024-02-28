package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/useCase"
	databaseInfra "github.com/speedianet/os/src/infra/database"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

type O11yController struct {
	transientDbSvc *databaseInfra.TransientDatabaseService
}

func NewO11yController(
	transientDbSvc *databaseInfra.TransientDatabaseService,
) *O11yController {
	return &O11yController{
		transientDbSvc: transientDbSvc,
	}
}

// O11yOverview  godoc
// @Summary      O11yOverview
// @Description  Show system information and resource usage.
// @Tags         o11y
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.O11yOverview
// @Router       /o11y/overview/ [get]
func (controller O11yController) O11yOverviewController(c echo.Context) error {
	o11yQueryRepo := o11yInfra.NewO11yQueryRepo(controller.transientDbSvc)
	o11yOverview, err := useCase.GetO11yOverview(o11yQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, o11yOverview)
}
