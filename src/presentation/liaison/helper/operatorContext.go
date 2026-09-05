package liaisonHelper

import (
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadOperatorContext(
	untrustedInput map[string]any,
) (tkValueObject.AccountId, tkValueObject.IpAddress, error) {
	operatorAccountId := sharedHelper.LocalOperatorAccountId
	if untrustedInput["operatorAccountId"] != nil {
		parsedAccountId, err := tkValueObject.NewAccountId(
			untrustedInput["operatorAccountId"],
		)
		if err != nil {
			return operatorAccountId, sharedHelper.LocalOperatorIpAddress, err
		}
		operatorAccountId = parsedAccountId
	}

	operatorIpAddress := sharedHelper.LocalOperatorIpAddress
	if untrustedInput["operatorIpAddress"] != nil {
		parsedIpAddress, err := tkValueObject.NewIpAddress(
			untrustedInput["operatorIpAddress"],
		)
		if err != nil {
			return operatorAccountId, operatorIpAddress, err
		}
		operatorIpAddress = parsedIpAddress
	}

	return operatorAccountId, operatorIpAddress, nil
}
