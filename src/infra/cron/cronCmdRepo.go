package cronInfra

import (
	"errors"
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkDto "github.com/goinfinite/tk/src/domain/dto"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
	tkInfra "github.com/goinfinite/tk/src/infra"
)

type CronCmdRepo struct {
	cronQueryRepo *CronQueryRepo
	fileClerk     tkInfra.FileClerk
}

func NewCronCmdRepo() *CronCmdRepo {
	return &CronCmdRepo{
		cronQueryRepo: NewCronQueryRepo(),
		fileClerk:     tkInfra.FileClerk{},
	}
}

func (repo *CronCmdRepo) rebuildCrontab(cronsEntities []entity.Cron) error {
	tmpCrontabDirPath, err := os.MkdirTemp("", "crontab-")
	if err != nil {
		return errors.New("CreateCrontabTempDirError: " + err.Error())
	}
	defer func() {
		removeErr := os.RemoveAll(tmpCrontabDirPath)
		if removeErr != nil {
			slog.Error(
				"DeleteCrontabTempDirFailed",
				slog.String("err", removeErr.Error()),
			)
		}
	}()

	tmpCrontabFilePath, err := tkValueObject.NewUnixAbsoluteFilePath(
		tmpCrontabDirPath+"/crontab", false,
	)
	if err != nil {
		return errors.New("DefineCrontabTempFilePathError: " + err.Error())
	}

	var crontabContent strings.Builder
	for _, cronEntity := range cronsEntities {
		crontabContent.WriteString(cronEntity.String())
		crontabContent.WriteString("\n")
	}

	crontabFilePermissions := os.FileMode(0644)
	err = repo.fileClerk.UpsertFile(tkInfra.FileUpsertSettings{
		FilePath:        tmpCrontabFilePath,
		Permissions:     &crontabFilePermissions,
		OverwritePolicy: &tkInfra.FileClerkOverwritePolicyReplace,
	}, []byte(crontabContent.String()))
	if err != nil {
		return errors.New("UpdateCrontabTempFileContentError: " + err.Error())
	}

	_, err = tkInfra.NewShell(tkInfra.ShellSettings{
		Command:           "crontab " + tmpCrontabFilePath.String(),
		ShouldUseSubShell: true,
	}).Run()
	if err != nil {
		return err
	}

	return nil
}

func (repo *CronCmdRepo) Create(
	createDto dto.CreateCron,
) (cronId valueObject.CronId, err error) {
	readRequestDto := dto.ReadCronsRequest{
		Pagination: tkDto.Pagination{
			ItemsPerPage: 1000,
		},
	}
	readResponseDto, err := repo.cronQueryRepo.Read(readRequestDto)
	if err != nil {
		return cronId, errors.New("ReadCronsError: " + err.Error())
	}
	cronsEntities := readResponseDto.Crons

	rawCronId := len(readResponseDto.Crons) + 1
	cronId, err = valueObject.NewCronId(rawCronId)
	if err != nil {
		return cronId, err
	}

	newCron := entity.NewCron(
		cronId, createDto.Schedule, createDto.Command, createDto.Comment,
	)
	cronsEntities = append(cronsEntities, newCron)

	return cronId, repo.rebuildCrontab(cronsEntities)
}

func (repo *CronCmdRepo) Update(updateDto dto.UpdateCron) error {
	readRequestDto := dto.ReadCronsRequest{
		Pagination: tkDto.Pagination{
			ItemsPerPage: 1000,
		},
	}
	readResponseDto, err := repo.cronQueryRepo.Read(readRequestDto)
	if err != nil {
		return errors.New("ReadCronsError: " + err.Error())
	}
	cronsEntities := readResponseDto.Crons

	desiredCronIndex := updateDto.Id.Uint64() - 1
	desiredCron := cronsEntities[desiredCronIndex]

	schedule := desiredCron.Schedule
	if updateDto.Schedule != nil {
		schedule = *updateDto.Schedule
	}

	command := desiredCron.Command
	if updateDto.Command != nil {
		command = *updateDto.Command
	}

	comment := desiredCron.Comment
	if updateDto.Comment != nil {
		comment = updateDto.Comment
	}
	if slices.Contains(updateDto.ClearableFields, "comment") {
		comment = nil
	}

	desiredCronWithUpdatedValues := entity.NewCron(
		desiredCron.Id, schedule, command, comment,
	)
	cronsEntities[desiredCronIndex] = desiredCronWithUpdatedValues

	return repo.rebuildCrontab(cronsEntities)
}

func (repo *CronCmdRepo) Delete(cronId valueObject.CronId) error {
	readRequestDto := dto.ReadCronsRequest{
		Pagination: tkDto.Pagination{
			ItemsPerPage: 1000,
		},
	}
	readResponseDto, err := repo.cronQueryRepo.Read(readRequestDto)
	if err != nil {
		return errors.New("ReadCronsError: " + err.Error())
	}

	cronsEntitiesToKeep := []entity.Cron{}
	cronFound := false
	for _, cronEntity := range readResponseDto.Crons {
		if cronEntity.Id == cronId {
			cronFound = true
			continue
		}

		cronsEntitiesToKeep = append(cronsEntitiesToKeep, cronEntity)
	}

	if !cronFound {
		return errors.New("CronNotFound")
	}

	return repo.rebuildCrontab(cronsEntitiesToKeep)
}
