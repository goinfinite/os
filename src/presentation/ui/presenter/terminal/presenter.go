package uiPresenter

import (
	"github.com/labstack/echo/v4"
)

type TerminalPresenter struct{}

func NewTerminalPresenter() *TerminalPresenter {
	return &TerminalPresenter{}
}

func (presenter *TerminalPresenter) FragmentHandler(echoContext echo.Context) error {
	echoContext.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return TerminalSessionsContent().Render(
		echoContext.Request().Context(), echoContext.Response().Writer,
	)
}
