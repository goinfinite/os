package useCase

import (
	"errors"
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

const trashDirPath string = "/app/.trash"

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
	trashPath, _ := valueObject.NewUnixFilePath(trashDirPath)
	err := uc.filesCmdRepo.Delete(trashPath)
	if err != nil {
		return err
	}

	return uc.CreateTrash()
}

func (uc DeleteUnixFiles) CreateTrash() error {
	trashPath, _ := valueObject.NewUnixFilePath(trashDirPath)

	_, err := uc.filesQueryRepo.GetOne(trashPath)
	if err == nil {
		return nil
	}

	trashDirPermissions, _ := valueObject.NewUnixFilePermissions("755")
	trashDirMimeType, _ := valueObject.NewMimeType("directory")
	createTrashDir := dto.NewCreateUnixFile(
		trashPath,
		&trashDirPermissions,
		trashDirMimeType,
	)

	return uc.filesCmdRepo.Create(createTrashDir)
}

func (uc DeleteUnixFiles) Execute(
	deleteUnixFiles dto.DeleteUnixFiles,
) error {
	for fileToDeleteIndex, fileToDelete := range deleteUnixFiles.SourcePaths {
		isRootPath := fileToDelete.String() == "/"
		if !isRootPath {
			continue
		}

		log.Printf("DeleteUnixFilesError: Path '/' cannot be deleted.")

		fileToDeleteAfterNotAllowedPathIndex := fileToDeleteIndex + 1
		filesToDeleteWithoutNotAllowedPath := slices.Delete(
			deleteUnixFiles.SourcePaths,
			fileToDeleteIndex,
			fileToDeleteAfterNotAllowedPathIndex,
		)

		deleteUnixFiles.SourcePaths = filesToDeleteWithoutNotAllowedPath
	}

	if (len(deleteUnixFiles.SourcePaths)) == 0 {
		return errors.New("ThereAreNoSourcePathsToDelete")
	}

	if deleteUnixFiles.HardDelete {
		for _, fileToDelete := range deleteUnixFiles.SourcePaths {
			err := uc.filesCmdRepo.Delete(fileToDelete)
			if err != nil {
				log.Printf("DeleteFileError: %s", err.Error())
				continue
			}

			log.Printf("File '%s' deleted.", fileToDelete.String())
		}

		return nil
	}

	err := uc.CreateTrash()
	if err != nil {
		return err
	}

	for _, fileToMoveToTrash := range deleteUnixFiles.SourcePaths {
		trashPathWithFileNameStr := trashDirPath + "/" + fileToMoveToTrash.GetFileName().String()
		trashPathWithFileName, _ := valueObject.NewUnixFilePath(trashPathWithFileNameStr)

		updateUnixFile := dto.NewUpdateUnixFile(
			fileToMoveToTrash,
			&trashPathWithFileName,
			nil,
			nil,
		)

		shouldOverwrite := true
		err = uc.filesCmdRepo.Move(updateUnixFile, shouldOverwrite)
		if err != nil {
			log.Printf(
				"MoveUnixFileToTrashError (%s): %s",
				fileToMoveToTrash.String(),
				err.Error(),
			)
			continue
		}

		log.Printf("File '%s' moved to trash.", fileToMoveToTrash.String())
	}

	return nil
}
