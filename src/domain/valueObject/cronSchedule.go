package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const cronScheduleCustomFrequencyRegex string = `^((@?(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(?:ns|us|µs|ms|s|m|h))+))$`
const cronScheduleFrequencyRegex string = `^(((?P<minute>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<hour>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<day>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<month>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<weekday>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )?)$`

type CronSchedule string

func NewCronSchedule(value string) (CronSchedule, error) {
	if shouldHaveAtSign(value) {
		hasAtSign := strings.HasPrefix(value, "@")
		if !hasAtSign {
			value = "@" + value
		}
		return CronSchedule(value), nil
	}

	schedule := CronSchedule(value)

	if !schedule.isValid() {
		return "", errors.New("InvalidCronSchedule")
	}

	return schedule, nil
}

func NewCronSchedulePanic(value string) CronSchedule {
	schedule, err := NewCronSchedule(value)
	if err != nil {
		panic(err)
	}
	return schedule
}

func shouldHaveAtSign(value string) bool {
	frequencyRegex := regexp.MustCompile(cronScheduleCustomFrequencyRegex)
	return frequencyRegex.MatchString(value)
}

func (schedule CronSchedule) isValid() bool {
	scheduleRe := regexp.MustCompile(cronScheduleFrequencyRegex)
	return scheduleRe.MatchString(string(schedule))
}

func (schedule CronSchedule) String() string {
	return string(schedule)
}
