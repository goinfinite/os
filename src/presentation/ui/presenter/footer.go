package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type FooterPresenter struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
	transientDbSvc  *internalDbInfra.TransientDatabaseService
	trailDbSvc      *internalDbInfra.TrailDatabaseService
}

func NewFooterPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	transientDbSvc *internalDbInfra.TransientDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *FooterPresenter {
	return &FooterPresenter{
		persistentDbSvc: persistentDbSvc,
		transientDbSvc:  transientDbSvc,
		trailDbSvc:      trailDbSvc,
	}
}

func (presenter *FooterPresenter) Handler(c echo.Context) error {
	o11yService := service.NewO11yService(presenter.transientDbSvc)

	o11yServiceOutput := o11yService.ReadOverview()
	if o11yServiceOutput.Status != service.Success {
		slog.Debug("FooterPresenterReadOverviewFailure")
		return nil
	}

	o11yOverviewEntity, assertOk := o11yServiceOutput.Body.(entity.O11yOverview)
	if !assertOk {
		slog.Debug("FooterPresenterAssertOverviewFailure")
		return nil
	}

	scheduledTaskService := service.NewScheduledTaskService(presenter.persistentDbSvc)

	scheduledTaskReadRequestBody := map[string]interface{}{
		"pageNumber":    0,
		"itemsPerPage":  5,
		"sortBy":        "id",
		"sortDirection": "desc",
	}
	scheduledTaskServiceOutput := scheduledTaskService.Read(scheduledTaskReadRequestBody)
	if scheduledTaskServiceOutput.Status != service.Success {
		slog.Debug("FooterPresenterReadScheduledTaskFailure")
		return nil
	}

	tasksResponseDto, assertOk := scheduledTaskServiceOutput.Body.(dto.ReadScheduledTasksResponse)
	if !assertOk {
		slog.Debug("FooterPresenterAssertScheduledTaskResponseFailure")
		return nil
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layout.Footer(o11yOverviewEntity, tasksResponseDto.Tasks).
		Render(c.Request().Context(), c.Response().Writer)
}
