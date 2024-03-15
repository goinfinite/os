package useCase

import (
	"errors"
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/repository"
	"github.com/speedianet/os/src/domain/valueObject"
)

type UpdateUnixFiles struct {
	filesCmdRepo repository.FilesCmdRepo
}

func NewUpdateUnixFiles(
	filesCmdRepo repository.FilesCmdRepo,
) UpdateUnixFiles {
	return UpdateUnixFiles{
		filesCmdRepo: filesCmdRepo,
	}
}

func (uc UpdateUnixFiles) updateFailureFactory(
	filePath valueObject.UnixFilePath,
	errMessage string,
) valueObject.UpdateProcessFailure {
	return valueObject.NewUpdateProcessFailure(
		filePath,
		errMessage,
	)
}

func (uc UpdateUnixFiles) updateUnixFilePermissions(
	sourcePath valueObject.UnixFilePath,
	permissions valueObject.UnixFilePermissions,
) error {
	err := uc.filesCmdRepo.UpdatePermissions(
		sourcePath,
		permissions,
	)
	if err != nil {
		log.Printf("UpdateUnixFilesPermissionsError: %s", err.Error())
		return errors.New("UpdateUnixFilesPermissionsInfraError")
	}

	log.Printf(
		"File '%s' (%s) permissions updated to '%s'.",
		sourcePath.GetFileName().String(),
		sourcePath.GetFileDir().String(),
		permissions.String(),
	)

	return nil
}

func (uc UpdateUnixFiles) updateUnixFilePath(
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
		log.Printf("MoveUnixFileError: %s", err.Error())
		return errors.New("MoveUnixFileInfraError")
	}

	log.Printf(
		"File '%s' moved from %s to '%s'.",
		sourcePath.GetFileName().String(),
		sourcePath.GetFileDir().String(),
		destinationPath.GetFileDir().String(),
	)

	return nil
}

func (uc UpdateUnixFiles) updateUnixFileContent(
	sourcePath valueObject.UnixFilePath,
	encodedContent valueObject.EncodedContent,
) error {
	err := uc.filesCmdRepo.UpdateContent(sourcePath, encodedContent)
	if err != nil {
		log.Printf("UpdateUnixFilesContentError: %s", err.Error())
		return errors.New("UpdateUnixFilesContentInfraError")
	}

	log.Printf(
		"File '%s' content updated.",
		sourcePath.GetFileName().String(),
	)

	return nil
}

func (uc UpdateUnixFiles) Execute(
	updateUnixFiles dto.UpdateUnixFiles,
) (dto.UpdateProcessReport, error) {
	updateProcessReport := dto.NewUpdateProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.UpdateProcessFailure{},
	)

	for _, sourcePath := range updateUnixFiles.SourcePaths {
		if updateUnixFiles.Permissions != nil {
			err := uc.updateUnixFilePermissions(
				sourcePath,
				*updateUnixFiles.Permissions,
			)
			if err != nil {
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					uc.updateFailureFactory(sourcePath, err.Error()),
				)
				continue
			}
		}

		if updateUnixFiles.DestinationPath != nil {
			err := uc.updateUnixFilePath(
				sourcePath,
				*updateUnixFiles.DestinationPath,
			)
			if err != nil {
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					uc.updateFailureFactory(sourcePath, err.Error()),
				)
				continue
			}
		}

		if updateUnixFiles.EncodedContent != nil {
			err := uc.updateUnixFileContent(
				sourcePath,
				*updateUnixFiles.EncodedContent,
			)
			if err != nil {
				updateProcessReport.FailedPathsWithReason = append(
					updateProcessReport.FailedPathsWithReason,
					uc.updateFailureFactory(sourcePath, err.Error()),
				)
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
