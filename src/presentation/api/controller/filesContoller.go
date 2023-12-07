package apiController

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

// GetFiles    godoc
// @Summary      GetFiles
// @Description  List files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        unixFilePath	query	string	true	"UnixFilePath"
// @Success      200 {array} entity.UnixFile
// @Router       /files/ [get]
func GetFilesController(c echo.Context) error {
	filesQueryRepo := infra.FilesQueryRepo{}
	filesList, err := useCase.GetFiles(
		filesQueryRepo,
		valueObject.NewUnixFilePathPanic(c.QueryParam("unixFilePath")),
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, filesList)
}

// AddFile    godoc
// @Summary      AddNewFile
// @Description  Add a new file.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addFileDto 	  body    dto.AddUnixFile  true  "NewFile"
// @Success      201 {object} object{} "FileCreated"
// @Router       /files/ [post]
func AddFileController(c echo.Context) error {
	requiredParams := []string{"filePath", "type"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	fileType := valueObject.NewUnixFileTypePanic(requestBody["type"].(string))

	successResponse := "FileCreated"

	filePermissions := valueObject.NewUnixFilePermissionsPanic("0644")
	if fileType.IsDir() {
		filePermissions = valueObject.NewUnixFilePermissionsPanic("0755")
		successResponse = "DirectoryCreated"
	}

	if requestBody["permissions"] != nil {
		filePermissions = valueObject.NewUnixFilePermissionsPanic(requestBody["permissions"].(string))
	}

	addUnixFileDto := dto.NewAddUnixFile(
		valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string)),
		filePermissions,
		fileType,
	)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.AddUnixFile(
		filesQueryRepo,
		filesCmdRepo,
		addUnixFileDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, successResponse)
}
