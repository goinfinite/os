package cronInfra

import (
	"errors"
	"strings"

	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
)

type CronQueryRepo struct {
}

func (repo CronQueryRepo) cronFactory(
	cronIndex int,
	cronLine string,
) (entity.Cron, error) {
	cronRegex := `^(?P<frequency>(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)|((((\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+) ?){5,7}))(?P<cmd>[^#\r\n]{1,1000})(?P<comment>#(.*)){0,1000}$`
	namedGroupMap := voHelper.FindNamedGroupsMatches(cronRegex, cronLine)

	var cron entity.Cron
	id, err := valueObject.NewCronId(cronIndex)
	if err != nil {
		return cron, errors.New("CronIdError")
	}

	if namedGroupMap["frequency"] == "" {
		return cron, errors.New("CronFrequencyError")
	}
	schedule, err := valueObject.NewCronSchedule(
		strings.TrimSpace(namedGroupMap["frequency"]),
	)
	if err != nil {
		return cron, errors.New("CronScheduleError")
	}

	if namedGroupMap["cmd"] == "" {
		return cron, errors.New("CronCommandError")
	}
	cmd, err := valueObject.NewUnixCommand(
		strings.TrimSpace(namedGroupMap["cmd"]),
	)
	if err != nil {
		return cron, errors.New("CronCommandError")
	}

	var cronCommentPtr *valueObject.CronComment
	if namedGroupMap["comment"] != "" {
		commentWithoutLeadingHash := strings.Trim(namedGroupMap["comment"], "#")
		cronComment, err := valueObject.NewCronComment(
			strings.TrimSpace(commentWithoutLeadingHash),
		)
		if err != nil {
			return cron, errors.New("CronCommentError")
		}
		cronCommentPtr = &cronComment
	}

	return entity.NewCron(id, schedule, cmd, cronCommentPtr), nil
}

func (repo CronQueryRepo) Read() ([]entity.Cron, error) {
	crons := []entity.Cron{}

	cronOut, err := infraHelper.RunCmd("crontab", "-l")
	if err != nil {
		if strings.Contains(err.Error(), "no crontab") {
			return crons, nil
		}
		return crons, errors.New("CrontabReadError: " + err.Error())
	}

	cronLines := strings.Split(cronOut, "\n")
	if len(cronLines) == 0 {
		return crons, nil
	}

	for cronIndex, cronLine := range cronLines {
		if cronLine == "" {
			continue
		}

		if strings.HasPrefix(cronLine, "#") {
			continue
		}
		cronLineIndex := cronIndex + 1
		cron, err := repo.cronFactory(cronLineIndex, cronLine)
		if err != nil {
			continue
		}
		crons = append(crons, cron)
	}

	return crons, nil
}

func (repo CronQueryRepo) ReadById(
	cronId valueObject.CronId,
) (cronEntity entity.Cron, err error) {
	crons, err := repo.Read()
	if err != nil {
		return cronEntity, err
	}

	if len(crons) == 0 {
		return cronEntity, errors.New("CronNotFound")
	}

	cronIdStr := cronId.String()
	for _, cron := range crons {
		if cron.Id.String() != cronIdStr {
			continue
		}

		return cron, nil
	}

	return cronEntity, errors.New("CronNotFound")
}
