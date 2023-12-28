package apiController

import (
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	"github.com/speedianet/os/src/infra"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
)

func getFilePathSliceFromBody(
	filePathBodyInput interface{},
) []valueObject.UnixFilePath {
	var filePaths []valueObject.UnixFilePath

	filePathBodyInputType := reflect.TypeOf(filePathBodyInput).Kind()

	switch filePathBodyInputType {
	case reflect.String:
		filePaths = append(
			filePaths,
			valueObject.NewUnixFilePathPanic(filePathBodyInput.(string)),
		)
	case reflect.Slice:
		for _, filePathInterface := range filePathBodyInput.([]interface{}) {
			filePathStr := filePathInterface.(string)
			filePath, err := valueObject.NewUnixFilePath(filePathStr)
			if err != nil {
				continue
			}

			filePaths = append(filePaths, filePath)
		}
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
	requiredParams := []string{"filePath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	fileType := valueObject.NewMimeTypePanic("generic")
	isDirType := false
	if requestBody["type"] != nil {
		fileTypeStr := requestBody["type"].(string)
		isDirType = strings.ToLower(fileTypeStr) == "directory"
	}

	if isDirType {
		fileType = valueObject.NewMimeTypePanic("directory")
	}

	successResponse := "FileCreated"

	filePermissions := valueObject.NewUnixFilePermissionsPanic("0644")
	if fileType.IsDir() {
		filePermissions = valueObject.NewUnixFilePermissionsPanic("0755")
		successResponse = "DirectoryCreated"
	}

	if requestBody["permissions"] != nil {
		filePermissions = valueObject.NewUnixFilePermissionsPanic(requestBody["permissions"].(string))
	}

	addUnixFileDto := dto.NewCreateUnixFile(
		valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string)),
		filePermissions,
		fileType,
	)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.CreateUnixFile(
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
// @Description  Update a dir/file path, name and/or permissions (Only filePath is required).
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateUnixFileDto 	  body dto.UpdateUnixFile  true  "UpdateFile"
// @Success      200 {object} object{} "FileUpdated"
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

	var encodedContentPtr *valueObject.EncodedContent
	if requestBody["encodedContent"] != nil {
		encodedContent := valueObject.NewEncodedContentPanic(requestBody["encodedContent"].(string))
		encodedContentPtr = &encodedContent
	}

	updateUnixFileDto := dto.NewUpdateUnixFile(
		filePath,
		destinationPathPtr,
		permissionsPtr,
		encodedContentPtr,
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

	return apiHelper.ResponseWrapper(c, http.StatusOK, "FileUpdated")
}

// CopyFile    godoc
// @Summary      CopyFile
// @Description  Copy a dir/file.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        copyFileDto 	  body    dto.CopyUnixFile  true  "NewFileCopy"
// @Success      201 {object} object{} "FileCopied"
// @Router       /files/copy/ [post]
func CopyFileController(c echo.Context) error {
	requiredParams := []string{"filePath", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePath := valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string))
	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	copyUnixFileDto := dto.NewCopyUnixFile(filePath, destinationPath)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.CopyUnixFile(
		filesQueryRepo,
		filesCmdRepo,
		copyUnixFileDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "FileCopied")
}

// DeleteFiles godoc
// @Summary      DeleteFiles
// @Description  Delete one or more directories/files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        filePath	body	[]string	true	"UnixFilePath"
// @Success      200 {object} object{} "DirectoriesAndFilesDeleted"
// @Router       /files/delete/ [put]
func DeleteFileController(c echo.Context) error {
	requiredParams := []string{"filePath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePaths := getFilePathSliceFromBody(requestBody["filePath"])

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	useCase.DeleteUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		filePaths,
	)

	return apiHelper.ResponseWrapper(c, http.StatusOK, "DirectoriesAndFilesDeleted")
}

// CompressFiles    godoc
// @Summary      CompressFiles
// @Description  Compress directories and files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        compressFilesDto 	  body    dto.CompressUnixFiles  true  "CompressFiles"
// @Success      200 {object} object{} "FilesAndDirectoriesCompressed"
// @Success      207 {object} object{} "FilesAndDirectoriesArePartialCompressed"
// @Router       /files/compress/ [post]
func CompressFilesController(c echo.Context) error {
	requiredParams := []string{"filePath", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePaths := getFilePathSliceFromBody(requestBody["filePath"])

	var compressionUnixTypePtr *valueObject.UnixCompressionType
	if requestBody["compressionType"] != nil {
		compressionUnixType := valueObject.NewUnixCompressionTypePanic(requestBody["compressionType"].(string))
		compressionUnixTypePtr = &compressionUnixType
	}

	compressUnixFilesDto := dto.NewCompressUnixFiles(
		filePaths,
		valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string)),
		compressionUnixTypePtr,
	)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	compressionProcessInfo, err := useCase.CompressUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		compressUnixFilesDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	httpStatus := http.StatusCreated

	hasFilePathsSuccessfullyCompressed := len(compressionProcessInfo.FilePathsSuccessfullyCompressed) > 0
	hasFailedPathsWithReason := len(compressionProcessInfo.FailedPathsWithReason) > 0
	isMultiStatus := hasFilePathsSuccessfullyCompressed && hasFailedPathsWithReason
	if isMultiStatus {
		httpStatus = http.StatusMultiStatus
	}

	return apiHelper.ResponseWrapper(c, httpStatus, compressionProcessInfo)
}

