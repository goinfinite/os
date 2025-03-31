package filesInfra

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type FilesQueryRepo struct{}

func (repo *FilesQueryRepo) unixFileFactory(
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
		unixFileContentStr, err := infraHelper.ReadFileContent(filePath.String())
		if err != nil {
			return unixFile, errors.New("FailedToReadFileContent: " + err.Error())
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

func (repo *FilesQueryRepo) simplifiedUnixFileFactory(
	unixFilePath valueObject.UnixFilePath,
) (simplifiedUnixFile entity.SimplifiedUnixFile, err error) {
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

func (repo *FilesQueryRepo) unixFileBranchFactory(
	branchAbsolutePath valueObject.UnixFilePath,
	shouldIncludeFiles bool,
) (fileBranch dto.UnixFileBranch, err error) {
	simplifiedBranchFileEntity, err := repo.simplifiedUnixFileFactory(branchAbsolutePath)
	if err != nil {
		return fileBranch, err
	}
	fileBranch = dto.NewUnixFileBranch(simplifiedBranchFileEntity)

	findCmdArgs := []string{
		"-L", branchAbsolutePath.String(),
		"-mindepth", "1",
		"-maxdepth", "1",
	}
	if !shouldIncludeFiles {
		findCmdArgs = append(findCmdArgs, "-type", "d")
	}

	rawBranchTwigs, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "find",
		Args:    findCmdArgs,
	})
	if err != nil {
		return fileBranch, err
	}

	if rawBranchTwigs == "" {
		return fileBranch, nil
	}

	rawFactorableFiles := strings.SplitSeq(rawBranchTwigs, "\n")
	for rawFilePath := range rawFactorableFiles {
		if rawFilePath == "" {
			continue
		}

		twigPath, err := valueObject.NewUnixFilePath(rawFilePath)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("rawTwigPath", rawFilePath),
			)
			continue
		}

		simplifiedTwigFileEntity, err := repo.simplifiedUnixFileFactory(twigPath)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("twigPath", twigPath.String()),
			)
			continue
		}
		fileBranch.Branches[simplifiedTwigFileEntity.Name] = dto.NewUnixFileBranch(simplifiedTwigFileEntity)
	}

	return fileBranch, nil
}

func (repo *FilesQueryRepo) unixFileTreeFactory(
	leafAbsolutePath valueObject.UnixFilePath,
) (treeTrunk dto.UnixFileBranch, err error) {
	rawTreeBranches := strings.SplitSeq(leafAbsolutePath.String(), "/")

	shouldIncludeFiles := false
	iterationBranch := treeTrunk
	iterationBranchPath := ""
	for rawBranchName := range rawTreeBranches {
		rawBranchName = strings.TrimSpace(rawBranchName)
		isTreeTrunk := rawBranchName == ""

		iterationBranchPath += rawBranchName + "/"
		branchFilePath, err := valueObject.NewUnixFilePath(iterationBranchPath)
		if err != nil {
			slog.Debug(
				err.Error(),
				slog.String("rawBranchPath", iterationBranchPath),
			)
			continue
		}

		treeBranch, err := repo.unixFileBranchFactory(branchFilePath, shouldIncludeFiles)
		if err != nil {
			slog.Error(
				err.Error(),
				slog.String("branchFilePath", branchFilePath.String()),
			)
			continue
		}
		if isTreeTrunk {
			treeTrunk = treeBranch
			iterationBranch = treeTrunk
			continue
		}

		iterationBranch.Branches[treeBranch.Name] = treeBranch
		iterationBranch = treeBranch
	}

	return treeTrunk, nil
}

func (repo *FilesQueryRepo) Read(
	requestDto dto.ReadFilesRequest,
) (responseDto dto.ReadFilesResponse, err error) {
	sourcePathStr := requestDto.SourcePath.String()

	if !infraHelper.FileExists(sourcePathStr) {
		return responseDto, errors.New("PathNotFound")
	}

	sourcePathInfo, err := os.Stat(sourcePathStr)
	if err != nil {
		return responseDto, errors.New("ReadSourcePathInfoError: " + err.Error())
	}

	factorableFilePaths := []valueObject.UnixFilePath{requestDto.SourcePath}

	if sourcePathInfo.IsDir() {
		factorableFilePathsWithoutSourcePath := factorableFilePaths[1:]
		factorableFilePaths = factorableFilePathsWithoutSourcePath

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

		rawDirectoryFilesList := strings.SplitSeq(rawDirectoryFiles, "\n")
		for fileToFactoryStr := range rawDirectoryFilesList {
			filePath, err := valueObject.NewUnixFilePath(fileToFactoryStr)
			if err != nil {
				slog.Error(
					"FactoryFileError",
					slog.String("filePath", filePath.String()),
					slog.String("err", err.Error()),
				)
				continue
			}

			factorableFilePaths = append(factorableFilePaths, filePath)
		}
	}

	shouldReturnContent := false
	if len(factorableFilePaths) == 1 {
		shouldReturnContent = true
	}

	fileEntities := []entity.UnixFile{}
	directoryEntities := []entity.UnixFile{}

	isSourcePathDir := sourcePathInfo.IsDir()
	for _, filePath := range factorableFilePaths {
		isFileTheSourcePath := filePath.String() == sourcePathStr
		if isFileTheSourcePath && isSourcePathDir {
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

		if fileEntity.MimeType.IsDir() {
			directoryEntities = append(directoryEntities, fileEntity)
			continue
		}

		fileEntities = append(fileEntities, fileEntity)
	}

	slices.SortStableFunc(fileEntities, func(a, b entity.UnixFile) int {
		return strings.Compare(a.Name.String(), b.Name.String())
	})
	slices.SortStableFunc(directoryEntities, func(a, b entity.UnixFile) int {
		return strings.Compare(a.Name.String(), b.Name.String())
	})
	fileEntities = append(directoryEntities, fileEntities...)

	responseDto = dto.ReadFilesResponse{Files: fileEntities}
	if requestDto.ShouldIncludeFileTree != nil && *requestDto.ShouldIncludeFileTree {
		filesTree, err := repo.unixFileTreeFactory(requestDto.SourcePath)
		if err != nil {
			return responseDto, err
		}

		responseDto.FileTree = &filesTree
	}

	return responseDto, nil
}

func (repo *FilesQueryRepo) ReadFirst(
	unixFilePath valueObject.UnixFilePath,
) (unixFileEntity entity.UnixFile, err error) {
	if !infraHelper.FileExists(unixFilePath.String()) {
		return unixFileEntity, errors.New("FileNotFound")
	}

	shouldReturnContent := false
	return repo.unixFileFactory(unixFilePath, shouldReturnContent)
}
