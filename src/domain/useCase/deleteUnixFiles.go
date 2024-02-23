package useCase

import (
	"log"
	"slices"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

func DeleteUnixFiles(
	filesQueryRepo repository.FilesQueryRepo,
	filesCmdRepo repository.FilesCmdRepo,
	deleteUnixFiles dto.DeleteUnixFiles,
) {
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
		return
	}

	trashPath, _ := valueObject.NewUnixFilePath("/app/.trash")
	_, err := filesQueryRepo.GetOne(trashPath)
	if err != nil {
		trashDirPermissions, _ := valueObject.NewUnixFilePermissions("775")
		trashDirMimeType, _ := valueObject.NewMimeType("directory")
		createTrashDir := dto.NewCreateUnixFile(
			trashPath,
			trashDirPermissions,
			trashDirMimeType,
		)

		filesCmdRepo.Create(createTrashDir)
	}

	for _, fileToMoveToTrash := range deleteUnixFiles.SourcePaths {
		trashPathWithFileNameStr := trashPath.String() + "/" + fileToMoveToTrash.GetFileName().String()
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
}
