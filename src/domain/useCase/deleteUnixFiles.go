package useCase

import (
	"log/slog"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

type DeleteUnixFiles struct {
	filesQueryRepo        repository.FilesQueryRepo
	filesCmdRepo          repository.FilesCmdRepo
	activityRecordCmdRepo repository.ActivityRecordCmdRepo
}

func NewDeleteUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	activityRecordCmdRepo repository.ActivityRecordCmdRepo,
) DeleteUnixFiles {
	return DeleteUnixFiles{
		filesQueryRepo:        filesQueryRepo,
		filesCmdRepo:          filesCmdRepo,
		activityRecordCmdRepo: activityRecordCmdRepo,
	}
}

func (uc DeleteUnixFiles) emptyTrash(
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) error {
	err := uc.filesCmdRepo.Delete(valueObject.UnixFilePathTrashDir)
	if err != nil {
		return err
	}

	return uc.CreateGeneralTrash(operatorAccountId, operatorIpAddress)
}

func (uc DeleteUnixFiles) CreateGeneralTrash(
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) error {
	_, err := uc.filesQueryRepo.ReadFirst(valueObject.UnixFilePathTrashDir)
	if err == nil {
		return nil
	}

	trashDirPermissions := valueObject.NewUnixDirDefaultPermissions()
	createGeneralTrashDir := dto.NewCreateUnixFile(
		valueObject.UnixFilePathTrashDir, &trashDirPermissions,
		valueObject.MimeTypeDirectory, operatorAccountId, operatorIpAddress,
	)

	return CreateUnixFile(
		uc.filesQueryRepo, uc.filesCmdRepo, uc.activityRecordCmdRepo,
		createGeneralTrashDir,
	)
}

func (uc DeleteUnixFiles) Execute(deleteDto dto.DeleteUnixFiles) error {
	for fileToDeleteIndex, fileToDelete := range deleteDto.SourcePaths {
		shouldCleanTrash := fileToDelete == valueObject.UnixFilePathTrashDir
		if shouldCleanTrash {
			err := uc.emptyTrash(
				deleteDto.OperatorAccountId, deleteDto.OperatorIpAddress,
			)
			if err != nil {
				slog.Debug("FailedToCleanTrash", slog.String("err", err.Error()))
			}

			fileToDeleteAfterTrashPathIndex := fileToDeleteIndex + 1
			filesToDeleteWithoutTrashPath := slices.Delete(
				deleteDto.SourcePaths, fileToDeleteIndex,
				fileToDeleteAfterTrashPathIndex,
			)

			deleteDto.SourcePaths = filesToDeleteWithoutTrashPath

			continue
		}

		if !fileToDelete.IsFileSystemRootDir() {
			continue
		}

		slog.Debug(
			"DeleteUnixFilesError", slog.String("err", "Path '/' cannot be deleted."),
		)

		fileToDeleteAfterNotAllowedPathIndex := fileToDeleteIndex + 1
		filesToDeleteWithoutNotAllowedPath := slices.Delete(
			deleteDto.SourcePaths, fileToDeleteIndex,
			fileToDeleteAfterNotAllowedPathIndex,
		)

		deleteDto.SourcePaths = filesToDeleteWithoutNotAllowedPath
	}

	if deleteDto.HardDelete {
		for _, fileToDelete := range deleteDto.SourcePaths {
			err := uc.filesCmdRepo.Delete(fileToDelete)
			if err != nil {
				slog.Debug("DeleteFileError", slog.String("err", err.Error()))
				continue
			}
		}

		return nil
	}

	err := uc.CreateGeneralTrash(deleteDto.OperatorAccountId, deleteDto.OperatorIpAddress)
	if err != nil {
		return err
	}

	shouldOverwrite := true
	for _, disposableFile := range deleteDto.SourcePaths {
		err = uc.filesCmdRepo.Move(
			dto.NewMoveUnixFile(disposableFile, valueObject.UnixFilePathTrashDir, shouldOverwrite),
		)
		if err != nil {
			slog.Debug(
				"MoveUnixFileToTrashError",
				slog.String("fileToMoveToTrash", disposableFile.String()),
				slog.String("err", err.Error()),
			)
			continue
		}

		slog.Info(
			"FileMovedToTrash", slog.String("filePath", disposableFile.String()),
		)
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).DeleteUnixFiles(deleteDto)
	return nil
}
