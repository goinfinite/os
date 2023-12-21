package infra

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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

	fileAbsPathStr := filePath.String() + "/" + fileInfo.Name()
	unixFileAbsPath, err := filepath.Abs(fileAbsPathStr)
	if err != nil {
		log.Printf("UnableToGetFileAbsolutePath: %s", err)
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
		unixFileName,
		unixFilePath,
		unixFileMimeType,
		unixFilePermissions,
		unixFileSize,
		&unixFileExtension,
		unixFileUid,
		unixFileUsername,
		unixFileGid,
		unixFileGroup,
		unixFileUpdatedAt,
		unixFileStreamPtr,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) Exists(
	unixFilePath valueObject.UnixFilePath,
) (bool, error) {
	_, err := os.Stat(unixFilePath.String())
	if os.IsNotExist(err) {
		log.Printf("PathDoesNotExists: %s", unixFilePath.String())
		return false, nil
	}
	if err != nil {
		log.Printf("PathExistsError: %s", err.Error())
		return false, errors.New("PathExistsError")
	}

	return true, nil
}

func (repo FilesQueryRepo) IsDir(
	unixFilePath valueObject.UnixFilePath,
) (bool, error) {
	unixFileInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		log.Printf("PathIsDirError: %s", err.Error())
		return false, errors.New("PathIsDirError")
	}

	return unixFileInfo.IsDir(), nil
}

func (repo FilesQueryRepo) Get(
	unixFilePath valueObject.UnixFilePath,
) ([]entity.UnixFile, error) {
	unixFileList := []entity.UnixFile{}

	exists, err := repo.Exists(unixFilePath)
	if err != nil {
		return unixFileList, err
	}
	if !exists {
		return unixFileList, errors.New("PathDoesNotExists")
	}

	isDir, err := repo.IsDir(unixFilePath)
	if err != nil {
		return unixFileList, err
	}
	if !isDir {
		fileInfo, err := os.Stat(unixFilePath.String())
		if err != nil {
			return unixFileList, errors.New("UnableToGetPathInfo")
		}

		unixFile, err := repo.unixFileFactory(unixFilePath, fileInfo)
		if err == nil {
			unixFileList = append(unixFileList, unixFile)
		}

		return unixFileList, nil
	}

	dirEntriesToAnalyzeList, err := os.ReadDir(unixFilePath.String())
	if err != nil {
		log.Printf("UnableToGetDirInfo: %s", err)
		return unixFileList, errors.New("UnableToGetDirInfo")
	}

	for _, dirEntry := range dirEntriesToAnalyzeList {
		inodeInfo, _ := dirEntry.Info()
		unixFile, err := repo.unixFileFactory(unixFilePath, inodeInfo)
		if err != nil {
			continue
		}

		unixFileList = append(unixFileList, unixFile)
	}

	return unixFileList, nil
}

func (repo FilesQueryRepo) GetOnlyFile(
	unixFilePath valueObject.UnixFilePath,
) (entity.UnixFile, error) {
	var unixFile entity.UnixFile

	exists, err := repo.Exists(unixFilePath)
	if err != nil {
		return unixFile, err
	}
	if !exists {
		return unixFile, errors.New("FileDoesNotExists")
	}

	isDir, err := repo.IsDir(unixFilePath)
	if err != nil {
		return unixFile, err
	}
	if isDir {
		return unixFile, errors.New("PathIsNotAFile")
	}

	fileInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		return unixFile, errors.New("UnableToGetFileInfo")
	}

	return repo.unixFileFactory(unixFilePath, fileInfo)
}
