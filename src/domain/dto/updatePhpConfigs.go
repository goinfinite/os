package dto

import (
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/valueObject"
	tkValueObject "github.com/goinfinite/tk/src/domain/valueObject"
)

type UpdatePhpConfigsRequest struct {
	Hostname          tkValueObject.Fqdn      `json:"hostname"`
	PhpVersion        valueObject.PhpVersion  `json:"version"`
	PhpModules        []entity.PhpModule      `json:"modules"`
	PhpSettings       []entity.PhpSetting     `json:"settings"`
	OperatorAccountId tkValueObject.AccountId `json:"-"`
	OperatorIpAddress tkValueObject.IpAddress `json:"-"`
}

type PhpModuleParsingFailure struct {
	Index  uint                       `json:"index"`
	Name   *valueObject.PhpModuleName `json:"name,omitempty"`
	Status *bool                      `json:"status,omitempty"`
	Reason valueObject.FailureReason  `json:"reason"`
}

type UpdatePhpConfigsResponse struct {
	TaskId                         *valueObject.ScheduledTaskId `json:"taskId,omitempty"`
	FailedModulesWithParsingErrors []PhpModuleParsingFailure    `json:"failedModulesWithParsingErrors"`
}

func NewUpdatePhpConfigsRequest(
	hostname tkValueObject.Fqdn,
	phpVersion valueObject.PhpVersion,
	phpModules []entity.PhpModule,
	phpSettings []entity.PhpSetting,
	operatorAccountId tkValueObject.AccountId,
	operatorIpAddress tkValueObject.IpAddress,
) UpdatePhpConfigsRequest {
	return UpdatePhpConfigsRequest{
		Hostname:          hostname,
		PhpVersion:        phpVersion,
		PhpModules:        phpModules,
		PhpSettings:       phpSettings,
		OperatorAccountId: operatorAccountId,
		OperatorIpAddress: operatorIpAddress,
	}
}

func NewUpdatePhpConfigsResponse(
	taskId *valueObject.ScheduledTaskId,
	failedModulesWithParsingErrors []PhpModuleParsingFailure,
) UpdatePhpConfigsResponse {
	if failedModulesWithParsingErrors == nil {
		failedModulesWithParsingErrors = []PhpModuleParsingFailure{}
	}

	return UpdatePhpConfigsResponse{
		TaskId:                         taskId,
		FailedModulesWithParsingErrors: failedModulesWithParsingErrors,
	}
}
