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
	filesQueryRepo repository.FilesQueryRepo
	filesCmdRepo   repository.FilesCmdRepo
}

func NewDeleteUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
) DeleteUnixFiles {
	return DeleteUnixFiles{
		filesQueryRepo: filesQueryRepo,
		filesCmdRepo:   filesCmdRepo,
	}
}

func (uc DeleteUnixFiles) emptyTrash() error {
	trashPath, _ := valueObject.NewUnixFilePath(TrashDirPath)
	err := uc.filesCmdRepo.Delete(trashPath)
	if err != nil {
		return err
	}

	return uc.CreateTrash()
}

func (uc DeleteUnixFiles) CreateTrash() error {
	trashPath, _ := valueObject.NewUnixFilePath(TrashDirPath)

	_, err := uc.filesQueryRepo.ReadFirst(trashPath)
	if err == nil {
		return nil
	}

	trashDirPermissions, _ := valueObject.NewUnixFilePermissions("755")
	trashDirMimeType, _ := valueObject.NewMimeType("directory")
	createTrashDir := dto.NewCreateUnixFile(
		trashPath, &trashDirPermissions, trashDirMimeType,
	)

	return uc.filesCmdRepo.Create(createTrashDir)
}

func (uc DeleteUnixFiles) Execute(deleteDto dto.DeleteUnixFiles) error {
	for fileToDeleteIndex, fileToDelete := range deleteDto.SourcePaths {
		shouldCleanTrash := fileToDelete.String() == TrashDirPath
		if shouldCleanTrash {
			err := uc.emptyTrash()
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

	err := uc.CreateTrash()
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

	return nil
}
