package uiPresenter

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	accountInfra "github.com/goinfinite/os/src/infra/account"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	"github.com/labstack/echo/v4"
)

type TerminalPresenter struct {
	accountQueryRepo *accountInfra.AccountQueryRepo
}

func NewTerminalPresenter(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *TerminalPresenter {
	return &TerminalPresenter{
		accountQueryRepo: accountInfra.NewAccountQueryRepo(persistentDbSvc),
	}
}

func (presenter *TerminalPresenter) readOperatorAccount(
	echoContext echo.Context,
) (accountEntity entity.Account, err error) {
	operatorAccountId, assertOk := echoContext.Get("operatorAccountId").(tkValueObject.AccountId)
	if !assertOk {
		return accountEntity, errors.New("OperatorAccountIdNotFound")
	}

	return presenter.accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &operatorAccountId,
	})
}

func (presenter *TerminalPresenter) FragmentHandler(echoContext echo.Context) error {
	accountEntity, err := presenter.readOperatorAccount(echoContext)
	if err != nil {
		slog.Error("ReadTerminalOperatorAccountError", slog.String("err", err.Error()))
		return echoContext.NoContent(http.StatusInternalServerError)
	}

	echoContext.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return TerminalSessionsContent(
		accountEntity.Username.String(), accountEntity.HomeDirectory.String(),
	).Render(echoContext.Request().Context(), echoContext.Response().Writer)
}
