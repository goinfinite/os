package cronInfra

import (
	"errors"
	"os"

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

func (repo *CronCmdRepo) updateCrontabFile(cronsList []entity.Cron) error {
	tmpCrontabFilePath := "/tmp/crontab"

	if !infraHelper.FileExists(tmpCrontabFilePath) {
		_, err := os.Create(tmpCrontabFilePath)
		if err != nil {
			return errors.New("CreateCrontabTempFileError: " + err.Error())
		}
	}

	crontabContent := ""
	for _, cron := range cronsList {
		crontabContent += cron.String() + "\n"
	}

	err := infraHelper.UpdateFile(tmpCrontabFilePath, crontabContent, true)
	if err != nil {
		return errors.New("UpdateCrontabTempFileContentError: " + err.Error())
	}

	_, err = infraHelper.RunCmdWithSubShell("crontab " + tmpCrontabFilePath)
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

	cronsList := readResponseDto.Crons

	rawCronId := len(readResponseDto.Crons) + 1
	cronId, err = valueObject.NewCronId(rawCronId)
	if err != nil {
		return cronId, err
	}

	newCron := entity.NewCron(
		cronId, createDto.Schedule, createDto.Command, createDto.Comment,
	)
	cronsList = append(cronsList, newCron)

	return cronId, repo.updateCrontabFile(cronsList)
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
	crons := readResponseDto.Crons

	desiredCronIndex := updateDto.Id.Uint64() - 1
	desiredCron := crons[desiredCronIndex]

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

	id, err := valueObject.NewCronId(desiredCronIndex)
	if err != nil {
		return err
	}

	desiredCronWithUpdatedValues := entity.NewCron(id, schedule, command, comment)

	crons[desiredCronIndex] = desiredCronWithUpdatedValues

	return repo.updateCrontabFile(crons)
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

	desiredCronIndex := cronId.Uint64() - 1

	cronsToKeep := append(
		readResponseDto.Crons[:desiredCronIndex],
		readResponseDto.Crons[desiredCronIndex+1:]...,
	)

	return repo.updateCrontabFile(cronsToKeep)
}
