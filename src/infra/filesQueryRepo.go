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

// ENHANCEMENT: Adicionar as outras flags de find no RunCmd e usar o retorno para montar o UnixFile dentro da factory.
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

	var unixFileExtensionPtr *valueObject.UnixFileExtension
	unixFileExtension, err := unixFilePath.GetFileExtension()
	unixFileExtensionPtr = &unixFileExtension
	if err != nil {
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
		unixFileName,
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

// TODO: Remover esse método e alterar todos os lugares onde chamam ele para usar para usar o Get ou o GetOnlyFile.
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

// TODO: Remover o IsDir() e utilizar o Get e o GetOnlyFile para validar se é diretório ou não.
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

	exists := infraHelper.FileExists(unixFilePath.String())
	if !exists {
		return unixFileList, errors.New("PathDoesNotExists")
	}

	filesToFactory := []valueObject.UnixFilePath{
		unixFilePath,
	}

	isDir, err := repo.IsDir(unixFilePath)
	if err != nil {
		return unixFileList, err
	}

	if isDir {
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
