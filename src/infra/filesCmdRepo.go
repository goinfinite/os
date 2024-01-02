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

type FilesCmdRepo struct {
	uploadProcessReport dto.UploadProcessReport
}

func (repo FilesCmdRepo) addUploadFailure(
	errMessage string,
	fileStreamHandler valueObject.FileStreamHandler,
) {
	failureReason, _ := valueObject.NewFileProcessingFailure(errMessage)
	repo.uploadProcessReport.FailedNamesWithReason = append(
		repo.uploadProcessReport.FailedNamesWithReason,
		valueObject.NewUploadProcessFailure(fileStreamHandler.Name, failureReason),
	)
}

func (repo FilesCmdRepo) uploadSingleFile(
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

	repo.uploadProcessReport.FileNamesSuccessfullyUploaded = append(
		repo.uploadProcessReport.FileNamesSuccessfullyUploaded,
		fileToUpload.Name,
	)

	return nil
}

func (repo FilesCmdRepo) Copy(copyUnixFile dto.CopyUnixFile) error {
	fileToCopyExists := infraHelper.FileExists(copyUnixFile.SourcePath.String())
	if !fileToCopyExists {
		return errors.New("FileToCopyNotFound")
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
	existingFiles := []string{}
	for _, sourcePath := range compressUnixFiles.SourcePaths {
		sourcePathExists := infraHelper.FileExists(sourcePath.String())
		if !sourcePathExists {
			log.Printf("SourcePathNotFound: %s", sourcePath.String())
			continue
		}

		existingFiles = append(existingFiles, sourcePath.String())
	}

	if len(existingFiles) < 1 {
		return dto.CompressionProcessReport{}, errors.New("NoExistingFilesToCompress")
	}

	compressionTypeStr := "zip"

	destinationPathExt, err := compressUnixFiles.DestinationPath.GetFileExtension()
	if err == nil {
		compressionTypeStr = destinationPathExt.String()

		if compressUnixFiles.CompressionType != nil {
			compressionTypeStr = compressUnixFiles.CompressionType.String()
		}
	}

	destinationPathWithoutExt := compressUnixFiles.DestinationPath.GetWithoutExtension()
	newDestinationPath, err := valueObject.NewUnixFilePath(
		destinationPathWithoutExt.String() + "." + compressionTypeStr,
	)
	if err != nil {
		return dto.CompressionProcessReport{}, errors.New(
			"CannotUpdateDestinationPathWithNewExtension",
		)
	}
	compressUnixFiles.DestinationPath = newDestinationPath

	_, err = valueObject.NewUnixCompressionType(compressionTypeStr)
	if err != nil {
		return dto.CompressionProcessReport{}, errors.New("UnsupportedCompressionType")
	}

	destinationPathExists := infraHelper.FileExists(newDestinationPath.String())
	if destinationPathExists {
		return dto.CompressionProcessReport{}, errors.New("DestinationPathAlreadyExists")
	}

	compressionBinary := "zip"
	compressionBinaryFlag := "-qr"
	if compressionTypeStr != "zip" {
		compressionBinary = "tar"
		compressionBinaryFlag = "-czf"
	}

	filesToCompress := strings.Join(existingFiles, " ")
	_, err = infraHelper.RunCmd(
		compressionBinary,
		compressionBinaryFlag,
		newDestinationPath.String(),
		filesToCompress,
	)
	if err != nil {
		return dto.CompressionProcessReport{}, err
	}

	compressionProcessReport := dto.NewCompressionProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.CompressionProcessFailure{},
		newDestinationPath,
	)
	for _, sourcePath := range compressUnixFiles.SourcePaths {
		if !slices.Contains(existingFiles, sourcePath.String()) {
			compressionProcessReport.FailedPathsWithReason = append(
				compressionProcessReport.FailedPathsWithReason,
				valueObject.NewCompressionProcessFailure(
					sourcePath,
					"SourcePathNotFound",
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
			log.Printf("DeleteFileError (%s): FileNotFound", fileToDelete.String())
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
		return errors.New("FileNotFound")
	}

	destinationPath := extractUnixFiles.DestinationPath

	destinationPathExists := infraHelper.FileExists(destinationPath.String())
	if destinationPathExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	compressBinary := "tar"
	compressBinaryFlag := "-xf"
	compressDestinationFlag := "-C"

	unixFilePathExtension, err := fileToExtract.GetFileExtension()
	if err != nil {
		return err
	}

	if unixFilePathExtension.String() == "zip" {
		compressBinary = "unzip"
		compressBinaryFlag = "-qq"
		compressDestinationFlag = "-d"
	}

	err = infraHelper.MakeDir(destinationPath.String())
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
		return errors.New("FileToMoveNotFound")
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

	fileToUpdate, err := queryRepo.GetOne(updateUnixFile.SourcePath)
	if err != nil {
		return err
	}

	if fileToUpdate.MimeType.IsDir() {
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

	repo.uploadProcessReport = dto.NewUploadProcessReport(
		[]valueObject.UnixFileName{},
		[]valueObject.UploadProcessFailure{},
		destinationPath,
	)

	destinationFile, err := queryRepo.GetOne(destinationPath)
	if err != nil {
		return repo.uploadProcessReport, err
	}

	if !destinationFile.MimeType.IsDir() {
		return repo.uploadProcessReport, errors.New("DestinationPathCannotBeAFile")
	}

	for _, fileToUpload := range uploadUnixFiles.FileStreamHandlers {
		err := repo.uploadSingleFile(
			destinationPath,
			fileToUpload,
		)
		if err != nil {
			for _, fileStreamHandler := range uploadUnixFiles.FileStreamHandlers {
				repo.addUploadFailure(err.Error(), fileStreamHandler)
			}
		}
	}

	return repo.uploadProcessReport, nil
}
