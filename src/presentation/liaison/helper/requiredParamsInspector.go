package liaisonHelper

import (
	"errors"
	"strings"
)

func RequiredParamsInspector(
	untrustedInput map[string]any,
	requiredParams []string,
) error {
	missingParams := []string{}
	for _, param := range requiredParams {
		if _, exists := untrustedInput[param]; !exists {
			missingParams = append(missingParams, param)
		}
	}

	if len(missingParams) == 0 {
		return nil
	}

	return errors.New("MissingParams: " + strings.Join(missingParams, ", "))
}
