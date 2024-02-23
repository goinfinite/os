package filesInfra

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type FilesQueryRepo struct{}

func (repo FilesQueryRepo) unixFileFactory(
	filePath valueObject.UnixFilePath,
	shouldReturnContent bool,
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	fileInfo, err := os.Stat(filePath.String())
	if err != nil {
		return unixFile, err
	}

	fileSysInfo := fileInfo.Sys().(*syscall.Stat_t)

	unixFileUid, err := valueObject.NewUnixUid(fileSysInfo.Uid)
	if err != nil {
		return unixFile, err
	}

	fileOwner, err := user.LookupId(unixFileUid.String())
	if err != nil {
		return unixFile, err
	}

	unixFileUsername, err := valueObject.NewUsername(fileOwner.Username)
	if err != nil {
		return unixFile, err
	}

	unixFileGid, err := valueObject.NewGroupId(fileSysInfo.Gid)
	if err != nil {
		return unixFile, err
	}

	fileGroupName, err := user.LookupGroupId(unixFileGid.String())
	if err != nil {
		return unixFile, err
	}

	unixFileGroup, err := valueObject.NewGroupName(fileGroupName.Name)
	if err != nil {
		return unixFile, err
	}

	unixFileAbsPath, err := filepath.Abs(filePath.String())
	if err != nil {
		return unixFile, err
	}

	unixFilePath, err := valueObject.NewUnixFilePath(unixFileAbsPath)
	if err != nil {
		return unixFile, err
	}

	var unixFileExtensionPtr *valueObject.UnixFileExtension
	unixFileExtension, err := unixFilePath.GetFileExtension()
	if err == nil {
		unixFileExtensionPtr = &unixFileExtension
	}

	unixFileMimeType := unixFileExtension.GetMimeType()
	if fileInfo.IsDir() {
		unixFileMimeType, _ = valueObject.NewMimeType("directory")
		unixFileExtensionPtr = nil
	}

	filePermissions := fileInfo.Mode().Perm()
	filePermissionsStr := fmt.Sprintf("%o", filePermissions)
	unixFilePermissions, err := valueObject.NewUnixFilePermissions(filePermissionsStr)
	if err != nil {
		return unixFile, err
	}

	unixFileSize := valueObject.Byte(fileInfo.Size())

	var unixFileContentPtr *valueObject.UnixFileContent
	if shouldReturnContent && unixFileSize.ToMiB() <= 5 {
		unixFileContentStr, err := infraHelper.GetFileContent(filePath.String())
		if err != nil {
			return unixFile, errors.New("FailedToGetFileContent: " + err.Error())
		}

		unixFileContent, err := valueObject.NewUnixFileContent(unixFileContentStr)
		if err != nil {
			return unixFile, err
		}

		unixFileContentPtr = &unixFileContent
	}

	unixFileUpdatedAt := valueObject.UnixTime(fileInfo.ModTime().Unix())

	unixFile = entity.NewUnixFile(
		unixFilePath.GetFileName(),
		unixFilePath,
		unixFileMimeType,
		unixFilePermissions,
		unixFileSize,
		unixFileExtensionPtr,
		unixFileContentPtr,
		unixFileUid,
		unixFileUsername,
		unixFileGid,
		unixFileGroup,
		unixFileUpdatedAt,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) Get(
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	unixFileList := []entity.UnixFile{}

	exists := infraHelper.FileExists(unixFilePath.String())
	if !exists {
		return unixFileList, errors.New("PathNotFound")
	}

	filePathEndsWithSlash := strings.HasSuffix(unixFilePath.String(), "/")
	if filePathEndsWithSlash {
		filePathWithoutSlashAtTheEnd := strings.TrimSuffix(unixFilePath.String(), "/")
		unixFilePath, _ = valueObject.NewUnixFilePath(filePathWithoutSlashAtTheEnd)
	}

	filesToFactory := []valueObject.UnixFilePath{
		unixFilePath,
	}

	fileInfo, _ := os.Stat(unixFilePath.String())
	if fileInfo.IsDir() {
		filesToFactoryWithoutDir := filesToFactory[1:]
		filesToFactory = filesToFactoryWithoutDir

		rawDirectoryFiles, err := infraHelper.RunCmd(
			"find",
			unixFilePath.String(),
			"-maxdepth",
			"1",
			"-printf",
			"%p\n",
		)
		if err != nil {
			return unixFileList, err
		}
		if len(rawDirectoryFiles) == 0 {
			return unixFileList, errors.New("UnableToGetDirFiles")
		}

		rawDirectoryFilesList := strings.Split(rawDirectoryFiles, "\n")
		for _, fileToFactoryStr := range rawDirectoryFilesList {
			filePath, err := valueObject.NewUnixFilePath(fileToFactoryStr)
			if err != nil {
				log.Printf(
					"FileToFactoryError (%s): %s",
					filePath.String(),
					err.Error(),
				)
				continue
			}

			filesToFactory = append(filesToFactory, filePath)
		}
	}

	shouldReturnContent := false
	if len(filesToFactory) == 1 {
		shouldReturnContent = true
	}

	for _, fileToFactory := range filesToFactory {
		filePathIsDir := fileInfo.IsDir() && (fileToFactory.String() == unixFilePath.String())
		if filePathIsDir {
			continue
		}

		unixFile, err := repo.unixFileFactory(fileToFactory, shouldReturnContent)

		if err != nil {
			log.Printf(
				"UnixFileFactoryError (%s): %s",
				fileToFactory.String(),
				err.Error(),
			)
			continue
		}

		unixFileList = append(unixFileList, unixFile)
	}

	return unixFileList, nil
}

func (repo FilesQueryRepo) GetOne(
	unixFilePath valueObject.UnixFilePath,
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	exists := infraHelper.FileExists(unixFilePath.String())
	if !exists {
		return unixFile, errors.New("FileNotFound")
	}

	shouldReturnContent := false
	return repo.unixFileFactory(unixFilePath, shouldReturnContent)
}
