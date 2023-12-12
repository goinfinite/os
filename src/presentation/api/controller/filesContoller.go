package apiController

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

func getInodeNameByFilePath(filePath valueObject.UnixFilePath) string {
	fileIsDir, _ := filePath.IsDir()
	inodeName := "File"
	if fileIsDir {
		inodeName = "Directory"
	}

	return inodeName
}

func getFilePathSliceFromBody(
	filePathBodyInput interface{},
) []valueObject.UnixFilePath {
	var filePaths []valueObject.UnixFilePath

	filePathsIsList := reflect.TypeOf(filePathBodyInput).Kind() == reflect.Slice
	if !filePathsIsList {
		panic(errors.New("FilePathIsNotASlice"))
	}

	for _, filePathInterface := range filePathBodyInput.([]interface{}) {
		filePathStr := filePathInterface.(string)
		filePath, err := valueObject.NewUnixFilePath(filePathStr)
		if err != nil {
			continue
		}

		filePaths = append(filePaths, filePath)
	}

	return filePaths
}

// GetFiles    godoc
// @Summary      GetFiles
// @Description  List dir/files.
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
// @Description  Add a new dir/file.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addFileDto 	  body    dto.AddUnixFile  true  "NewFile"
// @Success      201 {object} object{} "FileCreated/DirectoryCreated"
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

// UpdateFile godoc
// @Summary      UpdateFile
// @Description  Update a dir/file path, name and/or permissions (ONly filePath is required).
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateUnixFileDto 	  body dto.UpdateUnixFile  true  "UpdateFile"
// @Success      200 {object} object{} "FileUpdated/DirectoryUpdate"
// @Router       /files/ [put]
func UpdateFileController(c echo.Context) error {
	requiredParams := []string{"filePath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePath := valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string))

	var destinationPathPtr *valueObject.UnixFilePath
	if requestBody["destinationPath"] != nil {
		destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))
		destinationPathPtr = &destinationPath
	}

	var permissionsPtr *valueObject.UnixFilePermissions
	if requestBody["permissions"] != nil {
		permissions := valueObject.NewUnixFilePermissionsPanic(requestBody["permissions"].(string))
		permissionsPtr = &permissions
	}

	updateUnixFileDto := dto.NewUpdateUnixFile(
		filePath,
		destinationPathPtr,
		permissionsPtr,
	)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.UpdateUnixFile(
		filesQueryRepo,
		filesCmdRepo,
		updateUnixFileDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	inodeName := getInodeNameByFilePath(filePath)

	return apiHelper.ResponseWrapper(c, http.StatusOK, inodeName+"Updated")
}

// UpdateFile godoc
// @Summary      UpdateFileContent
// @Description  Update a file content.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateUnixFileContentDto 	  body dto.UpdateUnixFileContent  true  "UpdateFileContent"
// @Success      200 {object} object{} "FileContentUpdated"
// @Router       /files/content/ [put]
func UpdateFileContentController(c echo.Context) error {
	requiredParams := []string{"filePath", "encodedFileContent"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePath := valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string))
	fileContent := valueObject.NewUnixFileContentPanic(requestBody["encodedFileContent"].(string))

	updateUnixFileContentDto := dto.NewUpdateUnixFileContent(filePath, fileContent)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.UpdateUnixFileContent(
		filesQueryRepo,
		filesCmdRepo,
		updateUnixFileContentDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "FileContentUpdated")
}

// AddFileCopy    godoc
// @Summary      AddFileCopy
// @Description  Add a new dir/file copy.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        addFileCopyDto 	  body    dto.AddUnixFileCopy  true  "NewFileCopy"
// @Success      201 {object} object{} "FileCopyCreated/DirectoryCopyCreated"
// @Router       /files/copy/ [post]
func AddFileCopyController(c echo.Context) error {
	requiredParams := []string{"filePath", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePath := valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string))
	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	addUnixFileCopyDto := dto.NewAddUnixFileCopy(filePath, destinationPath)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.AddUnixFileCopy(
		filesQueryRepo,
		filesCmdRepo,
		addUnixFileCopyDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	inodeName := getInodeNameByFilePath(filePath)

	return apiHelper.ResponseWrapper(c, http.StatusCreated, inodeName+"CopyCreated")
}

// DeleteFiles godoc
// @Summary      DeleteFiles
// @Description  Delete one or more directories/files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        deleteFilesDto 	  body    dto.DeleteUnixFiles  true  "DeleteFile"
// @Success      200 {object} object{} "DirectoriesAndFilesDeleted"
// @Router       /files/delete/ [put]
func DeleteFileController(c echo.Context) error {
	requiredParams := []string{"filePaths"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePaths := getFilePathSliceFromBody(requestBody["filePaths"])

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	useCase.DeleteUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		filePaths,
	)

	return apiHelper.ResponseWrapper(c, http.StatusOK, "DirectoriesAndFilesDeleted")
}
