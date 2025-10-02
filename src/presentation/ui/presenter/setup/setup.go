package uiPresenter

import (
	"log/slog"
	"net/http"
	"strings"

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

func (presenter *SetupPresenter) Handler(echoContext echo.Context) error {
	uiBasePath, assertOk := echoContext.Get("uiBasePath").(string)
	if !assertOk {
		slog.Error("AssertUiBasePathFailed")
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	baseHref, assertOk := echoContext.Get("baseHref").(string)
	if !assertOk {
		slog.Error("AssertBaseHrefFailed")
		return echoContext.NoContent(http.StatusInternalServerError)
	}
	if len(baseHref) > 0 {
		baseHrefNoTrailing := strings.TrimSuffix(baseHref, "/")
		uiBasePath = baseHrefNoTrailing + uiBasePath
	}

	if !presenterHelper.ShouldEnableInitialSetup(presenter.accountLiaison) {
		return echoContext.Redirect(http.StatusFound, uiBasePath+"/login/")
	}

	setupLayoutSettings := layoutSetup.SetupLayoutSettings{BaseHref: baseHref}

	echoContext.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	echoContext.Response().Writer.WriteHeader(http.StatusOK)

	return layoutSetup.Setup(setupLayoutSettings).
		Render(echoContext.Request().Context(), echoContext.Response().Writer)
}
