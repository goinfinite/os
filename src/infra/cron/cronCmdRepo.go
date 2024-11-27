package cronInfra

import (
	"os"

	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	"github.com/google/uuid"
)

type CronCmdRepo struct {
	currentCrontab     []entity.Cron
	tmpCrontabFilename string
}

func NewCronCmdRepo() (*CronCmdRepo, error) {
	cronQueryRepo := CronQueryRepo{}

	currentCrontab, err := cronQueryRepo.Read()
	if err != nil {
		return nil, err
	}

	tmpCrontabFilename := "tmpCrontab_" + uuid.NewString()

	return &CronCmdRepo{
		currentCrontab:     currentCrontab,
		tmpCrontabFilename: tmpCrontabFilename,
	}, nil
}

func (repo *CronCmdRepo) createCrontabTmpFile() error {
	tmpCrontabFile, err := os.Create(repo.tmpCrontabFilename)
	if err != nil {
		return err
	}
	defer tmpCrontabFile.Close()

	return nil
}

func (repo *CronCmdRepo) installNewCrontab() error {
	err := repo.createCrontabTmpFile()
	if err != nil {
		return nil
	}

	var crontabContent string
	for _, cron := range repo.currentCrontab {
		crontabContent += cron.String() + "\n"
	}

	err = infraHelper.UpdateFile(repo.tmpCrontabFilename, crontabContent, true)
	if err != nil {
		return err
	}

	_, err = infraHelper.RunCmd("crontab", repo.tmpCrontabFilename)
	if err != nil {
		return err
	}

	return os.Remove(repo.tmpCrontabFilename)
}

func (repo *CronCmdRepo) Create(createCron dto.CreateCron) (valueObject.CronId, error) {
	cronsCount := len(repo.currentCrontab)
	newCronIndex := cronsCount + 1

	cronId, err := valueObject.NewCronId(newCronIndex)
	if err != nil {
		return cronId, err
	}

	newCron := entity.NewCron(
		cronId, createCron.Schedule, createCron.Command, createCron.Comment,
	)

	repo.currentCrontab = append(repo.currentCrontab, newCron)

	return cronId, repo.installNewCrontab()
}

func (repo *CronCmdRepo) Update(updateCron dto.UpdateCron) error {
	cronToUpdateId := updateCron.Id
	cronToUpdateListIndex := cronToUpdateId.Uint64() - 1

	newCronSchedule := repo.currentCrontab[cronToUpdateListIndex].Schedule
	if updateCron.Schedule != nil {
		newCronSchedule = *updateCron.Schedule
	}

	newCronCommand := repo.currentCrontab[cronToUpdateListIndex].Command
	if updateCron.Command != nil {
		newCronCommand = *updateCron.Command
	}

	newCronComment := repo.currentCrontab[cronToUpdateListIndex].Comment
	if updateCron.Comment != nil {
		newCronComment = updateCron.Comment
	}

	newCron := entity.NewCron(
		cronToUpdateId, newCronSchedule, newCronCommand, newCronComment,
	)
	repo.currentCrontab[cronToUpdateListIndex] = newCron

	return repo.installNewCrontab()
}

func (repo *CronCmdRepo) Delete(cronId valueObject.CronId) error {
	var cronsToKeep []entity.Cron
	for _, currentCron := range repo.currentCrontab {
		if cronId.Uint64() == currentCron.Id.Uint64() {
			continue
		}
		cronsToKeep = append(cronsToKeep, currentCron)
	}
	repo.currentCrontab = cronsToKeep

	return repo.installNewCrontab()
}
