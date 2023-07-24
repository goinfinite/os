package apiHelper

import (
	"strings"
)

func CheckMissingParams(
	requestBody map[string]interface{},
	requiredParams []string,
) {
	missingParams := []string{}
	for _, param := range requiredParams {
		if _, ok := requestBody[param]; !ok {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) > 0 {
		panic("MissingParams: " + strings.Join(missingParams, ", "))
	}
}
