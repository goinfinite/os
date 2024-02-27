package useCase

import (
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

const trashDirPath string = "/app/.trash"

func CreateTrash(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
) error {
	trashPath, _ := valueObject.NewUnixFilePath(trashDirPath)

	_, err := filesQueryRepo.GetOne(trashPath)
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

	return filesCmdRepo.Create(createTrashDir)
}

func DeleteUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
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

	if deleteUnixFiles.HardDelete {
		for _, fileToDelete := range deleteUnixFiles.SourcePaths {
			err := filesCmdRepo.Delete(fileToDelete)
			if err != nil {
				log.Printf("DeleteFileError: %s", err.Error())
				continue
			}

			log.Printf("File '%s' deleted.", fileToDelete.String())
		}

		return nil
	}

	err := CreateTrash(filesQueryRepo, filesCmdRepo)
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
		err = filesCmdRepo.Move(updateUnixFile, shouldOverwrite)
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
