package apiController

import (
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	filesInfra "github.com/goinfinite/os/src/infra/files"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	apiHelper "github.com/goinfinite/os/src/presentation/api/helper"
	"github.com/labstack/echo/v4"
)

type FilesController struct {
	filesQueryRepo        *filesInfra.FilesQueryRepo
	filesCmdRepo          *filesInfra.FilesCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewFilesController(
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *FilesController {
	return &FilesController{
		filesQueryRepo:        &filesInfra.FilesQueryRepo{},
		filesCmdRepo:          filesInfra.NewFilesCmdRepo(),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

// ReadFiles    godoc
// @Summary      ReadFiles
// @Description  List dir/files.
// @Tags         files
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        sourcePath	query	string	true	"SourcePath"
// @Param        shouldIncludeFileTree query  bool  false  "ShouldIncludeFileTree"
// @Success      200 {array} dto.ReadFilesResponse
// @Router       /v1/files/ [get]
func (controller *FilesController) Read(c echo.Context) error {
	requiredParams := []string{"sourcePath"}
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}
	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	sourcePath, err := valueObject.NewUnixFilePath(requestInputData["sourcePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	shouldIncludeFileTree := false
	if requestInputData["shouldIncludeFileTree"] != nil {
		shouldIncludeFileTree, err = voHelper.InterfaceToBool(
			requestInputData["shouldIncludeFileTree"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}

	readFilesRequestDto := dto.ReadFilesRequest{
		SourcePath:            sourcePath,
		ShouldIncludeFileTree: &shouldIncludeFileTree,
	}

	filesList, err := useCase.ReadFiles(controller.filesQueryRepo, readFilesRequestDto)
	if err != nil {
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}
	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	filePath, err := valueObject.NewUnixFilePath(requestInputData["filePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	fileType := valueObject.MimeTypeGeneric
	if requestInputData["mimeType"] != nil {
		fileType, err = valueObject.NewMimeType(requestInputData["mimeType"])
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}

		if fileType != valueObject.MimeTypeDirectory {
			fileType = valueObject.MimeTypeGeneric
		}
	}

	successResponse := "FileCreated"

	filePermissions := valueObject.NewUnixFileDefaultPermissions()
	if fileType.IsDir() {
		filePermissions = valueObject.NewUnixDirDefaultPermissions()
		successResponse = "DirectoryCreated"
	}

	if requestInputData["permissions"] != nil {
		filePermissions, err = valueObject.NewUnixFilePermissions(
			requestInputData["permissions"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	createDto := dto.NewCreateUnixFile(
		filePath, &filePermissions, fileType, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateUnixFile(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(c, http.StatusCreated, successResponse)
}

func (controller *FilesController) parseSourcePaths(
	rawSourcePathsUnknownType any,
) ([]valueObject.UnixFilePath, error) {
	sourcePaths := []valueObject.UnixFilePath{}

	rawSourcePathsStrSlice := []string{}
	switch rawSourcePathsValues := rawSourcePathsUnknownType.(type) {
	case string:
		rawSourcePathsStrSlice = []string{rawSourcePathsValues}
	case []string:
		rawSourcePathsStrSlice = rawSourcePathsValues
	case []interface{}:
		for _, rawSourcePath := range rawSourcePathsValues {
			rawSourcePathStr, err := voHelper.InterfaceToString(rawSourcePath)
			if err != nil {
				slog.Debug(err.Error(), slog.Any("rawSourcePath", rawSourcePath))
				continue
			}
			rawSourcePathsStrSlice = append(rawSourcePathsStrSlice, rawSourcePathStr)
		}
	default:
		return sourcePaths, errors.New("SourcePathsMustBeStringSlice")
	}

	for _, rawSourcePath := range rawSourcePathsStrSlice {
		sourcePath, err := valueObject.NewUnixFilePath(rawSourcePath)
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawSourcePath", rawSourcePath))
			continue
		}

		sourcePaths = append(sourcePaths, sourcePath)
	}

	return sourcePaths, nil
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["sourcePaths"] == nil {
		if _, exists := requestInputData["sourcePath"]; exists {
			requestInputData["sourcePaths"] = requestInputData["sourcePath"]
		}
	}

	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	sourcePaths, err := controller.parseSourcePaths(requestInputData["sourcePaths"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	var destinationPathPtr *valueObject.UnixFilePath
	if requestInputData["destinationPath"] != nil {
		destinationPath, err := valueObject.NewUnixFilePath(
			requestInputData["destinationPath"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		destinationPathPtr = &destinationPath
	}

	var permissionsPtr *valueObject.UnixFilePermissions
	if requestInputData["permissions"] != nil {
		permissions, err := valueObject.NewUnixFilePermissions(
			requestInputData["permissions"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		permissionsPtr = &permissions
	}

	var encodedContentPtr *valueObject.EncodedContent
	if requestInputData["encodedContent"] != nil {
		encodedContent, err := valueObject.NewEncodedContent(
			requestInputData["encodedContent"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		encodedContentPtr = &encodedContent
	}

	var ownershipPtr *valueObject.UnixFileOwnership
	if requestInputData["ownership"] != nil {
		ownership, err := valueObject.NewUnixFileOwnership(
			requestInputData["ownership"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		ownershipPtr = &ownership
	}

	var shouldFixPermissionsPtr *bool
	if requestInputData["shouldFixPermissions"] != nil {
		shouldFixPermissions, err := voHelper.InterfaceToBool(
			requestInputData["shouldFixPermissions"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		shouldFixPermissionsPtr = &shouldFixPermissions
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	updateUnixFileDto := dto.NewUpdateUnixFiles(
		sourcePaths, destinationPathPtr, permissionsPtr, encodedContentPtr,
		ownershipPtr, shouldFixPermissionsPtr, operatorAccountId, operatorIpAddress,
	)

	updateUnixFileUc := useCase.NewUpdateUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo, controller.activityRecordCmdRepo,
	)
	updateProcessInfo, err := updateUnixFileUc.Execute(updateUnixFileDto)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	hasSuccess := len(updateProcessInfo.FilePathsSuccessfullyUpdated) > 0
	hasFailures := len(updateProcessInfo.FailedPathsWithReason) > 0
	if !hasSuccess && hasFailures {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, updateProcessInfo)
	}
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		return apiHelper.ResponseWrapper(c, http.StatusMultiStatus, updateProcessInfo)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, updateProcessInfo)
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}
	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	sourcePath, err := valueObject.NewUnixFilePath(requestInputData["sourcePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	destinationPath, err := valueObject.NewUnixFilePath(requestInputData["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	shouldOverwrite := false
	if requestInputData["shouldOverwrite"] != nil {
		var err error
		shouldOverwrite, err = voHelper.InterfaceToBool(requestInputData["shouldOverwrite"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidShouldOverwrite",
			)
		}
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	copyDto := dto.NewCopyUnixFile(
		sourcePath, destinationPath, shouldOverwrite, operatorAccountId,
		operatorIpAddress,
	)

	err = useCase.CopyUnixFile(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, copyDto,
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["sourcePaths"] == nil {
		if _, exists := requestInputData["sourcePath"]; exists {
			requestInputData["sourcePaths"] = requestInputData["sourcePath"]
		}
	}

	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	sourcePaths, err := controller.parseSourcePaths(requestInputData["sourcePaths"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	hardDelete := false
	if requestInputData["hardDelete"] != nil {
		var err error
		hardDelete, err = voHelper.InterfaceToBool(requestInputData["hardDelete"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				c, http.StatusBadRequest, "InvalidHardDelete",
			)
		}
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	deleteDto := dto.NewDeleteUnixFiles(
		sourcePaths, hardDelete, operatorAccountId, operatorIpAddress,
	)

	deleteUnixFiles := useCase.NewDeleteUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo,
	)

	err = deleteUnixFiles.Execute(deleteDto)
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	if requestInputData["sourcePaths"] == nil {
		if _, exists := requestInputData["sourcePath"]; exists {
			requestInputData["sourcePaths"] = requestInputData["sourcePath"]
		}
	}

	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	sourcePaths, err := controller.parseSourcePaths(requestInputData["sourcePaths"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	destinationPath, err := valueObject.NewUnixFilePath(requestInputData["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	var compressionUnixTypePtr *valueObject.UnixCompressionType
	if requestInputData["compressionType"] != nil {
		compressionUnixType, err := valueObject.NewUnixCompressionType(
			requestInputData["compressionType"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		compressionUnixTypePtr = &compressionUnixType
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	compressDto := dto.NewCompressUnixFiles(
		sourcePaths, destinationPath, compressionUnixTypePtr, operatorAccountId,
		operatorIpAddress,
	)

	compressionProcessInfo, err := useCase.CompressUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, compressDto,
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}
	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	sourcePath, err := valueObject.NewUnixFilePath(requestInputData["sourcePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	destinationPath, err := valueObject.NewUnixFilePath(requestInputData["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	extractDto := dto.NewExtractUnixFiles(
		sourcePath, destinationPath, operatorAccountId, operatorIpAddress,
	)

	err = useCase.ExtractUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, extractDto,
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
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	apiHelper.CheckMissingParams(requestInputData, requiredParams)

	destinationPath, err := valueObject.NewUnixFilePath(requestInputData["destinationPath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	var filesToUpload []valueObject.FileStreamHandler
	for _, requestInputDataFile := range requestInputData["files"].(map[string]*multipart.FileHeader) {
		fileStreamHandler, err := valueObject.NewFileStreamHandler(requestInputDataFile)
		if err != nil {
			return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
		}
		filesToUpload = append(filesToUpload, fileStreamHandler)
	}

	operatorAccountId, err := valueObject.NewAccountId(
		requestInputData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := valueObject.NewIpAddress(
		requestInputData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	uploadDto := dto.NewUploadUnixFiles(
		destinationPath, filesToUpload, operatorAccountId, operatorIpAddress,
	)

	uploadProcessInfo, err := useCase.UploadUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, uploadDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, err.Error())
	}

	hasSuccess := len(uploadProcessInfo.FileNamesSuccessfullyUploaded) > 0
	hasFailures := len(uploadProcessInfo.FailedNamesWithReason) > 0
	if !hasSuccess && hasFailures {
		return apiHelper.ResponseWrapper(c, http.StatusInternalServerError, uploadProcessInfo)
	}
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		return apiHelper.ResponseWrapper(c, http.StatusMultiStatus, uploadProcessInfo)
	}

	return apiHelper.ResponseWrapper(c, http.StatusOK, uploadProcessInfo)
}

// DownloadFile    godoc
// @Summary      DownloadFile
// @Description  Download a file.
// @Tags         files
// @Accept       json
// @Produce      octet-stream
// @Security     Bearer
// @Param        sourcePath	query	string	true	"SourcePath"
// @Success      200 {file} file
// @Router       /v1/files/download/ [get]
func (controller *FilesController) Download(c echo.Context) error {
	requestInputData, err := apiHelper.ReadRequestInputData(c)
	if err != nil {
		return err
	}

	sourcePath, err := valueObject.NewUnixFilePath(requestInputData["sourcePath"])
	if err != nil {
		return apiHelper.ResponseWrapper(c, http.StatusBadRequest, err.Error())
	}

	return c.Attachment(sourcePath.String(), sourcePath.ReadFileName().String())
}
