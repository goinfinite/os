package serviceHelper

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func TimeParamsParser(
	timeParamNames []string,
	userInput map[string]interface{},
) map[string]*valueObject.UnixTime {
	timeParamPtrs := map[string]*valueObject.UnixTime{}

	for _, timeParamName := range timeParamNames {
		if userInput[timeParamName] == nil {
			continue
		}

		timeParam, err := valueObject.NewUnixTime(userInput[timeParamName])
		if err != nil {
			slog.Debug("InvalidTimeParam", slog.String("timeParamName", timeParamName))
			timeParamPtrs[timeParamName] = nil
			continue
		}

		timeParamPtrs[timeParamName] = &timeParam
	}

	return timeParamPtrs
}
