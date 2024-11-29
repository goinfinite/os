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

type AccountsPresenter struct {
	accountService *service.AccountService
}

func NewAccountsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountsPresenter {
	return &AccountsPresenter{
		accountService: service.NewAccountService(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *AccountsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.accountService.Read(
		map[string]interface{}{
			"shouldIncludeSecureAccessKeys": true,
		},
	)
	if responseOutput.Status != service.Success {
		return nil
	}

	accounts, assertOk := responseOutput.Body.([]entity.Account)
	if !assertOk {
		return nil
	}

	pageContent := page.AccountsIndex(accounts)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
