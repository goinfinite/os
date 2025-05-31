package cliHelper

import (
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	tkPresentation "github.com/goinfinite/tk/src/presentation"
)

func RequestInputParser(rawRequestInput map[string]any) tkPresentation.RequestInputParsed {
	return tkPresentation.RequestInputParser(
		tkPresentation.RequestInputSettings{
			RawRequestInput:        rawRequestInput,
			KnownParamConstructors: sharedHelper.KnownParamConstructors,
		},
	)
}
