package uiLayout

import (
	"errors"

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
	if componentSettings.EchoContext == nil {
		return errors.New("EchoContextIsMissing")
	}
	if componentSettings.PageContent == nil {
		return errors.New("PageContentIsMissing")
	}

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

	baseHref, assertOk := componentSettings.EchoContext.Get("baseHref").(string)
	if !assertOk {
		baseHref = "/"
	}

	return layoutMain.Main(layoutMain.MainLayoutSettings{
		PageContent: componentSettings.PageContent,
		BaseHref:    baseHref,
	}).Render(
		componentSettings.EchoContext.Request().Context(),
		componentSettings.EchoContext.Response().Writer,
	)
}
