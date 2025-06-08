package uiPresenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	layoutFooter "github.com/goinfinite/os/src/presentation/ui/layout/footer"
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
	o11yLiaison := liaison.NewO11yLiaison(presenter.transientDbSvc)

	o11yLiaisonOutput := o11yLiaison.ReadOverview()
	if o11yLiaisonOutput.Status != liaison.Success {
		slog.Debug("FooterPresenterReadOverviewFailure")
		return nil
	}

	o11yOverviewEntity, assertOk := o11yLiaisonOutput.Body.(entity.O11yOverview)
	if !assertOk {
		slog.Debug("FooterPresenterAssertOverviewFailure")
		return nil
	}

	scheduledTaskLiaison := liaison.NewScheduledTaskLiaison(presenter.persistentDbSvc)

	scheduledTaskReadRequestBody := map[string]interface{}{
		"pageNumber":    0,
		"itemsPerPage":  5,
		"sortBy":        "id",
		"sortDirection": "desc",
	}
	scheduledTaskLiaisonOutput := scheduledTaskLiaison.Read(scheduledTaskReadRequestBody)
	if scheduledTaskLiaisonOutput.Status != liaison.Success {
		slog.Debug("FooterPresenterReadScheduledTaskFailure")
		return nil
	}

	tasksResponseDto, assertOk := scheduledTaskLiaisonOutput.Body.(dto.ReadScheduledTasksResponse)
	if !assertOk {
		slog.Debug("FooterPresenterAssertScheduledTaskResponseFailure")
		return nil
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layoutFooter.Footer(o11yOverviewEntity, tasksResponseDto.Tasks).
		Render(c.Request().Context(), c.Response().Writer)
}
