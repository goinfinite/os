package presenter

import (
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
	rawWorkingDirPath string,
) page.FilesGroupedByType {
	workingDirPath, err := valueObject.NewUnixFilePath(rawWorkingDirPath)
	if err != nil {
		workingDirPath, _ = valueObject.NewUnixFilePath("/invalid/path")
	}

	filesGroupedByType := page.FilesGroupedByType{
		WorkingDirPath: workingDirPath.String(),
		Directories:    []entity.UnixFile{},
		Files:          []entity.UnixFile{},
	}
	filesList, err := useCase.ReadFiles(filesInfra.FilesQueryRepo{}, workingDirPath)
	if err != nil {
		return filesGroupedByType
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

	return filesGroupedByType
}

func (presenter *FileManagerPresenter) Handler(c echo.Context) error {
	rawWorkingDirPath := c.QueryParam("workingDirPath")
	if rawWorkingDirPath == "" {
		rawWorkingDirPath = valueObject.DefaultAppWorkingDir.String()
	}
	filesGroupedByType := presenter.readFilesGroupedByType(rawWorkingDirPath)

	pageContent := page.FileManagerIndex(filesGroupedByType)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
