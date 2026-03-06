package apiController

import (
	tkPresentation "github.com/goinfinite/tk/src/presentation"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
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
func (controller *FilesController) Read(echoContext echo.Context) error {
	requiredParams := []string{"sourcePath"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}
	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	sourcePath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["sourcePath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	shouldIncludeFileTree := false
	if requestData["shouldIncludeFileTree"] != nil {
		shouldIncludeFileTree, err = tkVoUtil.InterfaceToBool(
			requestData["shouldIncludeFileTree"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
	}

	readFilesRequestDto := dto.ReadFilesRequest{
		SourcePath:            sourcePath,
		ShouldIncludeFileTree: &shouldIncludeFileTree,
	}

	filesList, err := useCase.ReadFiles(controller.filesQueryRepo, readFilesRequestDto)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusOK, filesList)
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
func (controller *FilesController) Create(echoContext echo.Context) error {
	requiredParams := []string{"filePath"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}
	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	filePath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["filePath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	fileType := tkValueObject.MimeTypeGeneric
	if requestData["mimeType"] != nil {
		fileType, err = tkValueObject.NewMimeType(requestData["mimeType"])
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}

		if fileType != tkValueObject.MimeTypeDirectory {
			fileType = tkValueObject.MimeTypeGeneric
		}
	}

	successResponse := "FileCreated"

	filePermissions := valueObject.NewUnixFileDefaultPermissions()
	if fileType.IsDir() {
		filePermissions = valueObject.NewUnixDirDefaultPermissions()
		successResponse = "DirectoryCreated"
	}

	if requestData["permissions"] != nil {
		filePermissions, err = valueObject.NewUnixFilePermissions(
			requestData["permissions"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	createDto := dto.NewCreateUnixFile(
		filePath, &filePermissions, fileType, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateUnixFile(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, createDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusCreated, successResponse)
}

func (controller *FilesController) parseSourcePaths(
	rawSourcePathsUnknownType any,
) ([]tkValueObject.UnixAbsoluteFilePath, error) {
	sourcePaths := []tkValueObject.UnixAbsoluteFilePath{}

	rawSourcePathsStrSlice := []string{}
	switch rawSourcePathsValues := rawSourcePathsUnknownType.(type) {
	case string:
		rawSourcePathsStrSlice = []string{rawSourcePathsValues}
	case []string:
		rawSourcePathsStrSlice = rawSourcePathsValues
	case []interface{}:
		for _, rawSourcePath := range rawSourcePathsValues {
			rawSourcePathStr, err := tkVoUtil.InterfaceToString(rawSourcePath)
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
		sourcePath, err := tkValueObject.NewUnixAbsoluteFilePath(rawSourcePath, false)
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
func (controller *FilesController) Update(echoContext echo.Context) error {
	requiredParams := []string{"sourcePaths"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	if requestData["sourcePaths"] == nil {
		if _, exists := requestData["sourcePath"]; exists {
			requestData["sourcePaths"] = requestData["sourcePath"]
		}
	}

	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	sourcePaths, err := controller.parseSourcePaths(requestData["sourcePaths"])
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	var destinationPathPtr *tkValueObject.UnixAbsoluteFilePath
	if requestData["destinationPath"] != nil {
		destinationPath, err := tkValueObject.NewUnixAbsoluteFilePath(
			requestData["destinationPath"], false,
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		destinationPathPtr = &destinationPath
	}

	var permissionsPtr *valueObject.UnixFilePermissions
	if requestData["permissions"] != nil {
		permissions, err := valueObject.NewUnixFilePermissions(
			requestData["permissions"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		permissionsPtr = &permissions
	}

	var encodedContentPtr *valueObject.EncodedContent
	if requestData["encodedContent"] != nil {
		encodedContent, err := valueObject.NewEncodedContent(
			requestData["encodedContent"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		encodedContentPtr = &encodedContent
	}

	var ownershipPtr *tkValueObject.UnixFileOwnership
	if requestData["ownership"] != nil {
		ownership, err := tkValueObject.NewUnixFileOwnership(
			requestData["ownership"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		ownershipPtr = &ownership
	}

	var shouldFixPermissionsPtr *bool
	if requestData["shouldFixPermissions"] != nil {
		shouldFixPermissions, err := tkVoUtil.InterfaceToBool(
			requestData["shouldFixPermissions"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		shouldFixPermissionsPtr = &shouldFixPermissions
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
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
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	hasSuccess := len(updateProcessInfo.FilePathsSuccessfullyUpdated) > 0
	hasFailures := len(updateProcessInfo.FailedPathsWithReason) > 0
	if !hasSuccess && hasFailures {
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, updateProcessInfo)
	}
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		return apiHelper.ResponseWrapper(echoContext, http.StatusMultiStatus, updateProcessInfo)
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusOK, updateProcessInfo)
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
func (controller *FilesController) Copy(echoContext echo.Context) error {
	requiredParams := []string{"sourcePath", "destinationPath"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}
	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	sourcePath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["sourcePath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	destinationPath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["destinationPath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	shouldOverwrite := false
	if requestData["shouldOverwrite"] != nil {
		var err error
		shouldOverwrite, err = tkVoUtil.InterfaceToBool(requestData["shouldOverwrite"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				echoContext, http.StatusBadRequest, "InvalidShouldOverwrite",
			)
		}
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
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
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusCreated, "FileCopied")
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
func (controller *FilesController) Delete(echoContext echo.Context) error {
	requiredParams := []string{"sourcePaths"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	if requestData["sourcePaths"] == nil {
		if _, exists := requestData["sourcePath"]; exists {
			requestData["sourcePaths"] = requestData["sourcePath"]
		}
	}

	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	sourcePaths, err := controller.parseSourcePaths(requestData["sourcePaths"])
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	hardDelete := false
	if requestData["hardDelete"] != nil {
		var err error
		hardDelete, err = tkVoUtil.InterfaceToBool(requestData["hardDelete"])
		if err != nil {
			return apiHelper.ResponseWrapper(
				echoContext, http.StatusBadRequest, "InvalidHardDelete",
			)
		}
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
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
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusOK, "FilesDeleted")
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
func (controller *FilesController) Compress(echoContext echo.Context) error {
	requiredParams := []string{"sourcePaths", "destinationPath"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	if requestData["sourcePaths"] == nil {
		if _, exists := requestData["sourcePath"]; exists {
			requestData["sourcePaths"] = requestData["sourcePath"]
		}
	}

	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	sourcePaths, err := controller.parseSourcePaths(requestData["sourcePaths"])
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	destinationPath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["destinationPath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	var compressionUnixTypePtr *valueObject.UnixCompressionType
	if requestData["compressionType"] != nil {
		compressionUnixType, err := valueObject.NewUnixCompressionType(
			requestData["compressionType"],
		)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		compressionUnixTypePtr = &compressionUnixType
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
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
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	httpStatus := http.StatusCreated

	hasSuccess := len(compressionProcessInfo.FilePathsSuccessfullyCompressed) > 0
	hasFailures := len(compressionProcessInfo.FailedPathsWithReason) > 0
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		httpStatus = http.StatusMultiStatus
	}

	return apiHelper.ResponseWrapper(echoContext, httpStatus, compressionProcessInfo)
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
func (controller *FilesController) Extract(echoContext echo.Context) error {
	requiredParams := []string{"sourcePath", "destinationPath"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}
	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	sourcePath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["sourcePath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	destinationPath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["destinationPath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	extractDto := dto.NewExtractUnixFiles(
		sourcePath, destinationPath, operatorAccountId, operatorIpAddress,
	)

	err = useCase.ExtractUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, extractDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusCreated, "FilesExtracted")
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
func (controller *FilesController) Upload(echoContext echo.Context) error {
	requiredParams := []string{"destinationPath", "files"}
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	err := tkPresentation.RequiredParamsInspector(
		requestData, requiredParams,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(
			echoContext, http.StatusBadRequest, err.Error(),
		)
	}

	destinationPath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["destinationPath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	var filesToUpload []valueObject.FileStreamHandler
	for _, requestDataFile := range requestData["files"].(map[string]*multipart.FileHeader) {
		fileStreamHandler, err := valueObject.NewFileStreamHandler(requestDataFile)
		if err != nil {
			return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
		}
		filesToUpload = append(filesToUpload, fileStreamHandler)
	}

	operatorAccountId, err := tkValueObject.NewAccountId(
		requestData["operatorAccountId"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	operatorIpAddress, err := tkValueObject.NewIpAddress(
		requestData["operatorIpAddress"],
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	uploadDto := dto.NewUploadUnixFiles(
		destinationPath, filesToUpload, operatorAccountId, operatorIpAddress,
	)

	uploadProcessInfo, err := useCase.UploadUnixFiles(
		controller.filesQueryRepo, controller.filesCmdRepo,
		controller.activityRecordCmdRepo, uploadDto,
	)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, err.Error())
	}

	hasSuccess := len(uploadProcessInfo.FileNamesSuccessfullyUploaded) > 0
	hasFailures := len(uploadProcessInfo.FailedNamesWithReason) > 0
	if !hasSuccess && hasFailures {
		return apiHelper.ResponseWrapper(echoContext, http.StatusInternalServerError, uploadProcessInfo)
	}
	wasPartiallySuccessful := hasSuccess && hasFailures
	if wasPartiallySuccessful {
		return apiHelper.ResponseWrapper(echoContext, http.StatusMultiStatus, uploadProcessInfo)
	}

	return apiHelper.ResponseWrapper(echoContext, http.StatusOK, uploadProcessInfo)
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
func (controller *FilesController) Download(echoContext echo.Context) error {
	inputReader := tkPresentation.ApiRequestInputReader{}
	requestData, requestParsingErr := inputReader.Reader(echoContext)
	if requestParsingErr != nil {
		return requestParsingErr
	}

	sourcePath, err := tkValueObject.NewUnixAbsoluteFilePath(requestData["sourcePath"], false)
	if err != nil {
		return apiHelper.ResponseWrapper(echoContext, http.StatusBadRequest, err.Error())
	}

	return echoContext.Attachment(sourcePath.String(), sourcePath.ReadFileName(false).String())
}
