package presenter

import (
	"log/slog"
	"net/http"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	filesInfra "github.com/goinfinite/os/src/infra/files"
	uiHelper "github.com/goinfinite/os/src/presentation/ui/helper"
	"github.com/goinfinite/os/src/presentation/ui/page"
	"github.com/labstack/echo/v4"
)

type FileManagerPresenter struct{}

func NewFileManagerPresenter() *FileManagerPresenter {
	return &FileManagerPresenter{}
}

func (presenter *FileManagerPresenter) readFilesGroupedByType(
	rawDesiredSourcePath string,
) (filesGroupedByType page.FilesGroupedByType, err error) {
	if rawDesiredSourcePath == "" {
		rawDesiredSourcePath = "/"
	}

	desiredSourcePath, err := valueObject.NewUnixFilePath(rawDesiredSourcePath)
	if err != nil {
		return filesGroupedByType, err
	}

	filesList, err := useCase.ReadFiles(filesInfra.FilesQueryRepo{}, desiredSourcePath)
	if err != nil {
		return filesGroupedByType, err
	}

	filesGroupedByType = page.FilesGroupedByType{
		Directories: []entity.UnixFile{},
		Files:       []entity.UnixFile{},
	}
	for _, fileEntity := range filesList {
		if fileEntity.MimeType.IsDir() {
			filesGroupedByType.Directories = append(
				filesGroupedByType.Directories, fileEntity,
			)
			continue
		}

		filesGroupedByType.Files = append(filesGroupedByType.Files, fileEntity)
	}

	return filesGroupedByType, nil
}

func (presenter *FileManagerPresenter) Handler(c echo.Context) error {
	filesGroupedByType, err := presenter.readFilesGroupedByType(
		c.QueryParam("desiredSourcePath"),
	)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.FileManagerIndex(filesGroupedByType)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
