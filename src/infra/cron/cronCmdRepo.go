package cronInfra

import (
	"errors"
	"os"
	"slices"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type CronCmdRepo struct {
	cronQueryRepo *CronQueryRepo
}

func NewCronCmdRepo() *CronCmdRepo {
	return &CronCmdRepo{
		cronQueryRepo: NewCronQueryRepo(),
	}
}

func (repo *CronCmdRepo) rebuildCrontab(cronsEntities []entity.Cron) error {
	tmpCrontabFilePath := "/tmp/crontab"

	if !infraHelper.FileExists(tmpCrontabFilePath) {
		_, err := os.Create(tmpCrontabFilePath)
		if err != nil {
			return errors.New("CreateCrontabTempFileError: " + err.Error())
		}
	}

	crontabContent := ""
	for _, cronEntity := range cronsEntities {
		crontabContent += cronEntity.String() + "\n"
	}

	shouldOverwrite := true
	err := infraHelper.UpdateFile(tmpCrontabFilePath, crontabContent, shouldOverwrite)
	if err != nil {
		return errors.New("UpdateCrontabTempFileContentError: " + err.Error())
	}

	_, err = infraHelper.RunCmd(infraHelper.RunCmdSettings{
		Command:               "crontab " + tmpCrontabFilePath,
		ShouldRunWithSubShell: true,
	})
	if err != nil {
		return err
	}

	err = os.Remove(tmpCrontabFilePath)
	if err != nil {
		return errors.New("DeleteCrontabTempFileError: " + err.Error())
	}

	return nil
}

func (repo *CronCmdRepo) Create(
	createDto dto.CreateCron,
) (cronId valueObject.CronId, err error) {
	readRequestDto := dto.ReadCronsRequest{
		Pagination: dto.Pagination{
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
		Pagination: dto.Pagination{
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
		Pagination: dto.Pagination{
			ItemsPerPage: 1000,
		},
	}
	readResponseDto, err := repo.cronQueryRepo.Read(readRequestDto)
	if err != nil {
		return errors.New("ReadCronsError: " + err.Error())
	}
	cronsEntities := readResponseDto.Crons
	cronEntityIndex := cronId.Uint64() - 1

	cronsEntitiesToKeep := append(
		cronsEntities[:cronEntityIndex],
		cronsEntities[cronEntityIndex+1:]...,
	)

	return repo.rebuildCrontab(cronsEntitiesToKeep)
}
