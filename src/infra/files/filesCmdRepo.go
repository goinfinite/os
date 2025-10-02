package filesInfra

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/user"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type FilesCmdRepo struct {
	filesQueryRepo *FilesQueryRepo
}

func NewFilesCmdRepo() *FilesCmdRepo {
	return &FilesCmdRepo{
		filesQueryRepo: &FilesQueryRepo{},
	}
}

func (repo FilesCmdRepo) uploadFailureFactory(
	errMessage string,
	fileStreamHandler valueObject.FileStreamHandler,
) (uploadProcessFailure valueObject.UploadProcessFailure, err error) {
	failureReason, err := valueObject.NewFailureReason(errMessage)
	if err != nil {
		return uploadProcessFailure, err
	}

	return valueObject.NewUploadProcessFailure(
		fileStreamHandler.Name, failureReason,
	), nil
}

func (repo FilesCmdRepo) uploadSingleFile(
	destinationPath valueObject.UnixFilePath,
	fileToUpload valueObject.FileStreamHandler,
) error {
	destinationFilePath := destinationPath.String() + "/" + fileToUpload.Name.String()
	destinationEmptyFile, err := os.Create(destinationFilePath)
	if err != nil {
		return errors.New("CreateEmptyFileError: " + err.Error())
	}
	defer destinationEmptyFile.Close()

	fileToUploadStream, err := fileToUpload.Open()
	if err != nil {
		return errors.New("UnableToOpenFileStream: " + err.Error())
	}

	_, err = io.Copy(destinationEmptyFile, fileToUploadStream)
	if err != nil {
		return errors.New("CopyFileContentToDestinationError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) Copy(copyDto dto.CopyUnixFile) error {
	sourcePathStr := copyDto.SourcePath.String()
	fileToCopyExists := infraHelper.FileExists(sourcePathStr)
	if !fileToCopyExists {
		return errors.New("FileToCopyNotFound")
	}

	sourceFileName := copyDto.SourcePath.ReadFileName()
	destinationAbsolutePath := copyDto.DestinationPath.String() + "/" + sourceFileName.String()
	if !copyDto.ShouldOverwrite {
		destinationPathExists := infraHelper.FileExists(destinationAbsolutePath)
		if destinationPathExists {
			return errors.New("DestinationPathAlreadyExists")
		}
	}

	copyCmd := "rsync -avq " + sourcePathStr + " " + destinationAbsolutePath
	_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               copyCmd,
		ShouldRunWithSubShell: true,
	})
	return err
}

func (repo FilesCmdRepo) Compress(
	compressDto dto.CompressUnixFiles,
) (compressionProcessReport dto.CompressionProcessReport, err error) {
	compressibleFilesStr := []string{}
	incompressibleFilesStr := map[string]interface{}{}
	for _, sourcePath := range compressDto.SourcePaths {
		sourcePathExists := infraHelper.FileExists(sourcePath.String())
		if !sourcePathExists {
			incompressibleFilesStr[sourcePath.String()] = nil
			slog.Debug(
				"SourcePathNotFound", slog.String("sourcePath", sourcePath.String()),
			)
			continue
		}

		compressibleFilesStr = append(compressibleFilesStr, sourcePath.String())
	}

	if len(compressibleFilesStr) == 0 {
		return compressionProcessReport, errors.New("NoCompressibleFilesFound")
	}

	compressionTypeStr := "zip"

	destinationPathExt, err := compressDto.DestinationPath.ReadFileExtension()
	if err == nil {
		destinationPathExtStr := destinationPathExt.String()
		if destinationPathExtStr != "zip" {
			compressionTypeStr = "tgz"
		}
	}

	if compressDto.CompressionType != nil {
		compressionTypeStr = compressDto.CompressionType.String()
	}

	destinationPathWithoutExt := compressDto.DestinationPath.ReadWithoutExtension()
	compressionTypeAsExt := compressionTypeStr
	newDestinationPath, err := valueObject.NewUnixFilePath(
		destinationPathWithoutExt.String() + "." + compressionTypeAsExt,
	)
	if err != nil {
		return compressionProcessReport, errors.New(
			"CannotUpdateDestinationPathWithNewExtension",
		)
	}
	compressDto.DestinationPath = newDestinationPath

	_, err = valueObject.NewUnixCompressionType(compressionTypeStr)
	if err != nil {
		return compressionProcessReport, errors.New("UnsupportedCompressionType")
	}

	if infraHelper.FileExists(newDestinationPath.String()) {
		return compressionProcessReport, errors.New("DestinationPathAlreadyExists")
	}

	compressionBinary := "zip"
	compressionBinaryFlag := "-qr"
	if compressionTypeStr != "zip" {
		compressionBinary = "tar"
		compressionBinaryFlag = "-czf"
	}

	filesToCompress := strings.Join(compressibleFilesStr, " ")
	compressCmd := fmt.Sprintf(
		"%s %s %s %s",
		compressionBinary, compressionBinaryFlag,
		newDestinationPath.String(), filesToCompress,
	)
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               compressCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return compressionProcessReport, err
	}

	compressionProcessReport = dto.NewCompressionProcessReport(
		[]valueObject.UnixFilePath{},
		[]valueObject.CompressionProcessFailure{},
		newDestinationPath,
	)
	for _, sourcePath := range compressDto.SourcePaths {
		if _, isIncompressible := incompressibleFilesStr[sourcePath.String()]; isIncompressible {
			compressionProcessReport.FailedPathsWithReason = append(
				compressionProcessReport.FailedPathsWithReason,
				valueObject.NewCompressionProcessFailure(
					sourcePath, "SourcePathNotFound",
				),
			)

			continue
		}

		compressionProcessReport.FilePathsSuccessfullyCompressed = append(
			compressionProcessReport.FilePathsSuccessfullyCompressed, sourcePath,
		)
	}

	return compressionProcessReport, nil
}

