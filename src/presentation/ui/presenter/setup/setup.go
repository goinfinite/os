package uiPresenter

import (
	"net/http"

	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	layoutSetup "github.com/goinfinite/os/src/presentation/ui/layout/setup"
	presenterHelper "github.com/goinfinite/os/src/presentation/ui/presenter/helper"
	"github.com/labstack/echo/v4"
)

type SetupPresenter struct {
	accountLiaison *liaison.AccountLiaison
}

func NewSetupPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SetupPresenter {
	return &SetupPresenter{
		accountLiaison: liaison.NewAccountLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *SetupPresenter) Handler(c echo.Context) error {
	if !presenterHelper.ShouldEnableInitialSetup(presenter.accountLiaison) {
		return c.Redirect(http.StatusFound, "/login/")
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layoutSetup.Setup().
		Render(c.Request().Context(), c.Response().Writer)
}
