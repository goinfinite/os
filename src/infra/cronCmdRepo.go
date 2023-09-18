package infra

import (
	"os"
	"strconv"
	"time"

	"github.com/speedianet/sam/src/domain/dto"
	"github.com/speedianet/sam/src/domain/entity"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type CronCmdRepo struct {
}

func (repo CronCmdRepo) Add(addCron dto.AddCron) error {
	cronUnixTimestampStr := strconv.FormatInt(time.Now().Unix(), 10)
	cronjob := addCron.Schedule.String() + " " + addCron.Command.String() + " # " + addCron.Comment.String()

	err := editCrontab(cronjob, cronUnixTimestampStr, false)
	if err != nil {
		return err
	}

	return nil
}

func (repo CronCmdRepo) Update(cron entity.Cron, updateCron dto.UpdateCron) error {
	var cronjobSchedule string
	var cronjobCommand string
	var cronjobComment string

	cronUnixTimestampStr := strconv.FormatInt(time.Now().Unix(), 10)

	cronjobSchedule = cron.Schedule.String()
	if updateCron.Schedule != nil {
		cronjobSchedule = updateCron.Schedule.String()
	}

	cronjobCommand = cron.Command.String()
	if updateCron.Command != nil {
		cronjobCommand = updateCron.Command.String()
	}

	cronjobComment = cron.Comment.String()
	if updateCron.Comment != nil {
		cronjobComment = updateCron.Comment.String()
	}

	err := editCrontab(
		cronjobSchedule+" "+cronjobCommand+" # "+cronjobComment,
		cronUnixTimestampStr,
		false,
	)
	if err != nil {
		return err
	}

	err = removeCronjob(updateCron.Id.String(), cronUnixTimestampStr)
	if err != nil {
		return err
	}

	return nil
}

func editCrontab(crontabContent string, cronUnixTimestampStr string, delete bool) error {
	err := importCurrentCrontab(cronUnixTimestampStr)
	if err != nil {
		return err
	}

	shouldOverwrite := delete

	err = infraHelper.UpdateFile("cron_"+cronUnixTimestampStr, crontabContent+"\n", shouldOverwrite)
	if err != nil {
		return err
	}

	err = installNewCrontab(cronUnixTimestampStr)
	if err != nil {
		return err
	}

	return nil
}

func removeCronjob(line string, cronUnixTimestampStr string) error {
	crontabWithoutSpecificLine, err := infraHelper.RunCmd(
		"bash",
		"-c",
		"crontab -l | sed '"+line+"d'",
	)
	if err != nil {
		return err
	}

	err = editCrontab(crontabWithoutSpecificLine, cronUnixTimestampStr, true)
	if err != nil {
		return err
	}

	return nil
}

func importCurrentCrontab(timeSuffix string) error {
	crontabFileName := "cron_" + timeSuffix

	currentCrontab, err := infraHelper.RunCmd(
		"crontab",
		"-l",
	)
	if err != nil {
		return err
	}

	currentCrontabLen := len(currentCrontab)
	if currentCrontabLen > 0 {
		currentCrontab += "\n"
	}

	tmpCrontabFile, err := os.Create(crontabFileName)
	_, err = tmpCrontabFile.WriteString(currentCrontab)
	if err != nil {
		return err
	}
	_ = tmpCrontabFile.Close()

	return nil
}

func installNewCrontab(timeSuffix string) error {
	crontabFileName := "cron_" + timeSuffix

	_, err := infraHelper.RunCmd(
		"crontab",
		crontabFileName,
	)
	if err != nil {
		return err
	}

	err = os.Remove(crontabFileName)
	if err != nil {
		return err
	}

	return nil
}
