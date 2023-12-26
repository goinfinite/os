package infra

import (
	"errors"
	"fmt"
	"log"
	"mime"
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
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	fileInfo, err := os.Stat(filePath.String())
	if err != nil {
		return unixFile, err
	}

	fileSysInfo := fileInfo.Sys().(*syscall.Stat_t)

	fileUidStr := fmt.Sprint(fileSysInfo.Uid)
	unixFileUid, err := valueObject.NewUnixUid(fileSysInfo.Uid)
	if err != nil {
		return unixFile, err
	}

	fileOwner, err := user.LookupId(fileUidStr)
	if err != nil {
		log.Printf("UnableToGetFileGroupName: %s", err)
		return unixFile, errors.New("UnableToGetFileGroupName")
	}

	unixFileUsername, err := valueObject.NewUsername(fileOwner.Username)
	if err != nil {
		return unixFile, err
	}

	fileGidStr := fmt.Sprint(fileSysInfo.Gid)
	unixFileGid, err := valueObject.NewGroupId(fileSysInfo.Gid)
	if err != nil {
		return unixFile, err
	}

	fileGroupName, err := user.LookupGroupId(fileGidStr)
	if err != nil {
		log.Printf("UnableToGetFileGroupName: %s", err)
		return unixFile, errors.New("UnableToGetFileGroupName")
	}

	unixFileGroup, err := valueObject.NewGroupName(fileGroupName.Name)
	if err != nil {
		return unixFile, err
	}

	unixFileAbsPath, err := filepath.Abs(filePath.String())
	if err != nil {
		log.Printf("UnableToGetFileAbsolutePath: %s", err)
		return unixFile, errors.New("UnableToGetFileAbsolutePath")
	}

	unixFilePath, err := valueObject.NewUnixFilePath(unixFileAbsPath)
	if err != nil {
		return unixFile, err
	}

	var unixFileExtensionPtr *valueObject.UnixFileExtension
	unixFileExtension := unixFilePath.GetFileExtension()
	unixFileExtensionPtr = &unixFileExtension
	if unixFileExtension.IsEmpty() {
		unixFileExtensionPtr = nil
	}

	mimeTypeStr := "directory"
	isDir := fileInfo.IsDir()
	if !isDir {
		mimeTypeStr = "generic"

		mimeTypeWithCharset := mime.TypeByExtension("." + unixFileExtension.String())
		if len(mimeTypeWithCharset) > 1 {
			mimeTypeOnly := strings.Split(mimeTypeWithCharset, ";")[0]
			mimeTypeStr = mimeTypeOnly
		}
	}

	unixFileMimeType, err := valueObject.NewMimeType(mimeTypeStr)
	if err != nil {
		log.Print(err)
		return unixFile, err
	}

	filePermissions := fileInfo.Mode().Perm()
	filePermissionsStr := fmt.Sprintf("%o", filePermissions)
	unixFilePermissions, err := valueObject.NewUnixFilePermissions(filePermissionsStr)
	if err != nil {
		log.Print(err)
		return unixFile, err
	}

	unixFileSize := valueObject.Byte(fileInfo.Size())
	unixFileUpdatedAt := valueObject.UnixTime(fileInfo.ModTime().Unix())

	unixFileStreamPtr, err := os.Open(filePath.String())
	if err != nil {
		log.Printf("OpenFileError: %s", err.Error())
		return unixFile, err
	}
	defer unixFileStreamPtr.Close()

	unixFile = entity.NewUnixFile(
		unixFilePath.GetFileName(),
		unixFilePath,
		unixFileMimeType,
		unixFilePermissions,
		unixFileSize,
		unixFileExtensionPtr,
		unixFileUid,
		unixFileUsername,
		unixFileGid,
		unixFileGroup,
		unixFileUpdatedAt,
		unixFileStreamPtr,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) Get(
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	unixFileList := []entity.UnixFile{}

	exists := infraHelper.FileExists(unixFilePath.String())
	if !exists {
		return unixFileList, errors.New("PathDoesNotExists")
	}

	filesToFactory := []valueObject.UnixFilePath{
		unixFilePath,
	}

	fileInfo, _ := os.Stat(unixFilePath.String())
	if fileInfo.IsDir() {
		filesToFactoryWithoutDir := filesToFactory[1:]
		filesToFactory = filesToFactoryWithoutDir

		filesToFactoryStr, err := infraHelper.RunCmd(
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
		if len(filesToFactoryStr) == 0 {
			return unixFileList, errors.New("UnableToGetDirFiles")
		}

		filesToFactoryStrList := strings.Split(filesToFactoryStr, "\n")
		for _, fileToFactoryStr := range filesToFactoryStrList {
			filePath, err := valueObject.NewUnixFilePath(fileToFactoryStr)
			if err != nil {
				log.Printf("FileToFactoryError: %s", err.Error())
				continue
			}

			filesToFactory = append(filesToFactory, filePath)
		}
	}

	for _, fileToFactory := range filesToFactory {
		unixFile, err := repo.unixFileFactory(fileToFactory)

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

func (repo FilesQueryRepo) GetOnly(
	unixFilePath valueObject.UnixFilePath,
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	exists := infraHelper.FileExists(unixFilePath.String())
	if !exists {
		return unixFile, errors.New("FileDoesNotExists")
	}

	return repo.unixFileFactory(unixFilePath)
}
