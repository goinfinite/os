package apiHelper

import tkPresentation "github.com/goinfinite/tk/src/presentation"

func CheckMissingParams(
	requestBody map[string]interface{},
	requiredParams []string,
) {
	err := tkPresentation.RequiredParamsInspector(requestBody, requiredParams)
	if err != nil {
		panic(err.Error())
	}
}
