package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	"github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type SetupPresenter struct {
	accountService *service.AccountService
}

func NewSetupPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *SetupPresenter {
	return &SetupPresenter{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *SetupPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.accountService.Read(map[string]interface{}{})
	if responseOutput.Status != service.Success {
		return nil
	}

	typedOutputBody, assertOk := responseOutput.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		return nil
	}

	if len(typedOutputBody.Accounts) > 0 {
		return c.Redirect(http.StatusFound, "/login/")
	}

	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layout.SetupLayout().
		Render(c.Request().Context(), c.Response().Writer)
}
