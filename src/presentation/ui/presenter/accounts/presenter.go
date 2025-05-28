package uiPresenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/service"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
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
			"shouldIncludeSecureAccessPublicKeys": true,
		},
	)
	if responseOutput.Status != service.Success {
		return nil
	}

	typedOutputBody, assertOk := responseOutput.Body.(dto.ReadAccountsResponse)
	if !assertOk {
		return nil
	}

	pageContent := AccountsIndex(typedOutputBody.Accounts)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
