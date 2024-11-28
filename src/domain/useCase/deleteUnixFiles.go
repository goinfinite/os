package useCase

import (
	"log/slog"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/repository"
	"github.com/goinfinite/os/src/domain/valueObject"
)

const TrashDirPath string = "/app/.trash"

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
	trashPath, _ := valueObject.NewUnixFilePath(TrashDirPath)
	err := uc.filesCmdRepo.Delete(trashPath)
	if err != nil {
		return err
	}

	return uc.CreateGeneralTrash(operatorAccountId, operatorIpAddress)
}

func (uc DeleteUnixFiles) CreateGeneralTrash(
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) error {
	trashPath, _ := valueObject.NewUnixFilePath(TrashDirPath)

	_, err := uc.filesQueryRepo.ReadFirst(trashPath)
	if err == nil {
		return nil
	}

	trashDirPermissions, _ := valueObject.NewUnixFilePermissions("755")
	trashDirMimeType, _ := valueObject.NewMimeType("directory")
	createGeneralTrashDir := dto.NewCreateUnixFile(
		trashPath, &trashDirPermissions, trashDirMimeType, operatorAccountId,
		operatorIpAddress,
	)

	return CreateUnixFile(
		uc.filesQueryRepo, uc.filesCmdRepo, uc.activityRecordCmdRepo, createGeneralTrashDir,
	)
}

func (uc DeleteUnixFiles) Execute(deleteDto dto.DeleteUnixFiles) error {
	for fileToDeleteIndex, fileToDelete := range deleteDto.SourcePaths {
		shouldCleanTrash := fileToDelete.String() == TrashDirPath
		if shouldCleanTrash {
			err := uc.emptyTrash(
				deleteDto.OperatorAccountId, deleteDto.OperatorIpAddress,
			)
			if err != nil {
				slog.Debug("FailedToCleanTrash", slog.Any("err", err))
			}

			fileToDeleteAfterTrashPathIndex := fileToDeleteIndex + 1
			filesToDeleteWithoutTrashPath := slices.Delete(
				deleteDto.SourcePaths, fileToDeleteIndex,
				fileToDeleteAfterTrashPathIndex,
			)

			deleteDto.SourcePaths = filesToDeleteWithoutTrashPath

			continue
		}

		isRootPath := fileToDelete.String() == "/"
		if !isRootPath {
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
				slog.Debug("DeleteFileError", slog.Any("err", err))
				continue
			}
		}

		return nil
	}

	err := uc.CreateGeneralTrash(deleteDto.OperatorAccountId, deleteDto.OperatorIpAddress)
	if err != nil {
		return err
	}

	for _, fileToMoveToTrash := range deleteDto.SourcePaths {
		trashPathWithFileNameStr := TrashDirPath + "/" + fileToMoveToTrash.GetFileName().String()
		trashPathWithFileName, _ := valueObject.NewUnixFilePath(trashPathWithFileNameStr)
		shouldOverwrite := true
		moveDto := dto.NewMoveUnixFile(
			fileToMoveToTrash, trashPathWithFileName, shouldOverwrite,
		)

		err = uc.filesCmdRepo.Move(moveDto)
		if err != nil {
			slog.Debug(
				"MoveUnixFileToTrashError",
				slog.String("fileToMoveToTrash", fileToMoveToTrash.String()),
				slog.Any("err", err),
			)
			continue
		}

		slog.Info(
			"FileMovedToTrash", slog.String("filePath", fileToMoveToTrash.String()),
		)
	}

	NewCreateSecurityActivityRecord(uc.activityRecordCmdRepo).
		DeleteUnixFiles(deleteDto)

	return nil
}
