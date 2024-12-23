package presenter

import (
	"net/http"

	"github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type SetupPresenter struct{}

func NewSetupPresenter() *SetupPresenter {
	return &SetupPresenter{}
}

func (presenter *SetupPresenter) Handler(c echo.Context) error {
	c.Response().Writer.WriteHeader(http.StatusOK)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return layout.Setup().
		Render(c.Request().Context(), c.Response().Writer)
}
