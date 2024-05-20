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
	filesInfra "github.com/speedianet/os/src/infra/files"
	apiHelper "github.com/speedianet/os/src/presentation/api/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
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
// @Param        sourcePath	query	string	true	"SourcePath"
// @Success      200 {array} entity.UnixFile
// @Router       /files/ [get]
func GetFilesController(c echo.Context) error {
	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesList, err := useCase.GetFiles(
		filesQueryRepo,
		valueObject.NewUnixFilePathPanic(c.QueryParam("sourcePath")),
	)
	if err != nil {
		if err.Error() == "DirPathCannotEndWithSlash" {
			return apiHelper.ResponseWrapper(c, http.StatusNotFound, err.Error())
		}

		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, filesList)
}

// CreateFile    godoc
// @Summary      CreateNewFile
// @Description  Create a new dir/file.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        createFileDto 	  body    dto.CreateUnixFile  true  "NewFile"
// @Success      201 {object} object{} "FileCreated/DirectoryCreated"
// @Router       /files/ [post]
func CreateFileController(c echo.Context) error {
	requiredParams := []string{"filePath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	fileType := valueObject.NewMimeTypePanic("generic")
	isDirType := false
	if requestBody["mimeType"] != nil {
		fileTypeStr := requestBody["mimeType"].(string)
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

	createUnixFileDto := dto.NewCreateUnixFile(
		valueObject.NewUnixFilePathPanic(requestBody["filePath"].(string)),
		&filePermissions,
		fileType,
	)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err := useCase.CreateUnixFile(
		filesQueryRepo,
		filesCmdRepo,
		createUnixFileDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, successResponse)
}

// UpdateFile godoc
// @Summary      UpdateFile
// @Description  Move a dir/file, update name and/or permissions (Only sourcePath is required).
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        updateUnixFilesDto 	  body dto.UpdateUnixFiles  true  "Only sourcePaths are required."
// @Success      200 {object} object{} "FileUpdated"
// @Success      207 {object} object{} "FilesArePartialUpdated"
// @Router       /files/ [put]
func UpdateFileController(c echo.Context) error {
	requiredParams := []string{"sourcePaths"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePaths := getFilePathSliceFromBody(requestBody["sourcePaths"])

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

	updateUnixFileDto := dto.NewUpdateUnixFiles(
		sourcePaths,
		destinationPathPtr,
		permissionsPtr,
		encodedContentPtr,
	)

	filesCmdRepo := filesInfra.FilesCmdRepo{}

	updateUnixFileUc := useCase.NewUpdateUnixFiles(filesCmdRepo)
	updateProcessInfo, err := updateUnixFileUc.Execute(updateUnixFileDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	httpStatus := http.StatusOK

	hasSuccess := len(updateProcessInfo.FilePathsSuccessfullyUpdated) > 0
	hasFailures := len(updateProcessInfo.FailedPathsWithReason) > 0
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		httpStatus = http.StatusMultiStatus
	}

	return apiHelper.ResponseWrapper(c, httpStatus, updateProcessInfo)
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
	requiredParams := []string{"sourcePath", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePath := valueObject.NewUnixFilePathPanic(requestBody["sourcePath"].(string))
	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	shouldOverwrite := false
	if requestBody["shouldOverwrite"] != nil {
		var err error
		shouldOverwrite, err = sharedHelper.ParseBoolParam(
			requestBody["shouldOverwrite"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidShouldOverwrite",
			)
		}
	}

	copyUnixFileDto := dto.NewCopyUnixFile(sourcePath, destinationPath, shouldOverwrite)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

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
// @Param        sourcePaths	body	[]string	true	"SourcePath"
// @Success      200 {object} object{} "FilesDeleted"
// @Router       /files/delete/ [put]
func DeleteFileController(c echo.Context) error {
	requiredParams := []string{"sourcePaths"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePaths := getFilePathSliceFromBody(requestBody["sourcePaths"])

	hardDelete := false
	if requestBody["hardDelete"] != nil {
		var err error
		hardDelete, err = sharedHelper.ParseBoolParam(
			requestBody["hardDelete"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidHardDelete",
			)
		}
	}

	deleteUnixFilesDto := dto.NewDeleteUnixFiles(
		sourcePaths,
		hardDelete,
	)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	deleteUnixFiles := useCase.NewDeleteUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
	)

	err := deleteUnixFiles.Execute(deleteUnixFilesDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, "FilesDeleted")
}

// CompressFiles    godoc
// @Summary      CompressFiles
// @Description  Compress directories and files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        compressFilesDto 	  body    dto.CompressUnixFiles  true  "CompressFiles"
// @Success      200 {object} object{} "FilesCompressed"
// @Success      207 {object} object{} "FilesArePartialCompressed"
// @Router       /files/compress/ [post]
func CompressFilesController(c echo.Context) error {
	requiredParams := []string{"sourcePaths", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePaths := getFilePathSliceFromBody(requestBody["sourcePaths"])

	var compressionUnixTypePtr *valueObject.UnixCompressionType
	if requestBody["compressionType"] != nil {
		compressionUnixType := valueObject.NewUnixCompressionTypePanic(requestBody["compressionType"].(string))
		compressionUnixTypePtr = &compressionUnixType
	}

	compressUnixFilesDto := dto.NewCompressUnixFiles(
		sourcePaths,
		valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string)),
		compressionUnixTypePtr,
	)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	compressionProcessInfo, err := useCase.CompressUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		compressUnixFilesDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	httpStatus := http.StatusCreated

	hasSuccess := len(compressionProcessInfo.FilePathsSuccessfullyCompressed) > 0
	hasFailures := len(compressionProcessInfo.FailedPathsWithReason) > 0
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
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
// @Success      200 {object} object{} "FilesExtracted"
// @Router       /files/extract/ [put]
func ExtractFilesController(c echo.Context) error {
	requiredParams := []string{"sourcePath", "destinationPath"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePath := valueObject.NewUnixFilePathPanic(requestBody["sourcePath"].(string))
	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	extractUnixFilesDto := dto.NewExtractUnixFiles(sourcePath, destinationPath)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err := useCase.ExtractUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		extractUnixFilesDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, "FilesExtracted")
}

// UploadFiles    godoc
// @Summary      UploadFiles
// @Description  Upload files.
// @Tags         files
// @Accept       mpfd
// @Produce      json
// @Security     Bearer
// @Param        destinationPath	formData	string	true	"DestinationPath"
// @Param        file	formData	file	true	"FileToUpload"
// @Success      200 {object} object{} "FilesUploaded"
// @Success      207 {object} object{} "FilesPartialUploaded"
// @Router       /files/upload/ [post]
func UploadFilesController(c echo.Context) error {
	requiredParams := []string{"destinationPath", "files"}
	requestBody, _ := apiHelper.GetRequestBody(c)

	if requestBody["destinationPath"] == nil {
		requestBody["destinationPath"] = c.QueryParam("destinationPath")
	}

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	destinationPath := valueObject.NewUnixFilePathPanic(requestBody["destinationPath"].(string))

	var filesToUpload []valueObject.FileStreamHandler
	for _, requestBodyFile := range requestBody["files"].(map[string]*multipart.FileHeader) {
		fileStreamHandler := valueObject.NewFileStreamHandlerPanic(requestBodyFile)
		filesToUpload = append(filesToUpload, fileStreamHandler)
	}

	uploadUnixFilesDto := dto.NewUploadUnixFiles(destinationPath, filesToUpload)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	uploadProcessInfo, err := useCase.UploadUnixFiles(
		filesQueryRepo,
		filesCmdRepo,
		uploadUnixFilesDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	httpStatus := http.StatusCreated

	hasSuccess := len(uploadProcessInfo.FileNamesSuccessfullyUploaded) > 0
	hasFailures := len(uploadProcessInfo.FailedNamesWithReason) > 0
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		httpStatus = http.StatusMultiStatus
	}

	return apiHelper.ResponseWrapper(c, httpStatus, uploadProcessInfo)
}
