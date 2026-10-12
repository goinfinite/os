package valueObject

import (
	"errors"
	"regexp"
	"strings"

	tkVoUtil "github.com/goinfinite/tk/src/domain/valueObject/util"
)

var cronScheduleRegex = regexp.MustCompile(
	`^((?P<frequencyStr>(@(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(ns|us|µs|ms|s|m|h))+)) ?|((?P<minute>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<hour>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<day>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<month>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )((?P<weekday>(\d+,)+\d+|(\d+(\/|-)\d+)|\d+|\*|\*/\d+){1})(?: )?)$`,
)

var cronPredefinedScheduleRegex = regexp.MustCompile(
	`^((@?(annually|yearly|monthly|weekly|daily|hourly|reboot))|(@every (\d+(?:ns|us|µs|ms|s|m|h))+))$`,
)

type CronSchedule string

func NewCronSchedule(value any) (cronSchedule CronSchedule, err error) {
	stringValue, err := tkVoUtil.InterfaceToString(value)
	if err != nil {
		return cronSchedule, errors.New("CronScheduleMustBeString")
	}

	if shouldHaveAtSign(stringValue) {
		hasAtSign := strings.HasPrefix(stringValue, "@")
		if !hasAtSign {
			stringValue = "@" + stringValue
		}
	}

	if !cronScheduleRegex.MatchString(stringValue) {
		return cronSchedule, errors.New("InvalidCronSchedule")
	}

	return CronSchedule(stringValue), nil
}

func shouldHaveAtSign(value string) bool {
	return cronPredefinedScheduleRegex.MatchString(value)
}

func (vo CronSchedule) String() string {
	return string(vo)
}
