package uiPresenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	"github.com/goinfinite/os/src/presentation/liaison"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type AccountsPresenter struct {
	accountLiaison *liaison.AccountLiaison
}

func NewAccountsPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *AccountsPresenter {
	return &AccountsPresenter{
		accountLiaison: liaison.NewAccountLiaison(persistentDbSvc, trailDbSvc),
	}
}

func (presenter *AccountsPresenter) Handler(c echo.Context) error {
	responseOutput := presenter.accountLiaison.Read(
		map[string]interface{}{
			"shouldIncludeSecureAccessPublicKeys": true,
		},
	)
	if responseOutput.Status != liaison.Success {
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
