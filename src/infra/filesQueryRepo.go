package infra

import (
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

type FilesQueryRepo struct{}

func (repo FilesQueryRepo) unixFileFactory(
	filePath valueObject.UnixFilePath,
	fileInfo fs.FileInfo,
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	fileSysInfo := fileInfo.Sys().(*syscall.Stat_t)

	fileUidStr := fmt.Sprint(fileSysInfo.Uid)
	unixFileUid, err := valueObject.NewUnixUid(fileSysInfo.Uid)
	if err != nil {
		return unixFile, err
	}

	fileOwner, err := user.LookupId(fileUidStr)
	if err != nil {
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
		return unixFile, errors.New("UnableToGetFileGroupName")
	}

	unixFileGroup, err := valueObject.NewGroupName(fileGroupName.Name)
	if err != nil {
		return unixFile, err
	}

	fileAbsPathStr := filePath.String() + "/" + fileInfo.Name()
	unixFileAbsPath, err := filepath.Abs(fileAbsPathStr)
	if err != nil {
		return unixFile, errors.New("UnableToGetFileAbsolutePath")
	}

	unixFilePath, err := valueObject.NewUnixFilePath(unixFileAbsPath)
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

	if !fileInfo.IsDir() {
		mimeTypeByExtension := mime.TypeByExtension(unixFileExtension.String())
		unixFileMimeType, err = valueObject.NewMimeType(mimeTypeByExtension)
		if err != nil {
			return unixFile, err
		}
	}

	filePermissions := fileInfo.Mode().Perm()
	filePermissionsStr := fmt.Sprintf("%o", filePermissions)
	unixFilePermissions, err := valueObject.NewUnixFilePermissions(filePermissionsStr)
	if err != nil {
		return unixFile, err
	}

	unixFileSize := valueObject.Byte(fileInfo.Size())
	unixFileUpdatedAt := valueObject.UnixTime(fileInfo.ModTime().Unix())

	unixFile = entity.NewUnixFile(
		unixFileUid,
		unixFileUsername,
		unixFileGid,
		unixFileGroup,
		unixFileMimeType,
		unixFileName,
		unixFilePath,
		&unixFileExtension,
		unixFilePermissions,
		unixFileSize,
		unixFileUpdatedAt,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) Get(
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	unixFileSlice := []entity.UnixFile{}

	filePathInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		return unixFileSlice, errors.New("UnableToGetPathInfo")
	}

	unixPathIsDir := filePathInfo.IsDir()
	if !unixPathIsDir {
		unixFile, err := repo.unixFileFactory(unixFilePath, filePathInfo)
		if err == nil {
			unixFileSlice = append(unixFileSlice, unixFile)
		}
	}

	var filePathToAnalyzeList []fs.FileInfo
	if unixPathIsDir {
		unixDirInodesToAnalyzeSlice, err := os.ReadDir(unixFilePath.String())
		if err != nil {
			return unixFileSlice, errors.New("UnableToGetDirInfo")
		}

		for _, dirEntry := range unixDirInodesToAnalyzeSlice {
			dirInfo, _ := dirEntry.Info()
			filePathToAnalyzeList = append(filePathToAnalyzeList, dirInfo)
		}
	}

	for _, pathInfo := range filePathToAnalyzeList {
		unixFile, err := repo.unixFileFactory(unixFilePath, pathInfo)
		if err != nil {
			continue
		}

		unixFileSlice = append(unixFileSlice, unixFile)
	}

	return unixFileSlice, nil
}
