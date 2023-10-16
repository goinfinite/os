package valueObject

import (
	"errors"
	"regexp"
	"strings"
)

const cronSchedulePredefinedFrequencyRegex string = `^(?:(?:@?(?:annually|yearly|monthly|weekly|daily|hourly|reboot))|(?:@every (?:\d+(?:ns|us|µs|ms|s|m|h))+))$`
const cronScheduleFrequencyRegex string = `^(((?P<minute>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<hour>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<day>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<month>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<weekday>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )?)$`

type CronSchedule string

func NewCronSchedule(value string) (CronSchedule, error) {
	schedule := CronSchedule(value)

	if schedule.shouldHaveAtSign() {
		hasAtSign := strings.HasPrefix(string(schedule), "@")
		if !hasAtSign {
			schedule = CronSchedule("@" + value)
		}
	}

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

func (schedule CronSchedule) shouldHaveAtSign() bool {
	frequencyRegex := regexp.MustCompile(cronSchedulePredefinedFrequencyRegex)
	frequencyMatch := frequencyRegex.MatchString(string(schedule))

	return frequencyMatch
}

func (schedule CronSchedule) isValid() bool {
	frequencyRegex := regexp.MustCompile(cronSchedulePredefinedFrequencyRegex)
	frequencyMatch := frequencyRegex.MatchString(string(schedule))

	if frequencyMatch {
		return true
	}

	scheduleRe := regexp.MustCompile(cronScheduleFrequencyRegex)
	return scheduleRe.MatchString(string(schedule))
}

func (schedule CronSchedule) String() string {
	return string(schedule)
}
