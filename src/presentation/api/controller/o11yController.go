package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/useCase"
	o11yInfra "github.com/speedianet/os/src/infra/o11y"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// O11yOverview  godoc
// @Summary      O11yOverview
// @Description  Show system information and resource usage.
// @Tags         o11y
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200 {object} entity.O11yOverview
// @Router       /o11y/overview/ [get]
func O11yOverviewController(c echo.Context) error {
	o11yQueryRepo := o11yInfra.O11yQueryRepo{}
	o11yOverview, err := useCase.GetO11yOverview(o11yQueryRepo)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, o11yOverview)
}
