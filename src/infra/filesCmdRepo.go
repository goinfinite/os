package infra

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type FilesCmdRepo struct{}

func (repo FilesCmdRepo) Create(addUnixFile dto.AddUnixFile) error {
	filesExists := infraHelper.FileExists(addUnixFile.Path.String())
	if filesExists {
		return errors.New("PathAlreadyExists")
	}

	if !addUnixFile.Type.IsDir() {
		_, err := os.Create(addUnixFile.Path.String())
		if err != nil {
			log.Printf("CreateUnixFileError: %s", err)
			return errors.New("CreateUnixFileError")
		}

		return repo.UpdatePermissions(
			addUnixFile.Path,
			addUnixFile.Permissions,
		)
	}

	err := os.MkdirAll(addUnixFile.Path.String(), addUnixFile.Permissions.GetFileMode())
	if err != nil {
		log.Printf("CreateUnixDirectoryError: %s", err)
		return errors.New("CreateUnixDirectoryError")
	}

	return nil
}

func (repo FilesCmdRepo) Move(updateUnixFile dto.UpdateUnixFile) error {
	fileToMoveExists := infraHelper.FileExists(updateUnixFile.Path.String())
	if !fileToMoveExists {
		return errors.New("FileToMoveDoesNotExists")
	}

	destinationPathExists := infraHelper.FileExists(updateUnixFile.DestinationPath.String())
	if destinationPathExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	return os.Rename(
		updateUnixFile.Path.String(),
		updateUnixFile.DestinationPath.String(),
	)
}

func (repo FilesCmdRepo) Copy(copyUnixFile dto.CopyUnixFile) error {
	fileToCopyExists := infraHelper.FileExists(copyUnixFile.OriginPath.String())
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
		copyUnixFile.OriginPath.String(),
		copyUnixFile.DestinationPath.String(),
	)
	return err
}

func (repo FilesCmdRepo) UpdateContent(
	updateUnixFileContent dto.UpdateUnixFileContent,
) error {
	return infraHelper.UpdateFile(
		updateUnixFileContent.Path.String(),
		updateUnixFileContent.Content.GetDecodedContent(),
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

func (repo FilesCmdRepo) Compress(
	compressUnixFiles dto.CompressUnixFiles,
) (dto.CompressionProcessReport, error) {
	queryRepo := FilesQueryRepo{}

	_, err := queryRepo.GetOnly(compressUnixFiles.DestinationPath)
	if err != nil {
		return dto.CompressionProcessReport{}, err
	}

	compressBinary := "tar"
	compressBinaryFlag := "-czf"
	if compressUnixFiles.CompressionType.String() == "zip" {
		compressBinary = "zip"
		compressBinaryFlag = "-qr"
	}

	failedToCompressList := []valueObject.CompressionProcessFailure{}

	filesToCompressStr := compressUnixFiles.Paths[0].String()
	if len(compressUnixFiles.Paths) > 1 {
		var filesToCompressStrSlice []string
		for _, filePath := range compressUnixFiles.Paths {
			_, err := queryRepo.GetOnly(filePath)
			if err != nil {
				compressionProcessFailure := valueObject.NewCompressionProcessFailure(
					filePath,
					err.Error(),
				)
				failedToCompressList = append(failedToCompressList, compressionProcessFailure)

				continue
			}

			filesToCompressStrSlice = append(filesToCompressStrSlice, filePath.String())
		}

		filesToCompressStr = strings.Join(filesToCompressStrSlice, " ")
	}

	compressedFilePathWithoutExt := strings.Split(compressUnixFiles.DestinationPath.String(), ".")[0]
	compressedFilePathWithCompressionTypeAsExt := compressedFilePathWithoutExt + "." + compressUnixFiles.CompressionType.String()
	_, err = infraHelper.RunCmd(
		compressBinary,
		compressBinaryFlag,
		compressedFilePathWithCompressionTypeAsExt,
		filesToCompressStr,
	)

	compressedFilesList := compressUnixFiles.Paths

	if err != nil {
		compressedFilesList = []valueObject.UnixFilePath{}

		errMessage := fmt.Sprintf("CompressFilesError: %s", err.Error())

		log.Print(errMessage)

		for _, filePath := range compressUnixFiles.Paths {
			compressionProcessFailure := valueObject.NewCompressionProcessFailure(
				filePath,
				errMessage,
			)
			failedToCompressList = append(failedToCompressList, compressionProcessFailure)
		}
	}

	return dto.NewCompressionProcessReport(
		compressedFilesList,
		failedToCompressList,
		compressUnixFiles.DestinationPath,
	), nil
}

func (repo FilesCmdRepo) Extract(extractUnixFiles dto.ExtractUnixFiles) error {
	fileToExtract := extractUnixFiles.Path

	fileToExtractExists := infraHelper.FileExists(fileToExtract.String())
	if !fileToExtractExists {
		return errors.New("FileToExtractDoesNotExists")
	}

	destinationPath := extractUnixFiles.DestinationPath

	destinationPathExists := infraHelper.FileExists(destinationPath.String())
	if destinationPathExists {
		return errors.New("DEstinationPathAlreadyExists")
	}

	compressBinary := "tar"
	compressBinaryFlag := "-xf"
	compressDestinationFlag := "-C"

	unixFilePathExtension, _ := fileToExtract.GetFileExtension()
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

func (repo FilesCmdRepo) Delete(
	unixFilePath []valueObject.UnixFilePath,
) {
	for _, fileToDelete := range unixFilePath {
		fileExists := infraHelper.FileExists(fileToDelete.String())
		if !fileExists {
			log.Printf("DeleteFileError: FileDoesNotExists")
			continue
		}

		err := os.RemoveAll(fileToDelete.String())
		if err != nil {
			log.Printf("DeleteFileError: %s", err)
			continue
		}

		log.Printf("File '%s' deleted.", fileToDelete.String())
	}
}

func (repo FilesCmdRepo) Upload(
	unixFileDestinationPath valueObject.UnixFilePath,
	fileStreamHandler valueObject.FileStreamHandler,
) error {
	destinationFileName := unixFileDestinationPath.String() + "/" + fileStreamHandler.GetFileName().String()
	destinationEmptyFile, err := os.Create(destinationFileName)
	if err != nil {
		log.Printf("CreateEmptyFileToStoreUploadFileError: %s", err.Error())
		return errors.New("CreateEmptyFileToStoreUploadFileError")
	}
	defer destinationEmptyFile.Close()

	fileToUploadStream, err := fileStreamHandler.Open()
	if err != nil {
		log.Printf("UnableToOpenFileStream: %s", err.Error())
		return errors.New("UnableToOpenFileStream")
	}

	_, err = io.Copy(destinationEmptyFile, fileToUploadStream)
	if err != nil {
		log.Printf("CopyFileStreamHandlerContentToDestinationFileError: %s", err.Error())
		return errors.New("CopyFileStreamHandlerContentToDestinationFileError")
	}

	return nil
}
