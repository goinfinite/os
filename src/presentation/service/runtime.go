package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/entity"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	runtimeInfra "github.com/goinfinite/os/src/infra/runtime"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
	sharedHelper "github.com/goinfinite/os/src/presentation/shared/helper"
)

type RuntimeService struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	availabilityInspector *sharedHelper.ServiceAvailabilityInspector
	runtimeQueryRepo      runtimeInfra.RuntimeQueryRepo
	runtimeCmdRepo        *runtimeInfra.RuntimeCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewRuntimeService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *RuntimeService {
	return &RuntimeService{
		persistentDbSvc: persistentDbSvc,
		availabilityInspector: sharedHelper.NewServiceAvailabilityInspector(
			persistentDbSvc,
		),
		runtimeQueryRepo:      runtimeInfra.RuntimeQueryRepo{},
		runtimeCmdRepo:        runtimeInfra.NewRuntimeCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *RuntimeService) ReadPhpConfigs(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	phpConfigs, err := useCase.ReadPhpConfigs(service.runtimeQueryRepo, hostname)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, phpConfigs)
}

func (service *RuntimeService) UpdatePhpConfigs(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	if !service.availabilityInspector.IsAvailable(serviceName) {
		return NewServiceOutput(InfraError, sharedHelper.ServiceUnavailableError)
	}

	requiredParams := []string{"hostname", "version"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	phpVersion, err := valueObject.NewPhpVersion(input["version"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	phpModules := []entity.PhpModule{}
	if _, exists := input["modules"]; exists {
		var assertOk bool
		phpModules, assertOk = input["modules"].([]entity.PhpModule)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPhpModules")
		}
	}

	phpSettings := []entity.PhpSetting{}
	if _, exists := input["settings"]; exists {
		var assertOk bool
		phpSettings, assertOk = input["settings"].([]entity.PhpSetting)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidPhpSettings")
		}
	}

	operatorAccountId := LocalOperatorAccountId
	if input["operatorAccountId"] != nil {
		operatorAccountId, err = valueObject.NewAccountId(input["operatorAccountId"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	operatorIpAddress := LocalOperatorIpAddress
	if input["operatorIpAddress"] != nil {
		operatorIpAddress, err = valueObject.NewIpAddress(input["operatorIpAddress"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}

	updateDto := dto.NewUpdatePhpConfigs(
		hostname, phpVersion, phpModules, phpSettings, operatorAccountId,
		operatorIpAddress,
	)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	err = useCase.UpdatePhpConfigs(
		service.runtimeQueryRepo, service.runtimeCmdRepo, vhostQueryRepo,
		service.activityRecordCmdRepo, updateDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "PhpConfigsUpdated")
}
