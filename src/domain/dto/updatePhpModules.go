package dto

import (
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdatePhpModulesRequest struct {
	Hostname          tkValueObject.Fqdn      `json:"hostname"`
	PhpVersion        valueObject.PhpVersion  `json:"version"`
	PhpModules        []PhpModuleUpdate       `json:"modules"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

type PhpModuleUpdate struct {
	Name   valueObject.PhpModuleName `json:"name"`
	Status bool                      `json:"status"`
}

type PhpModuleUpdateFailure struct {
	Name   valueObject.PhpModuleName `json:"name"`
	Reason valueObject.FailureReason `json:"reason"`
	Status bool                      `json:"status"`
}

type UpdatePhpModulesResponse struct {
	ModulesSuccessfullyUpdated []PhpModuleUpdate        `json:"modulesSuccessfullyUpdated"`
	FailedModulesWithReason    []PhpModuleUpdateFailure `json:"failedModulesWithReason"`
}

func NewUpdatePhpModulesRequest(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	phpModules []PhpModuleUpdate,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdatePhpModulesRequest {
	return UpdatePhpModulesRequest{
		Hostname:          hostname,
		PhpVersion:        phpVersion,
		PhpModules:        phpModules,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}

func NewPhpModuleUpdate(
	name valueObject.PhpModuleName,
	status bool,
) PhpModuleUpdate {
	return PhpModuleUpdate{Name: name, Status: status}
}

func NewPhpModuleUpdateFailure(
	name valueObject.PhpModuleName,
	status bool,
	reason valueObject.FailureReason,
) PhpModuleUpdateFailure {
	return PhpModuleUpdateFailure{
		Name:   name,
		Reason: reason,
		Status: status,
	}
}

func NewUpdatePhpModulesResponse(
	modulesSuccessfullyUpdated []PhpModuleUpdate,
	failedModulesWithReason []PhpModuleUpdateFailure,
) UpdatePhpModulesResponse {
	return UpdatePhpModulesResponse{
		ModulesSuccessfullyUpdated: modulesSuccessfullyUpdated,
		FailedModulesWithReason:    failedModulesWithReason,
	}
}
