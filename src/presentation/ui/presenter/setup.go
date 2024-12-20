package presenter

import (
	"github.com/labstack/echo/v4"
)

type SetupPresenter struct{}

func NewSetupPresenter() *SetupPresenter {
	return &SetupPresenter{}
}

func (presenter *SetupPresenter) Handler(c echo.Context) error {
	return nil
}
