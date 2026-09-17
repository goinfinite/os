package useCaseHelper

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/repository"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

func ReadOperatorAccount(
	accountQueryRepo repository.AccountQueryRepo,
	operatorAccountId tkValueObject.AccountId,
) (operatorAccountEntity entity.Account, isSystemOperator bool, err error) {
	if operatorAccountId == tkValueObject.AccountIdSystem {
		return operatorAccountEntity, true, nil
	}

	operatorAccountEntity, err = accountQueryRepo.ReadFirst(dto.ReadAccountsRequest{
		AccountId: &operatorAccountId,
	})
	if err != nil {
		return operatorAccountEntity, false, err
	}

	return operatorAccountEntity, false, nil
}
