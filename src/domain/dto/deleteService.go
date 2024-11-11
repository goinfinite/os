package dto

import "github.com/goinfinite/os/src/domain/valueObject"

type DeleteService struct {
	Name              valueObject.ServiceName `json:"name"`
	OperatorAccountId valueObject.AccountId   `json:"-"`
	OperatorIpAddress valueObject.IpAddress   `json:"-"`
}

func NewDeleteService(
	name valueObject.ServiceName,
	operatorAccountId valueObject.AccountId,
	operatorIpAddress valueObject.IpAddress,
) DeleteService {
	return DeleteService{
		Name:              name,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}
