package infra

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
)

type FilesQueryRepo struct{}

func (repo FilesQueryRepo) unixFileFactory(
	isDir bool,
	absFilePath string,
	filePermissions string,
	fileSizeInBytes int64,
	fileModDate int64,
	fileUid uint32,
	fileGid uint32,
	fileOwner string,
	fileGroupName string,
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	unixFilePath, err := valueObject.NewUnixFilePath(absFilePath)
	if err != nil {
		return unixFile, err
	}

	unixFileName, err := unixFilePath.GetFileName()
	if err != nil {
		return unixFile, err
	}

	unixFileExtension, err := unixFilePath.GetFileExtension()
	if err != nil {
		return unixFile, err
	}

	unixFileMimeType, _ := valueObject.NewMimeType("directory")

	if !isDir {
		unixFileMimeType = infraHelper.GetFileExtensionMimeType(unixFileExtension)
	}

	unixFilePermissions, err := valueObject.NewUnixFilePermissions(filePermissions)
	if err != nil {
		return unixFile, err
	}

	unixFileSize := valueObject.Byte(fileSizeInBytes)
	unixFileUpdatedAt := valueObject.UnixTime(fileModDate)

	unixFileUidInt := int(fileUid)
	unixFileUid, err := valueObject.NewUnixUid(unixFileUidInt)
	if err != nil {
		return unixFile, err
	}

	unixFileGid, err := valueObject.NewGroupId(fileGid)
	if err != nil {
		return unixFile, err
	}

	unixFileUsername, err := valueObject.NewUsername(fileOwner)
	if err != nil {
		return unixFile, err
	}

	unixFileGroup, err := valueObject.NewGroupName(fileGroupName)
	if err != nil {
		return unixFile, err
	}

	unixFile = entity.NewUnixFile(
		unixFileUid,
		unixFileGid,
		unixFileMimeType,
		unixFileName,
		unixFilePath,
		&unixFileExtension,
		unixFilePermissions,
		unixFileSize,
		unixFileUpdatedAt,
		unixFileUsername,
		unixFileGroup,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) Get(
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	unixFileList := []entity.UnixFile{}

	unixFileInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		return unixFileList, errors.New("UnableToOpenFile")
	}

	unixFileIsDir := unixFileInfo.IsDir()

	unixFileAbsPath, err := filepath.Abs(unixFilePath.String())
	if err != nil {
		return unixFileList, errors.New("UnableToGetFileAbsolutePath")
	}

	unixFilePermissions := unixFileInfo.Mode().Perm()
	unixFilePermissionsStr := fmt.Sprintf("%o", unixFilePermissions)

	unixFileSysInfo := unixFileInfo.Sys().(*syscall.Stat_t)

	unixFileUidStr := fmt.Sprint(unixFileSysInfo.Uid)
	unixFileOwner, err := user.LookupId(unixFileUidStr)
	if err != nil {
		return unixFileList, errors.New("UnableToGetFileGroupName")
	}

	unixFileGidStr := fmt.Sprint(unixFileSysInfo.Gid)
	unixFileGroup, err := user.LookupGroupId(unixFileGidStr)
	if err != nil {
		return unixFileList, errors.New("UnableToGetFileGroupName")
	}

	unixFile, err := repo.unixFileFactory(
		unixFileIsDir,
		unixFileAbsPath,
		unixFilePermissionsStr,
		unixFileInfo.Size(),
		unixFileInfo.ModTime().Unix(),
		unixFileSysInfo.Uid,
		unixFileSysInfo.Gid,
		unixFileOwner.Username,
		unixFileGroup.Name,
	)

	unixFileList = append(unixFileList, unixFile)

	return unixFileList, nil
}
