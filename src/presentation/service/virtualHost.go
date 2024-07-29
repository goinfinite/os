package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	infraHelper "github.com/speedianet/os/src/infra/helper"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	servicesInfra "github.com/speedianet/os/src/infra/services"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
	mappingInfra "github.com/speedianet/os/src/infra/vhost/mapping"
	serviceHelper "github.com/speedianet/os/src/presentation/service/helper"
)

type VirtualHostService struct {
	persistentDbSvc *internalDbInfra.PersistentDatabaseService
}

func NewVirtualHostService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
) *VirtualHostService {
	return &VirtualHostService{
		persistentDbSvc: persistentDbSvc,
	}
}

func (service *VirtualHostService) Read() ServiceOutput {
	vhostsQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostsList, err := useCase.ReadVirtualHosts(vhostsQueryRepo)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, vhostsList)
}

func (service *VirtualHostService) Create(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"hostname"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	rawVhostType := "top-level"
	if input["type"] != nil {
		var assertOk bool
		rawVhostType, assertOk = input["type"].(string)
		if !assertOk {
			return NewServiceOutput(UserError, "InvalidVirtualHostType")
		}
	}
	vhostType, err := valueObject.NewVirtualHostType(rawVhostType)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var parentHostnamePtr *valueObject.Fqdn
	if input["parentHostname"] != nil {
		parentHostname, err := valueObject.NewFqdn(
			input["parentHostname"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		parentHostnamePtr = &parentHostname
	}

	dto := dto.NewCreateVirtualHost(hostname, vhostType, parentHostnamePtr)

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(service.persistentDbSvc)

	err = useCase.CreateVirtualHost(vhostQueryRepo, vhostCmdRepo, dto)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "VirtualHostCreated")
}

func (service *VirtualHostService) Delete(input map[string]interface{}) ServiceOutput {
	requiredParams := []string{"hostname"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	primaryVhost, err := infraHelper.GetPrimaryVirtualHost()
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	vhostCmdRepo := vhostInfra.NewVirtualHostCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteVirtualHost(
		vhostQueryRepo, vhostCmdRepo, primaryVhost, hostname,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "VirtualHostDeleted")
}

func (service *VirtualHostService) ReadWithMappings() ServiceOutput {
	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(service.persistentDbSvc)
	vhostsWithMappings, err := useCase.ReadVirtualHostsWithMappings(
		mappingQueryRepo,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, vhostsWithMappings)
}

func (service *VirtualHostService) CreateMapping(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"hostname", "path", "targetType"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	hostname, err := valueObject.NewFqdn(input["hostname"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	path, err := valueObject.NewMappingPath(input["path"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	rawMatchPattern := "begins-with"
	if input["matchPattern"] != nil {
		rawMatchPattern = input["matchPattern"].(string)
	}
	matchPattern, err := valueObject.NewMappingMatchPattern(rawMatchPattern)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	targetType, err := valueObject.NewMappingTargetType(input["targetType"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var targetValuePtr *valueObject.MappingTargetValue
	if input["targetValue"] != nil {
		targetValue, err := valueObject.NewMappingTargetValue(
			input["targetValue"], targetType,
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		targetValuePtr = &targetValue
	}

	var targetHttpResponseCodePtr *valueObject.HttpResponseCode
	if input["targetHttpResponseCode"] != nil {
		targetHttpResponseCode, err := valueObject.NewHttpResponseCode(
			input["targetHttpResponseCode"],
		)
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		targetHttpResponseCodePtr = &targetHttpResponseCode
	}

	dto := dto.NewCreateMapping(
		hostname, path, matchPattern, targetType, targetValuePtr,
		targetHttpResponseCodePtr,
	)

	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(service.persistentDbSvc)
	vhostQueryRepo := vhostInfra.NewVirtualHostQueryRepo(service.persistentDbSvc)
	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMapping(
		mappingQueryRepo, mappingCmdRepo, vhostQueryRepo, servicesQueryRepo, dto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingCreated")
}

func (service *VirtualHostService) DeleteMapping(
	input map[string]interface{},
) ServiceOutput {
	requiredParams := []string{"id"}
	err := serviceHelper.RequiredParamsInspector(input, requiredParams)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	id, err := valueObject.NewMappingId(input["id"])
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	mappingQueryRepo := mappingInfra.NewMappingQueryRepo(service.persistentDbSvc)
	mappingCmdRepo := mappingInfra.NewMappingCmdRepo(service.persistentDbSvc)

	err = useCase.DeleteMapping(mappingQueryRepo, mappingCmdRepo, id)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingDeleted")
}