func (repo FilesCmdRepo) Create(createDto dto.CreateUnixFile) error {
	filePathStr := createDto.FilePath.String()

	filesExists := infraHelper.FileExists(filePathStr)
	if filesExists {
		return errors.New("PathAlreadyExists")
	}

	unixUser, err := user.LookupId(createDto.OperatorAccountId.String())
	if err != nil {
		return errors.New("AccountNotFound")
	}

	fileOwnershipStr := unixUser.Username + ":" + unixUser.Username
	fileOwner, err := valueObject.NewUnixFileOwnership(fileOwnershipStr)
	if err != nil {
		return err
	}

	updateFileOwnerDto := dto.NewUpdateUnixFileOwnership(createDto.FilePath, fileOwner)

	if createDto.MimeType.IsDir() {
		err := os.MkdirAll(filePathStr, createDto.Permissions.GetFileMode())
		if err != nil {
			return err
		}

		return repo.UpdateOwnership(updateFileOwnerDto)
	}

	_, err = os.Create(filePathStr)
	if err != nil {
		return err
	}

	err = repo.UpdateOwnership(updateFileOwnerDto)
	if err != nil {
		return err
	}

	updatePermissionsDto := dto.NewUpdateUnixFilePermissions(
		createDto.FilePath, createDto.Permissions, nil,
	)
	return repo.UpdatePermissions(updatePermissionsDto)
}

func (repo FilesCmdRepo) Delete(unixFilePath valueObject.UnixFilePath) error {
	fileExists := infraHelper.FileExists(unixFilePath.String())
	if !fileExists {
		return errors.New("FileNotFound")
	}

	err := os.RemoveAll(unixFilePath.String())
	if err != nil {
		return errors.New("DeleteFileError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) Extract(extractDto dto.ExtractUnixFiles) error {
	fileToExtract := extractDto.SourcePath

	fileToExtractExists := infraHelper.FileExists(fileToExtract.String())
	if !fileToExtractExists {
		return errors.New("FileNotFound")
	}

	destinationPath := extractDto.DestinationPath

	destinationPathExists := infraHelper.FileExists(destinationPath.String())
	if destinationPathExists {
		return errors.New("DestinationPathAlreadyExists")
	}

	compressBinary := "tar"
	compressBinaryFlag := "-xf"
	compressDestinationFlag := "-C"

	unixFilePathExtension, err := fileToExtract.ReadFileExtension()
	if err != nil {
		return err
	}

	if unixFilePathExtension.String() == "zip" {
		compressBinary = "unzip"
		compressBinaryFlag = "-qq"
		compressDestinationFlag = "-d"
	}

	err = infraHelper.MakeDir(destinationPath.String())
	if err != nil {
		return err
	}

	compressCmd := fmt.Sprintf(
		"%s %s %s %s %s",
		compressBinary, compressBinaryFlag, fileToExtract.String(),
		compressDestinationFlag, destinationPath.String(),
	)
	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               compressCmd,
		ShouldRunWithSubShell: true,
	})
	return err
}

