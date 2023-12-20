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

type FilesCmdRepo struct {
	filesQueryRepo FilesQueryRepo
}

func NewFilesCmdRepo() FilesCmdRepo {
	return FilesCmdRepo{
		filesQueryRepo: FilesQueryRepo{},
	}
}

func (repo FilesCmdRepo) Add(addUnixFile dto.AddUnixFile) error {
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

func (repo FilesCmdRepo) Move(
	originPath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) error {
	err := os.Rename(
		originPath.String(),
		destinationPath.String(),
	)
	if err != nil {
		fileType := "File"
		fileIsDir, _ := repo.filesQueryRepo.IsDir(originPath)
		if fileIsDir {
			fileType = "Directory"
		}

		moveErrorStr := fmt.Sprintf("MoveUnix%sError", fileType)

		log.Printf("%s: %s", moveErrorStr, err)
		return errors.New(moveErrorStr)
	}

	return nil
}

func (repo FilesCmdRepo) Copy(copyUnixFile dto.CopyUnixFile) error {
	_, err := infraHelper.RunCmd(
		"rsync",
		"-avq",
		copyUnixFile.OriginPath.String(),
		copyUnixFile.DestinationPath.String(),
	)
	if err != nil {
		fileType := "File"
		fileIsDir, _ := repo.filesQueryRepo.IsDir(copyUnixFile.OriginPath)
		if fileIsDir {
			fileType = "Directory"
		}

		moveErrorStr := fmt.Sprintf("CopyUnix%sError", fileType)

		log.Printf("%s: %s", moveErrorStr, err)
		return errors.New(moveErrorStr)
	}

	return nil
}

func (repo FilesCmdRepo) UpdateContent(
	updateUnixFileContent dto.UpdateUnixFileContent,
) error {
	file, err := os.OpenFile(updateUnixFileContent.Path.String(), os.O_WRONLY, 0777)
	if err != nil {
		log.Printf("OpenFileError: %s", err)
		return errors.New("OpenFileError")
	}
	defer file.Close()

	err = file.Truncate(0)
	if err != nil {
		log.Printf("TruncateFileError: %s", err)
		return errors.New("TruncateFileError")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		log.Printf("SeekFileError: %s", err)
		return errors.New("SeekFileError")
	}

	_, err = file.WriteString(updateUnixFileContent.Content.GetDecodedContent())
	if err != nil {
		log.Printf("WriteFileError: %s", err)
		return errors.New("WriteFileError")
	}

	err = file.Sync()
	if err != nil {
		log.Printf("FileSyncError: %s", err)
		return errors.New("FileSyncError")
	}

	return nil
}

func (repo FilesCmdRepo) UpdatePermissions(
	unixFilePath valueObject.UnixFilePath,
	unixFilePermissions valueObject.UnixFilePermissions,
) error {
	err := os.Chmod(unixFilePath.String(), unixFilePermissions.GetFileMode())
	if err != nil {
		fileType := "File"
		fileIsDir, _ := repo.filesQueryRepo.IsDir(unixFilePath)
		if fileIsDir {
			fileType = "Directory"
		}

		chmodErrorStr := fmt.Sprintf("ChmodUnix%sError", fileType)

		log.Printf("%s: %s", chmodErrorStr, err)
		return errors.New(chmodErrorStr)
	}

	return nil
}

func (repo FilesCmdRepo) Compress(
	unixFilePaths []valueObject.UnixFilePath,
	unixFileDestinationPath valueObject.UnixFilePath,
	unixCompressionType valueObject.UnixCompressionType,
) error {
	compressBinary := "tar"
	compressBinaryFlag := "-czf"
	if unixCompressionType.String() == "zip" {
		compressBinary = "zip"
		compressBinaryFlag = "-qr"
	}

	filesToCompressStr := unixFilePaths[0].String()
	if len(unixFilePaths) > 1 {
		var filesToCompressStrSlice []string
		for _, filePath := range unixFilePaths {
			filesToCompressStrSlice = append(filesToCompressStrSlice, filePath.String())
		}

		filesToCompressStr = strings.Join(filesToCompressStrSlice, " ")
	}

	compressedFilePathWithoutExt := strings.Split(unixFileDestinationPath.String(), ".")[0]
	compressedFilePathWithCompressionTypeAsExt := compressedFilePathWithoutExt + "." + unixCompressionType.String()
	_, err := infraHelper.RunCmd(
		compressBinary,
		compressBinaryFlag,
		compressedFilePathWithCompressionTypeAsExt,
		filesToCompressStr,
	)

	if err != nil {
		log.Printf("CompressFilesError: %s", err.Error())
		return errors.New("CompressFilesError")
	}

	return nil
}

func (repo FilesCmdRepo) Extract(
	unixFilePath valueObject.UnixFilePath,
	unixFileDestinationPath valueObject.UnixFilePath,
) error {
	compressBinary := "tar"
	compressBinaryFlag := "-xf"
	compressDestinationFlag := "-C"

	unixFilePathExtension, _ := unixFilePath.GetFileExtension()
	if unixFilePathExtension.String() == "zip" {
		compressBinary = "unzip"
		compressBinaryFlag = "-qq"
		compressDestinationFlag = "-d"
	}

	err := infraHelper.MakeDir(unixFileDestinationPath.String())
	if err != nil {
		log.Printf("CreateExtractDirectoryError: %s", err.Error())
		return errors.New("CreateExtractDirectoryError")
	}

	_, err = infraHelper.RunCmd(
		compressBinary,
		compressBinaryFlag,
		unixFilePath.String(),
		compressDestinationFlag,
		unixFileDestinationPath.String(),
	)

	if err != nil {
		log.Printf("ExtractFilesError: %s", err.Error())
		return errors.New("ExtractFilesError")
	}

	return nil
}

func (repo FilesCmdRepo) Delete(
	unixFilePath valueObject.UnixFilePath,
) error {
	err := os.RemoveAll(unixFilePath.String())
	if err != nil {
		log.Printf("DeleteFileError: %s", err)
		return errors.New("DeleteFileError")
	}

	return nil
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

	fileStreamHandlerInstance, err := fileStreamHandler.Open()

	_, err = io.Copy(destinationEmptyFile, fileStreamHandlerInstance)
	if err != nil {
		log.Printf("CopyFileStreamHandlerContentToDestinationFileError: %s", err.Error())
		return errors.New("CopyFileStreamHandlerContentToDestinationFileError")
	}

	return nil
}
