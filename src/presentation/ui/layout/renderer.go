package uiLayout

import (
	"github.com/a-h/templ"
	layoutMain "github.com/goinfinite/os/src/presentation/ui/layout/main"
	"github.com/labstack/echo/v4"
)

type LayoutRendererSettings struct {
	EchoContext     echo.Context
	PageContent     templ.Component
	ResponseCode    int
	PreferredLayout templ.Component
}

func Renderer(componentSettings LayoutRendererSettings) error {
	componentSettings.EchoContext.Response().
		Writer.WriteHeader(componentSettings.ResponseCode)
	componentSettings.EchoContext.Response().
		Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	if componentSettings.PreferredLayout != nil {
		return componentSettings.PreferredLayout.
			Render(
				componentSettings.EchoContext.Request().Context(),
				componentSettings.EchoContext.Response().Writer,
			)
	}

	return layoutMain.Main(layoutMain.MainLayoutSettings{
		PageContent: componentSettings.PageContent,
	}).Render(
		componentSettings.EchoContext.Request().Context(),
		componentSettings.EchoContext.Response().Writer,
	)
}
