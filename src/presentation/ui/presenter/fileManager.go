package presenter

import (
	"net/http"

	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type FileManagerPresenter struct{}

func NewFileManagerPresenter() *FileManagerPresenter {
	return &FileManagerPresenter{}
}

func (presenter *FileManagerPresenter) Handler(c echo.Context) error {
	pageContent := page.FileManagerIndex()
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
