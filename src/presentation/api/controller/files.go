package apiController

import (
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	filesInfra "github.com/goinfinite/os/src/infra/files"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
)

type FilesController struct{}

func (controller *FilesController) parseSourcePaths(
	rawSourcePaths []interface{},
) ([]valueObject.UnixFilePath, error) {
	filePaths := []valueObject.UnixFilePath{}

	for pathIndex, rawSourcePath := range rawSourcePaths {
		filePath, err := valueObject.NewUnixFilePath(rawSourcePath)
		if err != nil {
			log.Printf("[PathIndex %d] %s", pathIndex, err.Error())
			continue
		}

		filePaths = append(filePaths, filePath)
	}

	return filePaths, nil
}

func NewFilesController() *FilesController {
	return &FilesController{}
}

// ReadFiles    godoc
// @Summary      ReadFiles
// @Description  List dir/files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        sourcePath	query	string	true	"SourcePath"
// @Success      200 {array} entity.UnixFile
// @Router       /v1/files/ [get]
func (controller *FilesController) Read(c echo.Context) error {
	filesQueryRepo := filesInfra.FilesQueryRepo{}

	sourcePath, err := valueObject.NewUnixFilePath(c.QueryParam("sourcePath"))
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	filesList, err := useCase.ReadFiles(filesQueryRepo, sourcePath)
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
// @Param        createFileDto 	  body    dto.CreateUnixFile  true  "permissions is optional. When not provided, permissions will be '644' for files and '755' for directories."
// @Success      201 {object} object{} "FileCreated/DirectoryCreated"
// @Router       /v1/files/ [post]
func (controller *FilesController) Create(c echo.Context) error {
	requiredParams := []string{"filePath"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	filePath, err := valueObject.NewUnixFilePath(requestBody["filePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	fileType, _ := valueObject.NewMimeType("generic")
	isDirType := false
	if requestBody["mimeType"] != nil {
		fileTypeStr := requestBody["mimeType"].(string)
		isDirType = strings.ToLower(fileTypeStr) == "directory"
	}

	if isDirType {
		fileType, _ = valueObject.NewMimeType("directory")
	}

	successResponse := "FileCreated"

	filePermissions, _ := valueObject.NewUnixFilePermissions("0644")
	if fileType.IsDir() {
		filePermissions, _ = valueObject.NewUnixFilePermissions("0755")
		successResponse = "DirectoryCreated"
	}

	if requestBody["permissions"] != nil {
		filePermissions, err = valueObject.NewUnixFilePermissions(
			requestBody["permissions"].(string),
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}

	createUnixFileDto := dto.NewCreateUnixFile(
		filePath,
		&filePermissions,
		fileType,
	)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err = useCase.CreateUnixFile(
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
// @Router       /v1/files/ [put]
func (controller *FilesController) Update(c echo.Context) error {
	requiredParams := []string{"sourcePaths"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	_, isSourcePathsMap := requestBody["sourcePaths"].(map[string]interface{})
	if isSourcePathsMap {
		requestBody["sourcePaths"] = []interface{}{requestBody["sourcePaths"]}
	}

	sourcePathsSlice, assertOk := requestBody["sourcePaths"].([]interface{})
	if !assertOk {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "SourcePathMustBeArray")
	}

	sourcePaths, err := controller.parseSourcePaths(sourcePathsSlice)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	var destinationPathPtr *valueObject.UnixFilePath
	if requestBody["destinationPath"] != nil {
		destinationPath, err := valueObject.NewUnixFilePath(requestBody["destinationPath"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		destinationPathPtr = &destinationPath
	}

	var permissionsPtr *valueObject.UnixFilePermissions
	if requestBody["permissions"] != nil {
		permissions, err := valueObject.NewUnixFilePermissions(requestBody["permissions"].(string))
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
		permissionsPtr = &permissions
	}

	var encodedContentPtr *valueObject.EncodedContent
	if requestBody["encodedContent"] != nil {
		encodedContent, err := valueObject.NewEncodedContent(requestBody["encodedContent"].(string))
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
		}
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
// @Param        copyFileDto 	  body    dto.CopyUnixFile  true  "All props are required."
// @Success      201 {object} object{} "FileCopied"
// @Router       /v1/files/copy/ [post]
func (controller *FilesController) Copy(c echo.Context) error {
	requiredParams := []string{"sourcePath", "destinationPath"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePath, err := valueObject.NewUnixFilePath(requestBody["sourcePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	destinationPath, err := valueObject.NewUnixFilePath(requestBody["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	shouldOverwrite := false
	if requestBody["shouldOverwrite"] != nil {
		var err error
		shouldOverwrite, err = voHelper.InterfaceToBool(requestBody["shouldOverwrite"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidShouldOverwrite",
			)
		}
	}

	copyUnixFileDto := dto.NewCopyUnixFile(sourcePath, destinationPath, shouldOverwrite)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err = useCase.CopyUnixFile(
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
// @Param        sourcePaths	body	[]string	true	"FilePaths to deleted."
// @Success      200 {object} object{} "FilesDeleted"
// @Router       /v1/files/delete/ [put]
func (controller *FilesController) Delete(c echo.Context) error {
	requiredParams := []string{"sourcePaths"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	_, isSourcePathsMap := requestBody["sourcePaths"].(map[string]interface{})
	if isSourcePathsMap {
		requestBody["sourcePaths"] = []interface{}{requestBody["sourcePaths"]}
	}

	sourcePathsSlice, assertOk := requestBody["sourcePaths"].([]interface{})
	if !assertOk {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "SourcePathMustBeArray")
	}

	sourcePaths, err := controller.parseSourcePaths(sourcePathsSlice)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	hardDelete := false
	if requestBody["hardDelete"] != nil {
		var err error
		hardDelete, err = voHelper.InterfaceToBool(requestBody["hardDelete"])
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

	err = deleteUnixFiles.Execute(deleteUnixFilesDto)
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
// @Param        compressFilesDto 	  body    dto.CompressUnixFiles  true  "All props are required."
// @Success      200 {object} object{} "FilesCompressed"
// @Success      207 {object} object{} "FilesArePartialCompressed"
// @Router       /v1/files/compress/ [post]
func (controller *FilesController) Compress(c echo.Context) error {
	requiredParams := []string{"sourcePaths", "destinationPath"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	_, isSourcePathsMap := requestBody["sourcePaths"].(map[string]interface{})
	if isSourcePathsMap {
		requestBody["sourcePaths"] = []interface{}{requestBody["sourcePaths"]}
	}

	sourcePathsSlice, assertOk := requestBody["sourcePaths"].([]interface{})
	if !assertOk {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, "SourcePathMustBeArray")
	}

	sourcePaths, err := controller.parseSourcePaths(sourcePathsSlice)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	destinationPath, err := valueObject.NewUnixFilePath(requestBody["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	var compressionUnixTypePtr *valueObject.UnixCompressionType
	if requestBody["compressionType"] != nil {
		compressionUnixType, err := valueObject.NewUnixCompressionType(
			requestBody["compressionType"].(string),
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		compressionUnixTypePtr = &compressionUnixType
	}

	compressUnixFilesDto := dto.NewCompressUnixFiles(
		sourcePaths,
		destinationPath,
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
// @Param        extractFilesDto 	  body    dto.ExtractUnixFiles  true  "All props are required."
// @Success      200 {object} object{} "FilesExtracted"
// @Router       /v1/files/extract/ [put]
func (controller *FilesController) Extract(c echo.Context) error {
	requiredParams := []string{"sourcePath", "destinationPath"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	sourcePath, err := valueObject.NewUnixFilePath(requestBody["sourcePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	destinationPath, err := valueObject.NewUnixFilePath(requestBody["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	extractUnixFilesDto := dto.NewExtractUnixFiles(sourcePath, destinationPath)

	filesQueryRepo := filesInfra.FilesQueryRepo{}
	filesCmdRepo := filesInfra.FilesCmdRepo{}

	err = useCase.ExtractUnixFiles(
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
// @Router       /v1/files/upload/ [post]
func (controller *FilesController) Upload(c echo.Context) error {
	requiredParams := []string{"destinationPath", "files"}
	requestBody, _ := apiHelper.ReadRequestBody(c)

	if requestBody["destinationPath"] == nil {
		requestBody["destinationPath"] = c.QueryParam("destinationPath")
	}

	apiHelper.CheckMissingParams(requestBody, requiredParams)

	destinationPath, err := valueObject.NewUnixFilePath(requestBody["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err)
	}

	var filesToUpload []valueObject.FileStreamHandler
	for _, requestBodyFile := range requestBody["files"].(map[string]*multipart.FileHeader) {
		fileStreamHandler, err := valueObject.NewFileStreamHandler(requestBodyFile)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
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
