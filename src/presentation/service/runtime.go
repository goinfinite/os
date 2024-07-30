package service

import (
	"log"

	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/entity"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	runtimeInfra "github.com/speedianet/os/src/infra/runtime"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
	sharedHelper "github.com/speedianet/os/src/presentation/shared/helper"
)

type RuntimeService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewRuntimeService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *RuntimeService {
	return &RuntimeService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *RuntimeService) ReadPhpConfigs(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	sharedHelper.StopIfServiceUnavailable(service.persistentDbSvc, serviceName)

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		log.Printf("%s --> '%+v'", err.Error(), input)
		return NewServiceOutput(UserError, err.Error())
	}

	runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
	phpConfigs, err := useCase.ReadPhpConfigs(runtimeQueryRepo, hostname)
	if err != nil {
		log.Printf("%s --> '%+v'", err.Error(), input)
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, phpConfigs)
}

func (service *RuntimeService) UpdatePhpConfigs(
	input map[string]interface{},
) ServiceOutput {
	serviceName, _ := valueObject.NewServiceName("php-webserver")
	sharedHelper.StopIfServiceUnavailable(service.persistentDbSvc, serviceName)

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

	dto := dto.NewUpdatePhpConfigs(hostname, phpVersion, phpModules, phpSettings)

	runtimeQueryRepo := runtimeInfra.RuntimeQueryRepo{}
	runtimeCmdRepo := runtimeInfra.NewRuntimeCmdRepo(service.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)

	err = useCase.UpdatePhpConfigs(
		runtimeQueryRepo, runtimeCmdRepo, vhostQueryRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "PhpConfigsUpdated")
}
