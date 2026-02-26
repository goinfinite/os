package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type DeleteService struct {
	Name              valueObject.ServiceName `json:"name"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

func NewDeleteService(
	name valueObject.ServiceName,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) DeleteService {
	return DeleteService{
		Name:              name,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
