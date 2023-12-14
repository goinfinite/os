package infraFactory

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"mime"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/valueObject"
)

func UnixFileFactory(
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

	unixFileMimeType, _ := valueObject.NewMimeType("directory")

	if !fileInfo.IsDir() {
		mimeType := "generic"

		mimeTypeWithCharset := mime.TypeByExtension("." + unixFileExtension.String())
		if len(mimeTypeWithCharset) > 1 {
			mimeTypeOnly := strings.Split(mimeTypeWithCharset, ";")[0]
			mimeType = mimeTypeOnly
		}

		unixFileMimeType, err = valueObject.NewMimeType(mimeType)
		if err != nil {
			log.Print(err)
			return unixFile, err
		}
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