// ExtractFiles godoc
// @Summary      ExtractFiles
// @Description  Extract directories and files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        extractFilesDto 	  body    dto.ExtractUnixFiles  true  "ExtractFiles"
// @Success      200 {object} object{} "ExtractFilesAndDirectories"
// @Router       /files/extract/ [put]
func ExtractFilesController(c echo.Context) error {
	requiredParams := []string{"filePath", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePath := valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string))
	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	extractUnixFilesDto := dto.NewExtractUnixFiles(filePath, destinationPath)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	err := useCase.ExtractUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		extractUnixFilesDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "ExtractFilesAndDirectories")
}

// UploadFiles    godoc
// @Summary      UploadFiles
// @Description  Upload files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        destinationPath	path	string	true	"DestinationPath"
// @Param        file	formData	file	true	"FileToUpload"
// @Success      200 {object} object{} "FilesUploaded"
// @Success      207 {object} object{} "FilesPartialUploaded"
// @Router       /files/upload/ [post]
func UploadFilesController(c echo.Context) error {
	requiredParams := []string{"destinationPath", "files"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	var filesToUpload []valueObject.FileStreamHandler
	for _, requestBodyFile := range requestBody["files"].(map[string]*multipart.FileHeader) {
		fileStreamHandler := valueObject.NewFileStreamHandlerPanic(requestBodyFile)
		filesToUpload = append(filesToUpload, fileStreamHandler)
	}

	uploadUnixFilesDto := dto.NewUploadUnixFiles(destinationPath, filesToUpload)

	filesQueryRepo := infra.FilesQueryRepo{}
	filesCmdRepo := infra.FilesCmdRepo{}

	uploadProcessInfo, err := useCase.UploadUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		uploadUnixFilesDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	httpStatus := http.StatusCreated

	hasFilePathsSuccessfullyUploaded := len(uploadProcessInfo.FilePathsSuccessfullyUploaded) > 0
	hasFailedPathsWithReason := len(uploadProcessInfo.FailedPathsWithReason) > 0
	isMultiStatus := hasFilePathsSuccessfullyUploaded && hasFailedPathsWithReason
	if isMultiStatus {
		httpStatus = http.StatusMultiStatus
	}

	return apiHelper.ResponseWrapper(c, httpStatus, uploadProcessInfo)
}
