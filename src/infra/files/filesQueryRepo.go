package filesInfra

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
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
		unixFileMimeType = valueObject.DirectoryMimeType
		unixFileExtensionPtr = nil
	}

	filePermissions := fileInfo.Mode().Perm()
	filePermissionsStr := fmt.Sprintf("%o", filePermissions)
	unixFilePermissions, err := valueObject.NewUnixFilePermissions(filePermissionsStr)
	if err != nil {
		return unixFile, err
	}

	unixFileSize, err := valueObject.NewByte(fileInfo.Size())
	if err != nil {
		return unixFile, err
	}

	var unixFileContentPtr *valueObject.UnixFileContent
	if shouldReturnContent && unixFileSize.ToMiB() <= valueObject.FileContentMaxSizeInMb {
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

	unixFileUpdatedAt := valueObject.NewUnixTimeWithGoTime(fileInfo.ModTime())

	unixFile = entity.NewUnixFile(
		unixFilePath.GetFileName(), unixFilePath, unixFileMimeType, unixFilePermissions,
		unixFileSize, unixFileExtensionPtr, unixFileContentPtr, unixFileUid,
		unixFileUsername, unixFileGid, unixFileGroup, unixFileUpdatedAt,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) simplifiedUnixFileFactory(
	unixFilePath valueObject.UnixFilePath,
) (simplifiedUnixFile entity.SimplifiedUnixFile, err error) {
	unixFilePathStr := unixFilePath.String()
	if !unixFilePath.IsRootPath() && strings.HasSuffix(unixFilePathStr, "/") {
		unixFilePathWithoutTrailingSlash := strings.TrimSuffix(unixFilePathStr, "/")
		unixFilePath, err = valueObject.NewUnixFilePath(
			unixFilePathWithoutTrailingSlash,
		)
		if err != nil {
			return simplifiedUnixFile, err
		}

		unixFilePathStr = unixFilePath.String()
	}

	fileInfo, err := os.Stat(unixFilePathStr)
	if err != nil {
		return simplifiedUnixFile, err
	}

	unixFileMimeType := valueObject.GenericMimeType
	if fileInfo.IsDir() {
		unixFileMimeType = valueObject.DirectoryMimeType
	}

	unixFileExtension, err := unixFilePath.GetFileExtension()
	if err == nil {
		unixFileMimeType = unixFileExtension.GetMimeType()
	}

	return entity.NewSimplifiedUnixFile(
		unixFilePath.GetFileName(), unixFilePath, unixFileMimeType,
	), nil
}

func (repo FilesQueryRepo) readUnixFileTree(
	desiredAbsolutePath valueObject.UnixFilePath,
	desiredParentPath valueObject.UnixFilePath,
) (unixFileTree dto.UnixFileTree, err error) {
	unixFileTreeRoot, err := repo.simplifiedUnixFileFactory(desiredParentPath)
	if err != nil {
		return unixFileTree, err
	}
	unixFileTree = dto.NewUnixFileTree(unixFileTreeRoot, []dto.UnixFileTree{})

	desiredParentPathStr := desiredParentPath.String()
	currentDesiredSourcePath := strings.TrimPrefix(
		desiredAbsolutePath.String(), desiredParentPathStr,
	)
	if currentDesiredSourcePath == "" {
		return unixFileTree, err
	}

	currentParentSourcePath := strings.Split(currentDesiredSourcePath, "/")[0]
	if currentParentSourcePath == "" {
		return unixFileTree, err
	}

	nextDesiredParentPath, err := valueObject.NewUnixFilePath(
		desiredParentPathStr + currentParentSourcePath + "/",
	)
	nextDesiredParentPathStr := nextDesiredParentPath.String()

	nextDesiredParentPathWithoutLeadingSlash := strings.TrimSuffix(
		nextDesiredParentPathStr, "/",
	)
	fileStat, err := os.Stat(nextDesiredParentPathWithoutLeadingSlash)
	if err != nil {
		return unixFileTree, err
	}
	nextDesiredParentPathIsDir := fileStat.IsDir()

	findCmdArgs := []string{
		"-L", desiredParentPathStr, "-mindepth", "1", "-maxdepth", "1",
	}
	if nextDesiredParentPathIsDir {
		findCmdArgs = append(
			findCmdArgs, "-type", "d", "!", "-path",
			nextDesiredParentPathWithoutLeadingSlash,
		)
	}

	rawUnixFileTree, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find",
		Args:    findCmdArgs,
	})
	if err != nil {
		return unixFileTree, err
	}

	for _, rawUnixFilePath := range strings.Split(rawUnixFileTree, "\n") {
		if rawUnixFilePath == "" {
			continue
		}

		unixFilePath, err := valueObject.NewUnixFilePath(rawUnixFilePath)
		if err != nil {
			slog.Debug(err.Error(), slog.String("rawUnixFilePath", rawUnixFilePath))
			continue
		}

		simplifiedFileEntity, err := repo.simplifiedUnixFileFactory(unixFilePath)
		if err != nil {
			slog.Debug(err.Error(), slog.String("unixFilePath", rawUnixFilePath))
			continue
		}

		unixFileTree.AddUnixFile(simplifiedFileEntity)
	}

	if !nextDesiredParentPathIsDir {
		return unixFileTree, err
	}

	unixFileSubTree, err := repo.readUnixFileTree(
		desiredAbsolutePath, nextDesiredParentPath,
	)
	if err != nil {
		return unixFileTree, err
	}

	unixFileTree.AddSubTree(unixFileSubTree)
	return unixFileTree, err
}

