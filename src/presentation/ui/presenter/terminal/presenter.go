package uiPresenter

import (
	"net/http"

	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
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

func (presenter *TerminalPresenter) Handler(echoContext echo.Context) error {
	pageContent := TerminalIndex()
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  echoContext,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
