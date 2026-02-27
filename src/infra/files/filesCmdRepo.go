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
	tkInfra "github.com/goinfinite/tk/src/infra"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type FilesCmdRepo struct {
	filesQueryRepo *FilesQueryRepo
	fileClerk      tkInfra.FileClerk
}

func NewFilesCmdRepo() *FilesCmdRepo {
	return &FilesCmdRepo{
		filesQueryRepo: NewFilesQueryRepo(),
		fileClerk:      tkInfra.FileClerk{},
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
	destinationPath tkValueObject.UnixAbsoluteFilePath,
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
	fileToCopyExists := repo.fileClerk.FileExists(sourcePathStr)
	if !fileToCopyExists {
		return errors.New("FileToCopyNotFound")
	}

	sourceFileName := copyDto.SourcePath.ReadFileName(false)
	destinationAbsolutePath := copyDto.DestinationPath.String() + "/" + sourceFileName.String()
	if !copyDto.ShouldOverwrite {
		destinationPathExists := repo.fileClerk.FileExists(destinationAbsolutePath)
		if destinationPathExists {
			return errors.New("DestinationPathAlreadyExists")
		}
	}

	copyCmd := "rsync -avq " + sourcePathStr + " " + destinationAbsolutePath
	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:          copyCmd,
		ShouldUseSubShell: true,
	}).Run()
	return err
}

func (repo FilesCmdRepo) Compress(
	compressDto dto.CompressUnixFiles,
) (compressionProcessReport dto.CompressionProcessReport, err error) {
	compressibleFilesStr := []string{}
	incompressibleFilesStr := map[string]interface{}{}
	for _, sourcePath := range compressDto.SourcePaths {
		sourcePathExists := repo.fileClerk.FileExists(sourcePath.String())
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

	destinationPathWithoutExt := compressDto.DestinationPath.ReadWithoutExtension(false)
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

	if repo.fileClerk.FileExists(newDestinationPath.String()) {
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
	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command:          compressCmd,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return compressionProcessReport, err
	}

	compressionProcessReport = dto.NewCompressionProcessReport(
		[]tkValueObject.UnixAbsoluteFilePath{},
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

	filesExists := repo.fileClerk.FileExists(filePathStr)
	if filesExists {
		return errors.New("PathAlreadyExists")
	}

	unixUser, err := user.LookupId(createDto.OperatorAccountId.String())
	if err != nil {
		return errors.New("AccountNotFound")
	}

	fileOwnershipStr := unixUser.Username + ":" + unixUser.Username
	fileOwner, err := tkValueObject.NewUnixFileOwnership(fileOwnershipStr)
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

func (repo FilesCmdRepo) Delete(unixFilePath tkValueObject.UnixAbsoluteFilePath) error {
	fileExists := repo.fileClerk.FileExists(unixFilePath.String())
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

	fileToExtractExists := repo.fileClerk.FileExists(fileToExtract.String())
	if !fileToExtractExists {
		return errors.New("FileNotFound")
	}

	destinationPath := extractDto.DestinationPath

	destinationPathExists := repo.fileClerk.FileExists(destinationPath.String())
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

	err = repo.fileClerk.CreateDir(destinationPath.String())
	if err != nil {
		return err
	}

	compressCmd := fmt.Sprintf(
		"%s %s %s %s %s",
		compressBinary, compressBinaryFlag, fileToExtract.String(),
		compressDestinationFlag, destinationPath.String(),
	)
	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command:          compressCmd,
		ShouldUseSubShell: true,
	}).Run()
	return err
}

func (repo FilesCmdRepo) Move(moveDto dto.MoveUnixFile) error {
	sourcePathStr := moveDto.SourcePath.String()
	if !repo.fileClerk.FileExists(sourcePathStr) {
		return errors.New("SourceFileNotFound")
	}

	if moveDto.DestinationPath == valueObject.UnixFilePathTrashDir {
		fileNameStr := moveDto.SourcePath.ReadFileName(false).String()
		destinationPathStr := moveDto.DestinationPath.String()
		rawTrashFilePath := destinationPathStr + "/" + fileNameStr
		trashFilePath, err := valueObject.NewUnixFilePath(rawTrashFilePath)
		if err != nil {
			return errors.New("DefineTrashFilePathError: " + err.Error())
		}

		trashFilePathStr := trashFilePath.String()
		if repo.fileClerk.FileExists(trashFilePathStr) {
			uniqueTrashPathStr := trashFilePathStr + "-" + tkValueObject.NewUnixTimeNow().String()
			uniqueTrashFilePath, err := valueObject.NewUnixFilePath(uniqueTrashPathStr)
			if err != nil {
				return errors.New("DefineUniqueTrashFilePathError: " + err.Error())
			}

			trashFilePath = uniqueTrashFilePath
			trashFilePathStr = uniqueTrashFilePath.String()
		}

		return os.Rename(sourcePathStr, trashFilePathStr)
	}

	destinationPathStr := moveDto.DestinationPath.String()
	if repo.fileClerk.FileExists(destinationPathStr) {
		if !moveDto.ShouldOverwrite {
			return errors.New("DestinationPathAlreadyExists")
		}

		err := repo.Delete(moveDto.DestinationPath)
		if err != nil {
			return errors.New("DeletePreviousDestinationFileError: " + err.Error())
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

	return repo.fileClerk.UpdateFileContent(
		updateContentDto.SourcePath.String(), decodedContent, true,
	)
}

func (repo FilesCmdRepo) UpdateOwnership(
	updateOwnershipDto dto.UpdateUnixFileOwnership,
) error {
	sourcePathStr := updateOwnershipDto.SourcePath.String()
	if !repo.fileClerk.FileExists(sourcePathStr) {
		return errors.New("FileNotFound")
	}

	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command: "chown",
		Args:    []string{updateOwnershipDto.Ownership.String(), sourcePathStr},
	}).Run()
	if err != nil {
		return errors.New("UpdateFileOwnershipError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) UpdatePermissions(
	updatePermissionsDto dto.UpdateUnixFilePermissions,
) error {
	sourcePathStr := updatePermissionsDto.SourcePath.String()
	if !repo.fileClerk.FileExists(sourcePathStr) {
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

	_, err := tkInfra.NewShell(tkInfra.ShellSettings{
		Command:          updatePermissionsCmd,
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return errors.New("UpdatePermissionsError: " + err.Error())
	}

	return nil
}

func (repo FilesCmdRepo) Upload(
	uploadDto dto.UploadUnixFiles,
) (dto.UploadProcessReport, error) {
	uploadProcessReport := dto.NewUploadProcessReport(
		[]tkValueObject.UnixFileName{},
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
