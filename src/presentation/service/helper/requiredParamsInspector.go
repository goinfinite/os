package serviceHelper

import (
	"errors"
	"strings"
)

func RequiredParamsInspector(
	input map[string]interface{},
	requiredParams []string,
) error {
	missingParams := []string{}
	for _, param := range requiredParams {
		if _, exists := input[param]; !exists {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) == 0 {
		return nil
	}

	return errors.New("MissingParams: " + strings.Join(missingParams, ", "))
}
