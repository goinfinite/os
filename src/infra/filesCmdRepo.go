package infra

import (
	"errors"
	"io"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type FilesCmdRepo struct{}

var uploadProcessReport dto.UploadProcessReport

func uploadProcessReportFailureListFactory(
	errMessage string,
	fileStreamHandlers []valueObject.FileStreamHandler,
) []valueObject.UploadProcessFailure {
	uploadProcessReportFailureList := []valueObject.UploadProcessFailure{}

	for _, fileStreamHandler := range fileStreamHandlers {
		failureReason, _ := valueObject.NewFileProcessingFailure(errMessage)
		uploadProcessReportFailureList = append(
			uploadProcessReportFailureList,
			valueObject.NewUploadProcessFailure(fileStreamHandler.Name, failureReason),
		)
	}

	return uploadProcessReportFailureList
}

func uploadSingleFile(
	destinationPath valueObject.UnixFilePath,
	fileToUpload valueObject.FileStreamHandler,
) error {
	destinationFilePath := destinationPath.String() + "/" + fileToUpload.Name.String()
	destinationEmptyFile, err := os.Create(destinationFilePath)
	if err != nil {
		return errors.New("CreateEmptyFileToStoreUploadFileError: " + err.Error())
	}
	defer destinationEmptyFile.Close()

	fileToUploadStream, err := fileToUpload.Open()
	if err != nil {
		return errors.New("UnableToOpenFileStream: " + err.Error())
	}

	_, err = io.Copy(destinationEmptyFile, fileToUploadStream)
	if err != nil {
		return errors.New("CopyFileStreamHandlerContentToDestinationFileError: " + err.Error())
	}

	uploadProcessReport.FileNamesSuccessfullyUploaded = append(
		uploadProcessReport.FileNamesSuccessfullyUploaded,
		fileToUpload.Name,
	)

	return nil
}

func (repo FilesCmdRepo) Copy(copyUnixFile dto.CopyUnixFile) error {
	fileToCopyExists := infraHelper.FileExists(copyUnixFile.SourcePath.String())
	if !fileToCopyExists {
		return errors.New("FileToCopyDoesNotExists")
	}

	destinationPathExists := infraHelper.FileExists(copyUnixFile.DestinationPath.String())
	if destinationPathExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	_, err := infraHelper.RunCmd(
		"rsync",
		"-avq",
		copyUnixFile.SourcePath.String(),
		copyUnixFile.DestinationPath.String(),
	)
	return err
}

func (repo FilesCmdRepo) Compress(
	compressUnixFiles dto.CompressUnixFiles,
) (dto.CompressionProcessReport, error) {
	sourcePathsThatExists := []string{}
	for _, sourcePath := range compressUnixFiles.SourcePaths {
		sourcePathExists := infraHelper.FileExists(sourcePath.String())
		if !sourcePathExists {
			log.Printf("SourcePathDoesNotExists: %s", sourcePath.String())
			continue
		}

		sourcePathsThatExists = append(sourcePathsThatExists, sourcePath.String())
	}

	if len(sourcePathsThatExists) < 1 {
		return dto.CompressionProcessReport{}, errors.New("NoExistingFilesToCompress")
	}

	compressionTypePtr := compressUnixFiles.CompressionType
	if compressionTypePtr == nil {
		destinationPathExt := compressUnixFiles.DestinationPath.GetFileExtension()
		destinationPathExtStr := destinationPathExt.String()
		if destinationPathExt.String() == "" {
			destinationPathExtStr = "zip"
		}

		compressionType, err := valueObject.NewUnixCompressionType(destinationPathExtStr)
		if err != nil {
			return dto.CompressionProcessReport{}, err
		}

		compressionTypePtr = &compressionType
	}

	compressBinary := "zip"
	compressBinaryFlag := "-qr"
	compressExtFile := ".zip"
	if compressionTypePtr.String() != "zip" {
		compressBinary = "tar"
		compressBinaryFlag = "-czf"
		compressExtFile = ".tar.gz"
	}

	destinationPathWithoutExt := compressUnixFiles.DestinationPath.GetWithoutExtension()
	destinationPathWithCompressionTypeAsExtStr := destinationPathWithoutExt.String() + compressExtFile
	destinationPathWithCompressionTypeAsExt, err := valueObject.NewUnixFilePath(
		destinationPathWithCompressionTypeAsExtStr,
	)
	if err != nil {
		return dto.CompressionProcessReport{}, err
	}

	destinationPathExists := infraHelper.FileExists(destinationPathWithCompressionTypeAsExt.String())
	if destinationPathExists {
		return dto.CompressionProcessReport{}, errors.New("DestinationPathAlreadyExists")
	}

	filesToCompress := strings.Join(sourcePathsThatExists, " ")
	_, err = infraHelper.RunCmd(
		compressBinary,
		compressBinaryFlag,
		destinationPathWithCompressionTypeAsExt.String(),
		filesToCompress,
	)
	if err != nil {
		return dto.CompressionProcessReport{}, err
	}

	compressionProcessReport := dto.NewCompressionProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.CompressionProcessFailure{},
		destinationPathWithCompressionTypeAsExt,
	)
	for _, sourcePath := range compressUnixFiles.SourcePaths {
		if !slices.Contains(sourcePathsThatExists, sourcePath.String()) {
			compressionProcessReport.FailedPathsWithReason = append(
				compressionProcessReport.FailedPathsWithReason,
				valueObject.NewCompressionProcessFailure(
					sourcePath,
					"SourcePathDoesNotExists",
				),
			)
		}

		compressionProcessReport.FilePathsSuccessfullyCompressed = append(
			compressionProcessReport.FilePathsSuccessfullyCompressed,
			sourcePath,
		)
	}

	return compressionProcessReport, nil
}

