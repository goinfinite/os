package presenter

import (
	"errors"
	"log/slog"
	"net/http"
	"os/user"

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

func (presenter *FileManagerPresenter) readAccountHomeDir(
	rawAccountId any,
) (accountHomeDir valueObject.UnixFilePath, err error) {
	accountId, err := valueObject.NewAccountId(rawAccountId)
	if err != nil {
		return accountHomeDir, err
	}

	accountIdStr := accountId.String()
	user, err := user.LookupId(accountIdStr)
	if err != nil {
		return accountHomeDir, errors.New("LookupAccountIdError: " + err.Error())
	}

	return valueObject.NewUnixFilePath(user.HomeDir)
}

func (presenter *FileManagerPresenter) readFilesGroupedByType(
	rawDesiredSourcePath string,
) (filesGroupedByType page.FilesGroupedByType, err error) {
	desiredSourcePath, err := valueObject.NewUnixFilePath(rawDesiredSourcePath)
	if err != nil {
		return filesGroupedByType, err
	}

	filesList, err := useCase.ReadFiles(filesInfra.FilesQueryRepo{}, desiredSourcePath)
	if err != nil {
		return filesGroupedByType, err
	}

	filesGroupedByType = page.FilesGroupedByType{
		SourcePath:  desiredSourcePath.String(),
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
	rawAccountId := c.Get("accountId")
	accountHomeDir, err := presenter.readAccountHomeDir(rawAccountId)
	if err != nil {
		slog.Debug(err.Error(), slog.Any("rawAccountId", rawAccountId))
		return nil
	}

	rawDesiredSourcePath := c.QueryParam("desiredSourcePath")
	if rawDesiredSourcePath == "" {
		rawDesiredSourcePath = accountHomeDir.String()
	}
	filesGroupedByType, err := presenter.readFilesGroupedByType(rawDesiredSourcePath)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	pageContent := page.FileManagerIndex(accountHomeDir, filesGroupedByType)
	return uiHelper.Render(c, pageContent, http.StatusOK)
}
