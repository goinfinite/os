package infra

import (
	"os"

	"github.com/google/uuid"
	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	"github.com/speedianet/sam/src/domain/valueObject"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type CronCmdRepo struct {
	currentCrontab     []entity.Cron
	tmpCrontabFilename string
}

func NewCronCmdRepo() (*CronCmdRepo, error) {
	cronQueryRepo := CronQueryRepo{}

	currentCrontab, err := cronQueryRepo.Get()
	if err != nil {
		return nil, err
	}

	tmpCrontabFilename := "tmpCrontab_" + uuid.NewString()

	return &CronCmdRepo{
		currentCrontab:     currentCrontab,
		tmpCrontabFilename: tmpCrontabFilename,
	}, nil
}

func (repo CronCmdRepo) createCrontabTmpFile() error {
	tmpCrontabFile, err := os.Create(repo.tmpCrontabFilename)
	if err != nil {
		return err
	}
	defer tmpCrontabFile.Close()

	return nil
}

func fromCronEntityToCronStr(cron entity.Cron) string {
	return cron.Schedule.String() + " " +
		cron.Command.String() + " # " +
		cron.Comment.String() + "\n"
}

func (repo CronCmdRepo) installNewCrontab() error {
	err := repo.createCrontabTmpFile()
	if err != nil {
		return nil
	}

	var crontabContent string
	for _, currentCrontabContent := range repo.currentCrontab {
		crontabContent += fromCronEntityToCronStr(currentCrontabContent)
	}

	err = infraHelper.UpdateFile(repo.tmpCrontabFilename, crontabContent, true)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd(
		"crontab",
		repo.tmpCrontabFilename,
	)
	if err != nil {
		return err
	}

	return os.Remove(repo.tmpCrontabFilename)
}

func (repo CronCmdRepo) Add(addCron dto.AddCron) error {
	cronLineIndex := len(repo.currentCrontab) + 1
	cronId, err := valueObject.NewCronId(cronLineIndex)
	if err != nil {
		return err
	}

	newCron := entity.NewCron(
		cronId,
		addCron.Schedule,
		addCron.Command,
		addCron.Comment,
	)

	repo.currentCrontab = append(repo.currentCrontab, newCron)

	return repo.installNewCrontab()
}

func (repo CronCmdRepo) Update(updateCron dto.UpdateCron) error {
	cronIndex := updateCron.Id.Get() - 1

	if updateCron.Schedule != nil {
		repo.currentCrontab[cronIndex].Schedule = *updateCron.Schedule
	}

	if updateCron.Command != nil {
		repo.currentCrontab[cronIndex].Command = *updateCron.Command
	}

	if updateCron.Comment != nil {
		repo.currentCrontab[cronIndex].Comment = updateCron.Comment
	}

	err := repo.installNewCrontab()
	return err
}

func (repo CronCmdRepo) Delete(cronId valueObject.CronId) error {
	cronIndex := cronId.Get() - 1

	repo.currentCrontab = append(
		repo.currentCrontab[:cronIndex],
		repo.currentCrontab[cronIndex+1:]...,
	)

	err := repo.installNewCrontab()
	return err
}
