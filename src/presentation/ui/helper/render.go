package uiHelper

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"

	"github.com/goinfinite/os/src/presentation/ui/layout"
)

func Render(c echo.Context, pageContent templ.Component, statusCode int) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return layout.Main(pageContent).Render(c.Request().Context(), c.Response().Writer)
}
