package valueObject

import (
	"errors"
	"regexp"
)

const cronScheduleRegex string = `^((?P<frequencyStr>(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)) ?|((?P<minute>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<hour>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<day>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<month>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<weekday>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )?)$`

type CronSchedule string

func NewCronSchedule(value string) (CronSchedule, error) {
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

func (schedule CronSchedule) isValid() bool {
	re := regexp.MustCompile(cronScheduleRegex)
	return re.MatchString(string(schedule))
}

func (schedule CronSchedule) String() string {
	return string(schedule)
}
