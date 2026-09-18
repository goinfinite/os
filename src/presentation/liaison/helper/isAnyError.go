package liaisonHelper

import (
	"errors"
	"slices"
)

func IsAnyError(err error, targets ...error) bool {
	return slices.ContainsFunc(targets, func(target error) bool {
		return errors.Is(err, target)
	})
}