func (repo FilesCmdRepo) Create(createUnixFile dto.CreateUnixFile) error {
	filesExists := infraHelper.FileExists(createUnixFile.SourcePath.String())
	if filesExists {
		return errors.New("PathAlreadyExists")
	}

	if !createUnixFile.MimeType.IsDir() {
		_, err := os.Create(createUnixFile.SourcePath.String())
		if err != nil {
			return err
		}

		return repo.UpdatePermissions(
			createUnixFile.SourcePath,
			createUnixFile.Permissions,
		)
	}

	err := os.MkdirAll(createUnixFile.SourcePath.String(), createUnixFile.Permissions.GetFileMode())
	if err != nil {
		return err
	}

	return nil
}

func (repo FilesCmdRepo) Delete(
	unixFilePathList []valueObject.UnixFilePath,
) {
	for _, fileToDelete := range unixFilePathList {
		fileExists := infraHelper.FileExists(fileToDelete.String())
		if !fileExists {
			log.Printf("DeleteFileError (%s): FileDoesNotExists", fileToDelete.String())
			continue
		}

		err := os.RemoveAll(fileToDelete.String())
		if err != nil {
			log.Printf("DeleteFileError (%s): %s", fileToDelete.String(), err)
			continue
		}

		log.Printf("File '%s' deleted.", fileToDelete.String())
	}
}

func (repo FilesCmdRepo) Extract(extractUnixFiles dto.ExtractUnixFiles) error {
	fileToExtract := extractUnixFiles.SourcePath

	fileToExtractExists := infraHelper.FileExists(fileToExtract.String())
	if !fileToExtractExists {
		return errors.New("FileDoesNotExists")
	}

	destinationPath := extractUnixFiles.DestinationPath

	destinationPathExists := infraHelper.FileExists(destinationPath.String())
	if destinationPathExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	compressBinary := "tar"
	compressBinaryFlag := "-xf"
	compressDestinationFlag := "-C"

	unixFilePathExtension := fileToExtract.GetFileExtension()
	if unixFilePathExtension.String() == "zip" {
		compressBinary = "unzip"
		compressBinaryFlag = "-qq"
		compressDestinationFlag = "-d"
	}

	err := infraHelper.MakeDir(destinationPath.String())
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd(
		compressBinary,
		compressBinaryFlag,
		fileToExtract.String(),
		compressDestinationFlag,
		destinationPath.String(),
	)
	return err
}

func (repo FilesCmdRepo) Move(updateUnixFile dto.UpdateUnixFile) error {
	fileToMoveExists := infraHelper.FileExists(updateUnixFile.SourcePath.String())
	if !fileToMoveExists {
		return errors.New("FileToMoveDoesNotExists")
	}

	destinationPathExists := infraHelper.FileExists(updateUnixFile.DestinationPath.String())
	if destinationPathExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	return os.Rename(
		updateUnixFile.SourcePath.String(),
		updateUnixFile.DestinationPath.String(),
	)
}

func (repo FilesCmdRepo) UpdateContent(
	updateUnixFile dto.UpdateUnixFile,
) error {
	queryRepo := FilesQueryRepo{}

	fileToUpdateContent, err := queryRepo.GetOne(updateUnixFile.SourcePath)
	if err != nil {
		return err
	}

	if fileToUpdateContent.MimeType.IsDir() {
		return errors.New("PathIsADir")
	}

	decodedContent, err := updateUnixFile.EncodedContent.GetDecodedContent()
	if err != nil {
		return err
	}

	return infraHelper.UpdateFile(
		updateUnixFile.SourcePath.String(),
		decodedContent,
		true,
	)
}

func (repo FilesCmdRepo) UpdatePermissions(
	unixFilePath valueObject.UnixFilePath,
	unixFilePermissions valueObject.UnixFilePermissions,
) error {
	queryRepo := FilesQueryRepo{}

	_, err := queryRepo.Get(unixFilePath)
	if err != nil {
		return err
	}

	return os.Chmod(unixFilePath.String(), unixFilePermissions.GetFileMode())
}

func (repo FilesCmdRepo) Upload(
	uploadUnixFiles dto.UploadUnixFiles,
) (dto.UploadProcessReport, error) {
	queryRepo := FilesQueryRepo{}

	destinationPath := uploadUnixFiles.DestinationPath

	uploadProcessReport = dto.NewUploadProcessReport(
		[]valueObject.UnixFileName{},
		[]valueObject.UploadProcessFailure{},
		destinationPath,
	)

	destinationFile, err := queryRepo.GetOne(destinationPath)
	if err != nil {
		return uploadProcessReport, err
	}

	if !destinationFile.MimeType.IsDir() {
		return uploadProcessReport, errors.New("DestinationPathCannotBeAFile")
	}

	for _, fileToUpload := range uploadUnixFiles.FileStreamHandlers {
		err := uploadSingleFile(
			destinationPath,
			fileToUpload,
		)
		if err != nil {
			uploadProcessReport.FailedNamesWithReason = append(
				uploadProcessReport.FailedNamesWithReason,
				uploadProcessReportFailureListFactory(
					err.Error(),
					uploadUnixFiles.FileStreamHandlers,
				)...,
			)
		}
	}

	return uploadProcessReport, nil
}