func (repo FilesQueryRepo) Read(
	requestDto dto.ReadFilesRequest,
) (responseDto dto.ReadFilesResponse, err error) {
	sourcePathStr := requestDto.SourcePath.String()
	exists := infraHelper.FileExists(sourcePathStr)
	if !exists {
		return responseDto, errors.New("PathNotFound")
	}

	filePathHasTrailingSlash := strings.HasSuffix(sourcePathStr, "/")
	if filePathHasTrailingSlash && !requestDto.SourcePath.IsRootPath() {
		filePathWithoutTrailingSlashStr := strings.TrimSuffix(sourcePathStr, "/")
		filePathWithoutTrailingSlash, _ := valueObject.NewUnixFilePath(
			filePathWithoutTrailingSlashStr,
		)
		sourcePathStr = filePathWithoutTrailingSlash.String()
	}

	sourcePathInfo, err := os.Stat(sourcePathStr)
	if err != nil {
		return responseDto, errors.New("ReadSourcePathInfoError")
	}

	filesToFactory := []valueObject.UnixFilePath{requestDto.SourcePath}

	if sourcePathInfo.IsDir() {
		filesToFactoryWithoutSourcePath := filesToFactory[1:]
		filesToFactory = filesToFactoryWithoutSourcePath

		rawDirectoryFiles, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
			Command: "find",
			Args:    []string{"-L", sourcePathStr, "-maxdepth", "1", "-printf", "%p\n"},
		})
		if err != nil {
			return responseDto, errors.New("ReadDirectoryError: " + err.Error())
		}

		if len(rawDirectoryFiles) == 0 {
			return responseDto, errors.New("ReadDirectoryError")
		}

		rawDirectoryFilesList := strings.Split(rawDirectoryFiles, "\n")
		for _, fileToFactoryStr := range rawDirectoryFilesList {
			filePath, err := valueObject.NewUnixFilePath(fileToFactoryStr)
			if err != nil {
				slog.Error(
					"FileToFactoryError", slog.String("filePath", filePath.String()),
					slog.String("err", err.Error()),
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

	fileEntities := []entity.UnixFile{}
	for _, filePath := range filesToFactory {
		isFileTheSourcePath := filePath.String() == sourcePathStr
		if isFileTheSourcePath && sourcePathInfo.IsDir() {
			continue
		}

		fileEntity, err := repo.unixFileFactory(filePath, shouldReturnContent)
		if err != nil {
			slog.Error(
				"UnixFileFactoryError", slog.String("filePath", filePath.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		fileEntities = append(fileEntities, fileEntity)
	}

	responseDto = dto.ReadFilesResponse{Files: fileEntities}
	if requestDto.ShouldIncludeFileTree != nil && *requestDto.ShouldIncludeFileTree {
		responseDto.FileTree, err = repo.readUnixFileTree(
			requestDto.SourcePath, valueObject.FileSystemRootDir,
		)
		if err != nil {
			return responseDto, err
		}
	}

	return responseDto, nil
}

func (repo FilesQueryRepo) ReadFirst(
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
