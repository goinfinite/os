package useCase

import (
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

type UpdateUnixFile struct {
	filesCmdRepo repository.FilesCmdRepo
}

func NewUpdateUnixFile(
	filesCmdRepo repository.FilesCmdRepo,
) UpdateUnixFile {
	return UpdateUnixFile{
		filesCmdRepo: filesCmdRepo,
	}
}

func (uc UpdateUnixFile) updateFailureFactory(
	filePath valueObject.UnixFilePath,
	errMessage string,
) valueObject.UpdateProcessFailure {
	return valueObject.NewUpdateProcessFailure(
		filePath,
		errMessage,
	)
}

func (uc UpdateUnixFile) updateUnixFilePermissions(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
) error {
	err := uc.filesCmdRepo.UpdatePermissions(
		sourcePath,
		permissions,
	)
	if err != nil {
		return err
	}

	log.Printf(
		"File '%s' (%s) permissions updated to '%s'.",
		sourcePath.GetFileName().String(),
		sourcePath.GetFileDir().String(),
		permissions.String(),
	)

	return nil
}

func (uc UpdateUnixFile) updateUnixFilePath(
	sourcePath valueObject.UnixFilePath,
	destinationPath valueObject.UnixFilePath,
) error {
	shouldOverwrite := false
	err := uc.filesCmdRepo.Move(
		sourcePath,
		destinationPath,
		shouldOverwrite,
	)
	if err != nil {
		return err
	}

	log.Printf(
		"File '%s' moved from %s to '%s'.",
		sourcePath.GetFileName().String(),
		sourcePath.GetFileDir().String(),
		destinationPath.GetFileDir().String(),
	)

	return nil
}

func (uc UpdateUnixFile) updateUnixFileContent(
	sourcePath valueObject.UnixFilePath,
	encodedContent valueObject.EncodedContent,
) error {
	err := uc.filesCmdRepo.UpdateContent(sourcePath, encodedContent)
	if err != nil {
		return err
	}

	log.Printf(
		"File '%s' content updated.",
		sourcePath.GetFileName().String(),
	)

	return nil
}

func (uc UpdateUnixFile) Execute(
	updateUnixFile dto.UpdateUnixFile,
) (dto.UpdateProcessReport, error) {
	updateProcessReport := dto.NewUpdateProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.UpdateProcessFailure{},
	)

	for _, sourcePath := range updateUnixFile.SourcePaths {
		if updateUnixFile.Permissions != nil {
			err := uc.updateUnixFilePermissions(
				sourcePath,
				*updateUnixFile.Permissions,
			)
			if err != nil {
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					uc.updateFailureFactory(sourcePath, "UpdateUnixFilePermissionsInfraError"),
				)
				log.Printf("UpdateUnixFilePermissionsError: %s", err.Error())
				continue
			}
		}

		if updateUnixFile.DestinationPath != nil {
			err := uc.updateUnixFilePath(
				sourcePath,
				*updateUnixFile.DestinationPath,
			)
			if err != nil {
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					uc.updateFailureFactory(sourcePath, "MoveUnixFileInfraError"),
				)
				log.Printf("MoveUnixFileError: %s", err.Error())
				continue
			}
		}

		if updateUnixFile.EncodedContent != nil {
			err := uc.filesCmdRepo.UpdateContent(
				sourcePath,
				*updateUnixFile.EncodedContent,
			)
			if err != nil {
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					uc.updateFailureFactory(sourcePath, "UpdateUnixFileContentInfraError"),
				)
				log.Printf("UpdateUnixFileContentError: %s", err.Error())
				continue
			}
		}

		updateProcessReport.FilePathsSuccessfullyUpdated = append(
			updateProcessReport.FilePathsSuccessfullyUpdated,
			sourcePath,
		)
	}

	return updateProcessReport, nil
}
