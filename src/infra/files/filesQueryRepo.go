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
	unixFileExtension, err := unixFilePath.ReadFileExtension()
	if err == nil {
		unixFileExtensionPtr = &unixFileExtension
	}

	unixFileMimeType := unixFileExtension.GetMimeType()
	if fileInfo.IsDir() {
		unixFileMimeType = valueObject.MimeTypeDirectory
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
		unixFilePath.ReadFileName(), unixFilePath, unixFileMimeType, unixFilePermissions,
		unixFileSize, unixFileExtensionPtr, unixFileContentPtr, unixFileUid,
		unixFileUsername, unixFileGid, unixFileGroup, unixFileUpdatedAt,
	)

	return unixFile, nil
}

func (repo FilesQueryRepo) simplifiedUnixFileFactory(
	unixFilePath valueObject.UnixFilePath,
) (simplifiedUnixFile entity.SimplifiedUnixFile, err error) {
	unixFilePath = unixFilePath.ReadWithoutTrailingSlash()

	fileInfo, err := os.Stat(unixFilePath.String())
	if err != nil {
		return simplifiedUnixFile, err
	}

	unixFileMimeType := valueObject.MimeTypeGeneric
	if fileInfo.IsDir() {
		unixFileMimeType = valueObject.MimeTypeDirectory
	}

	unixFileExtension, err := unixFilePath.ReadFileExtension()
	if err == nil {
		unixFileMimeType = unixFileExtension.GetMimeType()
	}

	return entity.NewSimplifiedUnixFile(
		unixFilePath.ReadFileName(), unixFilePath, unixFileMimeType,
	), nil
}

func (repo FilesQueryRepo) readUnixFileBranch(
	desiredAbsolutePath valueObject.UnixFilePath,
	shouldIncludeFiles bool,
) (treeBranch dto.UnixFileBranch, err error) {
	treeDirectory, err := repo.simplifiedUnixFileFactory(desiredAbsolutePath)
	if err != nil {
		return treeBranch, err
	}
	treeBranch = dto.NewUnixFileBranch(treeDirectory)

	findCmdArgs := []string{
		"-L", desiredAbsolutePath.String(), "-mindepth", "1", "-maxdepth", "1",
	}
	if !shouldIncludeFiles {
		findCmdArgs = append(findCmdArgs, "-type", "d")
	}

	rawDirectoryTree, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find",
		Args:    findCmdArgs,
	})
	if err != nil {
		return treeBranch, err
	}

	directoriesToBrowse := strings.Split(rawDirectoryTree, "\n")
	if len(directoriesToBrowse) == 0 {
		return treeBranch, err
	}

	for _, rawDirectoryPath := range directoriesToBrowse {
		if rawDirectoryPath == "" {
			continue
		}

		directoryPath, err := valueObject.NewUnixFilePath(rawDirectoryPath)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("directoryPath", rawDirectoryPath),
			)
			continue
		}

		treeDirectory, err = repo.simplifiedUnixFileFactory(directoryPath)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("directoryPath", rawDirectoryPath),
			)
			continue
		}
		treeBranch.Branches[treeDirectory.Name] = dto.NewUnixFileBranch(treeDirectory)
	}

	return treeBranch, nil
}

func (repo FilesQueryRepo) readUnixFileTree(
	desiredAbsolutePath valueObject.UnixFilePath,
) (trunkBranch dto.UnixFileBranch, err error) {
	fileSystemRootDirectory, err := repo.simplifiedUnixFileFactory(
		valueObject.UnixFilePathFileSystemRootDir,
	)
	if err != nil {
		return trunkBranch, err
	}
	trunkBranch = dto.NewUnixFileBranch(fileSystemRootDirectory)

	directoriesToBrowse := strings.Split(
		desiredAbsolutePath.ReadWithoutTrailingSlash().String(), "/",
	)

	directoriesToBrowseLength := len(directoriesToBrowse)
	if directoriesToBrowseLength == 0 {
		return trunkBranch, err
	}

	currentDirectoryPath := ""
	currentBranch := trunkBranch
	for rawDirectoryIndex, rawDirectoryName := range directoriesToBrowse {
		currentDirectoryPath += rawDirectoryName + "/"

		isRootDirectory := rawDirectoryIndex == 0
		if isRootDirectory {
			currentDirectoryPath = "/"
		}

		directoryPath, err := valueObject.NewUnixFilePath(currentDirectoryPath)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("directoryPath", currentDirectoryPath),
			)
			continue
		}

		isTheLastDirectory := directoriesToBrowseLength == (rawDirectoryIndex + 1)
		treeBranch, err := repo.readUnixFileBranch(directoryPath, isTheLastDirectory)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("directoryPath", currentDirectoryPath),
			)
			continue
		}

		currentBranch.Branches[treeBranch.Name] = treeBranch
		currentBranch = currentBranch.Branches[treeBranch.Name]
	}

	return trunkBranch, nil
}

func (repo FilesQueryRepo) Read(
	requestDto dto.ReadFilesRequest,
) (responseDto dto.ReadFilesResponse, err error) {
	sourcePath := requestDto.SourcePath.ReadWithoutTrailingSlash()
	sourcePathStr := sourcePath.String()

	exists := infraHelper.FileExists(sourcePathStr)
	if !exists {
		return responseDto, errors.New("PathNotFound")
	}

	sourcePathInfo, err := os.Stat(sourcePathStr)
	if err != nil {
		return responseDto, errors.New("ReadSourcePathInfoError")
	}

	filesToFactory := []valueObject.UnixFilePath{sourcePath}

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
		unixFileTree, err := repo.readUnixFileTree(sourcePath)
		if err != nil {
			return responseDto, err
		}

		responseDto.FileTree = &unixFileTree
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