func (repo FilesCmdRepo) Move(moveDto dto.MoveUnixFile) error {
	sourcePathStr := moveDto.SourcePath.String()
	if !infraHelper.FileExists(sourcePathStr) {
		return errors.New("SourceFileNotFound")
	}

	if moveDto.DestinationPath == valueObject.UnixFilePathTrashDir {
		fileNameStr := moveDto.SourcePath.ReadFileName().String()
		destinationPathStr := moveDto.DestinationPath.String()
		rawTrashFilePath := destinationPathStr + "/" + fileNameStr
		trashFilePath, err := valueObject.NewUnixFilePath(rawTrashFilePath)
		if err != nil {
			return errors.New("MoveFileToTrashError: " + err.Error())
		}

		trashFilePathStr := trashFilePath.String()
		if infraHelper.FileExists(trashFilePathStr) {
			err := repo.Delete(trashFilePath)
			if err != nil {
				return errors.New("RemoveTrashFileWithSameNameError: " + err.Error())
			}

			slog.Debug("TrashFileWithSameNameRemoved", slog.String("filePath", trashFilePathStr))
		}

		return os.Rename(sourcePathStr, trashFilePathStr)
	}

	destinationPathStr := moveDto.DestinationPath.String()
	if infraHelper.FileExists(destinationPathStr) {
		if !moveDto.ShouldOverwrite {
			return errors.New("DestinationPathAlreadyExists")
		}

		err := repo.Delete(moveDto.DestinationPath)
		if err != nil {
			return errors.New("TrashFileError: " + err.Error())
		}
	}

	return os.Rename(sourcePathStr, destinationPathStr)
}

func (repo FilesCmdRepo) UpdateContent(
	updateContentDto dto.UpdateUnixFileContent,
) error {
	fileToUpdate, err := repo.filesQueryRepo.ReadFirst(updateContentDto.SourcePath)
	if err != nil {
		return err
	}

	if fileToUpdate.MimeType.IsDir() {
		return errors.New("PathIsADir")
	}

	decodedContent, err := updateContentDto.Content.GetDecodedContent()
	if err != nil {
		return err
	}

	return infraHelper.UpdateFile(
		updateContentDto.SourcePath.String(), decodedContent, true,
	)
}

func (repo FilesCmdRepo) UpdateOwnership(
	updateOwnershipDto dto.UpdateUnixFileOwnership,
) error {
	sourcePathStr := updateOwnershipDto.SourcePath.String()
	if !infraHelper.FileExists(sourcePathStr) {
		return errors.New("FileNotFound")
	}

	_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command: "chown",
		Args:    []string{updateOwnershipDto.Ownership.String(), sourcePathStr},
	})
	if err != nil {
		return errors.New("UpdateFileOwnershipError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) UpdatePermissions(
	updatePermissionsDto dto.UpdateUnixFilePermissions,
) error {
	sourcePathStr := updatePermissionsDto.SourcePath.String()
	if !infraHelper.FileExists(sourcePathStr) {
		return errors.New("FileOrDirNotFound")
	}

	updatePermissionsCmd := "find " + sourcePathStr + " -exec chmod " +
		updatePermissionsDto.FilePermissions.String() + " {} \\;"

	if updatePermissionsDto.DirectoryPermissions != nil {
		updatePermissionsCmd = fmt.Sprintf(
			"find %s -type d -exec chmod %s {} \\; && find %s -type f -exec chmod %s {} \\;",
			sourcePathStr, updatePermissionsDto.DirectoryPermissions.String(), sourcePathStr,
			updatePermissionsDto.FilePermissions.String(),
		)
	}

	_, err := infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               updatePermissionsCmd,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return errors.New("UpdatePermissionsError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) Upload(
	uploadDto dto.UploadUnixFiles,
) (dto.UploadProcessReport, error) {
	uploadProcessReport := dto.NewUploadProcessReport(
		[]valueObject.UnixFileName{},
		[]valueObject.UploadProcessFailure{},
		uploadDto.DestinationPath,
	)

	destinationFile, err := repo.filesQueryRepo.ReadFirst(uploadDto.DestinationPath)
	if err != nil {
		return uploadProcessReport, errors.New("DestinationFileNotFound")
	}

	if !destinationFile.MimeType.IsDir() {
		return uploadProcessReport, errors.New("DestinationPathCannotBeAFile")
	}

	for _, fileToUpload := range uploadDto.FileStreamHandlers {
		err := repo.uploadSingleFile(uploadDto.DestinationPath, fileToUpload)

		if err != nil {
			uploadFailure, err := repo.uploadFailureFactory(err.Error(), fileToUpload)
			if err != nil {
				slog.Debug("ReportUploadFailureError", slog.String("err", err.Error()))
			}

			uploadProcessReport.FailedNamesWithReason = append(
				uploadProcessReport.FailedNamesWithReason,
				uploadFailure,
			)

			continue
		}

		uploadProcessReport.FileNamesSuccessfullyUploaded = append(
			uploadProcessReport.FileNamesSuccessfullyUploaded,
			fileToUpload.Name,
		)

	}

	return uploadProcessReport, nil
}
