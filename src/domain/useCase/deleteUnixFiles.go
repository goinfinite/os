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

	if deleteUnixFiles.PermanentDelete {
		filesCmdRepo.Delete(deleteUnixFiles)
		return
	}

	for _, fileToMoveToTrash := range deleteUnixFiles.SourcePaths {
		trashPath, _ := valueObject.NewUnixFilePath("/.trash")
		_, err := filesQueryRepo.GetOne(trashPath)
		if err != nil {
			log.Print("TrashNotFound")
		}

		trashPathWithFileNameStr := trashPath.String() + "/" + fileToMoveToTrash.GetFileName().String()
		trashPathWithFileName, _ := valueObject.NewUnixFilePath(trashPathWithFileNameStr)

		updateUnixFile := dto.NewUpdateUnixFile(
			fileToMoveToTrash,
			&trashPathWithFileName,
			nil,
			nil,
		)

		err = filesCmdRepo.Move(updateUnixFile)
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
