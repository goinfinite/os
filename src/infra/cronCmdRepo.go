package infra

import (
	"os"
	"strconv"
	"time"

	"github.com/speedianet/sam/src/domain/dto"
	infraHelper "github.com/speedianet/sam/src/infra/helper"
)

type CronCmdRepo struct {
}

func (repo CronCmdRepo) Add(addCron dto.AddCron) error {
	cronUnixTimestampStr := strconv.FormatInt(time.Now().Unix(), 10)

	err := importCurrentCrontab(cronUnixTimestampStr)
	if err != nil {
		return err
	}

	cronJob := addCron.Schedule.String() + " " + addCron.Command.String() + " # " + addCron.Comment.String()

	err = infraHelper.UpdateFile("newCrontab_"+cronUnixTimestampStr, cronJob+"\n", false)
	if err != nil {
		return err
	}

	err = installNewCrontab(cronUnixTimestampStr)
	if err != nil {
		return err
	}

	return nil
}

func importCurrentCrontab(timeSuffix string) error {
	crontabFileName := "newCrontab_" + timeSuffix

	currentCrontab, err := infraHelper.RunCmd(
		"crontab",
		"-l",
	)
	if err != nil {
		return err
	}

	tmpCrontabFile, err := os.Create(crontabFileName)
	_ = tmpCrontabFile.Close()

	currentCrontabLen := len(currentCrontab)
	if currentCrontabLen > 0 {
		currentCrontab += "\n"
	}

	err = infraHelper.UpdateFile(crontabFileName, currentCrontab, true)
	if err != nil {
		return err
	}

	return nil
}

func installNewCrontab(timeSuffix string) error {
	crontabFileName := "newCrontab_" + timeSuffix

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
