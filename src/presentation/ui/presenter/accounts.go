package presenter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/entity"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	"github.com/speedianet/os/src/presentation/service"
	uiHelper "github.com/speedianet/os/src/presentation/ui/helper"
	"github.com/speedianet/os/src/presentation/ui/page"
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
	responseOutput := presenter.accountService.Read()
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
