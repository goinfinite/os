package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const cronScheduleFrequencyRegex string = `^(?:(?:@?(?:annually|yearly|monthly|weekly|daily|hourly|reboot))|(?:@every (?:\d+(?:ns|us|µs|ms|s|m|h))+))$`
const cronScheduleRegex string = `^(((?P<minute>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<hour>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<day>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<month>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<weekday>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )?)$`

type CronSchedule string

func NewCronSchedule(value string) (CronSchedule, error) {
	schedule := CronSchedule(value)
	if !schedule.isValid() {
		return "", errors.New("InvalidCronSchedule")
	}

	if !schedule.hasAtSign() {
		schedule = CronSchedule("@" + value)
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

func (schedule CronSchedule) hasAtSign() bool {
	frequencyRegex := regexp.MustCompile(cronScheduleFrequencyRegex)
	frequencyGroup := frequencyRegex.FindStringSubmatch(string(schedule))

	if len(frequencyGroup) > 0 {
		return strings.HasPrefix(string(schedule), "@")
	}

	return false
}

func (schedule CronSchedule) isValid() bool {
	frequencyRegex := regexp.MustCompile(cronScheduleFrequencyRegex)
	frequencyMatch := frequencyRegex.MatchString(string(schedule))

	if frequencyMatch {
		return true
	}

	scheduleRe := regexp.MustCompile(cronScheduleRegex)
	return scheduleRe.MatchString(string(schedule))
}

func (schedule CronSchedule) String() string {
	return string(schedule)
}
