package valueObject

import (
	"errors"
	"regexp"
	"strings"

	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
)

const cronScheduleRegex string = `^((?P<frequencyStr>(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)) ?|((?P<minute>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<hour>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<day>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<month>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<weekday>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )?)$`

type CronSchedule string

func NewCronSchedule(value interface{}) (cronSchedule CronSchedule, err error) {
	stringValue, err := voHelper.InterfaceToString(value)
	if err != nil {
		return cronSchedule, errors.New("CronScheduleMustBeString")
	}

	if shouldHaveAtSign(stringValue) {
		hasAtSign := strings.HasPrefix(stringValue, "@")
		if !hasAtSign {
			stringValue = "@" + stringValue
		}
	}

	re := regexp.MustCompile(cronScheduleRegex)
	if !re.MatchString(stringValue) {
		return cronSchedule, errors.New("InvalidCronSchedule")
	}

	return CronSchedule(stringValue), nil
}

func shouldHaveAtSign(value string) bool {
	cronPredefinedScheduleRegex := `^((@?(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(?:ns|us|µs|ms|s|m|h))+))$`
	frequencyRegex := regexp.MustCompile(cronPredefinedScheduleRegex)
	return frequencyRegex.MatchString(value)
}

func (vo CronSchedule) String() string {
	return string(vo)
}
