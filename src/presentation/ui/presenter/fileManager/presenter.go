package uiPresenter

import (
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	filesInfra "github.com/goinfinite/os/src/infra/files"
	uiLayout "github.com/goinfinite/os/src/presentation/ui/layout"
	"github.com/labstack/echo/v4"
)

type FileManagerPresenter struct{}

func NewFileManagerPresenter() *FileManagerPresenter {
	return &FileManagerPresenter{}
}

func (presenter *FileManagerPresenter) readUnixFilesByWorkingDir(
	workingDirPath tkValueObject.UnixAbsoluteFilePath,
) dto.ReadFilesResponse {
	shouldIncludeFileTree := true
	readFilesRequestDto := dto.ReadFilesRequest{
		SourcePath:            workingDirPath,
		ShouldIncludeFileTree: &shouldIncludeFileTree,
	}
	readFilesResponseDto, err := useCase.ReadFiles(
		&filesInfra.FilesQueryRepo{}, readFilesRequestDto,
	)
	if err != nil {
		return readFilesResponseDto
	}

	return readFilesResponseDto
}

func (presenter *FileManagerPresenter) Handler(c echo.Context) error {
	rawWorkingDirPath := c.QueryParam("workingDirPath")
	if rawWorkingDirPath == "" {
		rawWorkingDirPath = valueObject.UnixFilePathAppWorkingDir.String()
	}

	workingDirPath, err := tkValueObject.NewUnixAbsoluteFilePath(rawWorkingDirPath, false)
	if err != nil {
		workingDirPath, _ = tkValueObject.NewUnixAbsoluteFilePath("/invalid/path", false)
	}

	readFilesResponseDto := presenter.readUnixFilesByWorkingDir(workingDirPath)

	pageContent := FileManagerIndex(workingDirPath, readFilesResponseDto)
	return uiLayout.Renderer(uiLayout.LayoutRendererSettings{
		EchoContext:  c,
		PageContent:  pageContent,
		ResponseCode: http.StatusOK,
	})
}
