package filesInfra

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

func (repo FilesCmdRepo) uploadFailureFactory(
	errMessage string,
	fileStreamHandler valueObject.FileStreamHandler,
) (valueObject.UploadProcessFailure, error) {
	failureReason, err := valueObject.NewFailureReason(errMessage)
	if err != nil {
		return valueObject.UploadProcessFailure{}, err
	}

	return valueObject.NewUploadProcessFailure(
		fileStreamHandler.Name,
		failureReason,
	), nil
}

func (repo FilesCmdRepo) uploadSingleFile(
	destinationPath valueObject.UnixFilePath,
	fileToUpload valueObject.FileStreamHandler,
) error {
	destinationFilePath := destinationPath.String() + "/" + fileToUpload.Name.String()
	destinationEmptyFile, err := os.Create(destinationFilePath)
	if err != nil {
		return errors.New("CreateEmptyFileError: " + err.Error())
	}
	defer destinationEmptyFile.Close()

	fileToUploadStream, err := fileToUpload.Open()
	if err != nil {
		return errors.New("UnableToOpenFileStream: " + err.Error())
	}

	_, err = io.Copy(destinationEmptyFile, fileToUploadStream)
	if err != nil {
		return errors.New("CopyFileContentToDestinationError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) Copy(copyUnixFile dto.CopyUnixFile) error {
	fileToCopyExists := infraHelper.FileExists(copyUnixFile.SourcePath.String())
	if !fileToCopyExists {
		return errors.New("FileToCopyNotFound")
	}

	if !copyUnixFile.ShouldOverwrite {
		destinationPathExists := infraHelper.FileExists(copyUnixFile.DestinationPath.String())
		if destinationPathExists {
			return errors.New("DestinationPathAlreadyExists")
		}
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
		destinationPathExtStr := destinationPathExt.String()
		if destinationPathExtStr != "zip" {
			compressionTypeStr = "tgz"
		}
	}

	if compressUnixFiles.CompressionType != nil {
		compressionTypeStr = compressUnixFiles.CompressionType.String()
	}

	destinationPathWithoutExt := compressUnixFiles.DestinationPath.GetWithoutExtension()
	compressionTypeAsExt := compressionTypeStr
	newDestinationPath, err := valueObject.NewUnixFilePath(
		destinationPathWithoutExt.String() + "." + compressionTypeAsExt,
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
	_, err = infraHelper.RunCmdWithSubShell(
		compressionBinary + " " +
			compressionBinaryFlag + " " +
			newDestinationPath.String() + " " +
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

			continue
		}

		compressionProcessReport.FilePathsSuccessfullyCompressed = append(
			compressionProcessReport.FilePathsSuccessfullyCompressed,
			sourcePath,
		)
	}

	return compressionProcessReport, nil
}

func (repo FilesCmdRepo) Create(createUnixFile dto.CreateUnixFile) error {
	filesExists := infraHelper.FileExists(createUnixFile.FilePath.String())
	if filesExists {
		return errors.New("PathAlreadyExists")
	}

	if !createUnixFile.MimeType.IsDir() {
		_, err := os.Create(createUnixFile.FilePath.String())
		if err != nil {
			return err
		}

		return repo.UpdatePermissions(
			createUnixFile.FilePath,
			createUnixFile.Permissions,
		)
	}

	err := os.MkdirAll(createUnixFile.FilePath.String(), createUnixFile.Permissions.GetFileMode())
	if err != nil {
		return err
	}

	return nil
}

func (repo FilesCmdRepo) Delete(unixFilePath valueObject.UnixFilePath) error {
	fileExists := infraHelper.FileExists(unixFilePath.String())
	if !fileExists {
		return errors.New("DeleteFileError: FileNotFound")
	}

	err := os.RemoveAll(unixFilePath.String())
	if err != nil {
		return errors.New("DeleteFileError: " + err.Error())
	}

	return nil
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

func (repo FilesCmdRepo) Move(
	unixSrcFilePath valueObject.UnixFilePath,
	unixDestinationDir valueObject.UnixFilePath,
	shouldOverwrite bool,
) error {
	fileToMoveExists := infraHelper.FileExists(unixSrcFilePath.String())
	if !fileToMoveExists {
		return errors.New("FileToMoveNotFound")
	}

	fileDestinationAbsolutePathStr := unixDestinationDir.String() + "/" + unixSrcFilePath.GetFileName().String()
	fileDestinationAbsolutePath, err := valueObject.NewUnixFilePath(fileDestinationAbsolutePathStr)
	if err != nil {
		return errors.New(err.Error() + ": " + fileDestinationAbsolutePathStr)
	}

	if infraHelper.FileExists(fileDestinationAbsolutePathStr) {
		if !shouldOverwrite {
			return errors.New("DestinationPathAlreadyExists")
		}

		err := repo.Delete(fileDestinationAbsolutePath)
		if err != nil {
			return errors.New("FailedToReplaceTrashFile: " + err.Error())
		}
	}

	return os.Rename(
		unixSrcFilePath.String(),
		fileDestinationAbsolutePathStr,
	)
}

func (repo FilesCmdRepo) UpdateContent(
	unixSrcFilePath valueObject.UnixFilePath,
	unixFileEncodedContent valueObject.EncodedContent,
) error {
	queryRepo := FilesQueryRepo{}

	fileToUpdate, err := queryRepo.GetOne(unixSrcFilePath)
	if err != nil {
		return err
	}

	if fileToUpdate.MimeType.IsDir() {
		return errors.New("PathIsADir")
	}

	decodedContent, err := unixFileEncodedContent.GetDecodedContent()
	if err != nil {
		return err
	}

	return infraHelper.UpdateFile(
		unixSrcFilePath.String(),
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

	uploadProcessReport := dto.NewUploadProcessReport(
		[]valueObject.UnixFileName{},
		[]valueObject.UploadProcessFailure{},
		destinationPath,
	)

	destinationFile, err := queryRepo.GetOne(destinationPath)
	if err != nil {
		return uploadProcessReport, errors.New("DestinationFileNotFound")
	}

	if !destinationFile.MimeType.IsDir() {
		return uploadProcessReport, errors.New("DestinationPathCannotBeAFile")
	}

	for _, fileToUpload := range uploadUnixFiles.FileStreamHandlers {
		err := repo.uploadSingleFile(
			destinationPath,
			fileToUpload,
		)

		if err != nil {
			uploadFailure, err := repo.uploadFailureFactory(err.Error(), fileToUpload)
			if err != nil {
				log.Printf("AddUploadFailureError: %s", err.Error())
			}

			uploadProcessReport.FailedNamesWithReason = append(
				uploadProcessReport.FailedNamesWithReason,
				uploadFailure,
			)

			continue
		}

		uploadProcessReport.FileNamesSuccessfullyUploaded = append(
			uploadProcessReport.FileNamesSuccessfullyUploaded,
			fileToUpload.Name,
		)

	}

	return uploadProcessReport, nil
}
