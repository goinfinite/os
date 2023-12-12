package infra

import (
	"errors"
	"fmt"
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
