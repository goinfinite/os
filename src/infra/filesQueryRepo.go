package infra

import (
	"errors"
	"fmt"
	"io/fs"
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
	unixFileSlice := []entity.UnixFile{}

	var unixPathInfoSlice []fs.FileInfo

	unixPathInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		return unixFileSlice, errors.New("UnableToGetPathInfo")
	}

	unixPathInfoSlice = append(unixPathInfoSlice, unixPathInfo)

	unixFileIsDir := unixPathInfo.IsDir()
	if unixFileIsDir {
		unixPathInfoSlice = unixPathInfoSlice[1:]

		unixDirEntriesSlice, err := os.ReadDir(unixFilePath.String())
		if err != nil {
			return unixFileSlice, errors.New("UnableToGetDirInfo")
		}

		for _, dirEntry := range unixDirEntriesSlice {
			dirInfo, _ := dirEntry.Info()
			unixPathInfoSlice = append(unixPathInfoSlice, dirInfo)
		}
	}

	for _, pathInfo := range unixPathInfoSlice {
		filePath := unixFilePath.String() + "/" + pathInfo.Name()

		unixFileAbsPath, err := filepath.Abs(filePath)
		if err != nil {
			return unixFileSlice, errors.New("UnableToGetFileAbsolutePath")
		}

		unixFilePermissions := pathInfo.Mode().Perm()
		unixFilePermissionsStr := fmt.Sprintf("%o", unixFilePermissions)

		unixFileSysInfo := pathInfo.Sys().(*syscall.Stat_t)

		unixFileUidStr := fmt.Sprint(unixFileSysInfo.Uid)
		unixFileOwner, err := user.LookupId(unixFileUidStr)
		if err != nil {
			return unixFileSlice, errors.New("UnableToGetFileGroupName")
		}

		unixFileGidStr := fmt.Sprint(unixFileSysInfo.Gid)
		unixFileGroup, err := user.LookupGroupId(unixFileGidStr)
		if err != nil {
			return unixFileSlice, errors.New("UnableToGetFileGroupName")
		}

		unixFile, err := repo.unixFileFactory(
			unixFileIsDir,
			unixFileAbsPath,
			unixFilePermissionsStr,
			pathInfo.Size(),
			pathInfo.ModTime().Unix(),
			unixFileSysInfo.Uid,
			unixFileSysInfo.Gid,
			unixFileOwner.Username,
			unixFileGroup.Name,
		)

		unixFileSlice = append(unixFileSlice, unixFile)
	}

	return unixFileSlice, nil
}
