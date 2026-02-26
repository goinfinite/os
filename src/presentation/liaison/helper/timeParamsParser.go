package liaisonHelper

import (
	"log/slog"

	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func TimeParamsParser(
	timeParamNames []string,
	untrustedInput map[string]any,
) map[string]*tkValueObject.UnixTime {
	timeParamPtrs := map[string]*tkValueObject.UnixTime{}

	for _, timeParamName := range timeParamNames {
		if untrustedInput[timeParamName] == nil {
			continue
		}

		timeParam, err := tkValueObject.NewUnixTime(untrustedInput[timeParamName])
		if err != nil {
			slog.Debug("InvalidTimeParam", slog.String("timeParamName", timeParamName))
			timeParamPtrs[timeParamName] = nil
			continue
		}

		timeParamPtrs[timeParamName] = &timeParam
	}

	return timeParamPtrs
}
