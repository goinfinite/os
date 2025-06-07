package liaisonHelper

import (
	"log/slog"

	"github.com/goinfinite/os/src/domain/valueObject"
)

func TimeParamsParser(
	timeParamNames []string,
	untrustedInput map[string]any,
) map[string]*valueObject.UnixTime {
	timeParamPtrs := map[string]*valueObject.UnixTime{}

	for _, timeParamName := range timeParamNames {
		if untrustedInput[timeParamName] == nil {
			continue
		}

		timeParam, err := valueObject.NewUnixTime(untrustedInput[timeParamName])
		if err != nil {
			slog.Debug("InvalidTimeParam", slog.String("timeParamName", timeParamName))
			timeParamPtrs[timeParamName] = nil
			continue
		}

		timeParamPtrs[timeParamName] = &timeParam
	}

	return timeParamPtrs
}
