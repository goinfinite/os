package infra

import (
	"bufio"
	"compress/gzip"
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
		log.Printf("CreateUnixFileError: %s", err)
		return errors.New("CreateUnixFileError")
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
		fileIsDir, _ := originPath.IsDir()
		if fileIsDir {
			fileType = "Directory"
		}

		moveErrorStr := fmt.Sprintf("MoveUnix%sError", fileType)

		log.Printf("%s: %s", moveErrorStr, err)
		return errors.New(moveErrorStr)
	}

	return nil
}

func (repo FilesCmdRepo) Copy(addUnixFileCopy dto.AddUnixFileCopy) error {
	_, err := infraHelper.RunCmd(
		"rsync",
		"-avq",
		addUnixFileCopy.OriginPath.String(),
		addUnixFileCopy.DestinationPath.String(),
	)
	if err != nil {
		fileType := "File"
		fileIsDir, _ := addUnixFileCopy.OriginPath.IsDir()
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
		fileIsDir, _ := unixFilePath.IsDir()
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
	unixFilePath valueObject.UnixFilePath,
	unixFileDestinationPath valueObject.UnixFilePath,
	unixCompressionType valueObject.UnixCompressionType,
) error {
	fileToCompress, err := os.Open(unixFilePath.String())
	if err != nil {
		log.Printf("OpenFileToCompressError: %s", err)
		return errors.New("OpenFileToCompressError")
	}

	fileToCompressReader := bufio.NewReader(fileToCompress)

	fileToCompressBytes, err := io.ReadAll(fileToCompressReader)
	if err != nil {
		log.Printf("ReadFileToCompressBytesError: %s", err)
		return errors.New("ReadFileToCompressBytesError")
	}

	compressedFilePathWithoutExt := strings.Split(unixFileDestinationPath.String(), ".")[0]
	compressedFilePathWithCompressionTypeAsExt := compressedFilePathWithoutExt + "." + unixCompressionType.String()

	compressedFile, err := os.Create(compressedFilePathWithCompressionTypeAsExt)
	if err != nil {
		log.Printf("CreateCompressedEmptyFileError: %s", err)
		return errors.New("CreateCompressedEmptyFileError")
	}

	gzipWriter := gzip.NewWriter(compressedFile)
	defer gzipWriter.Close()

	_, err = gzipWriter.Write(fileToCompressBytes)
	if err != nil {
		log.Printf("WriteFileToCompressBytesInCompressedFileError: %s", err)
		return errors.New("WriteFileToCompressBytesInCompressedFileError")
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
