package service

import (
	"github.com/goinfinite/os/src/domain/dto"
	"github.com/goinfinite/os/src/domain/useCase"
	"github.com/goinfinite/os/src/domain/valueObject"
	voHelper "github.com/goinfinite/os/src/domain/valueObject/helper"
	activityRecordInfra "github.com/goinfinite/os/src/infra/activityRecord"
	infraHelper "github.com/goinfinite/os/src/infra/helper"
	internalDbInfra "github.com/goinfinite/os/src/infra/internalDatabase"
	servicesInfra "github.com/goinfinite/os/src/infra/services"
	vhostInfra "github.com/goinfinite/os/src/infra/vhost"
	mappingInfra "github.com/goinfinite/os/src/infra/vhost/mapping"
	serviceHelper "github.com/goinfinite/os/src/presentation/service/helper"
)

type VirtualHostService struct {
	persistentDbSvc       *internalDbInfra.PersistentDatabaseService
	vhostQueryRepo        *vhostInfra.VirtualHostQueryRepo
	vhostCmdRepo          *vhostInfra.VirtualHostCmdRepo
	mappingQueryRepo      *mappingInfra.MappingQueryRepo
	mappingCmdRepo        *mappingInfra.MappingCmdRepo
	activityRecordCmdRepo *activityRecordInfra.ActivityRecordCmdRepo
}

func NewVirtualHostService(
	persistentDbSvc *internalDbInfra.PersistentDatabaseService,
	trailDbSvc *internalDbInfra.TrailDatabaseService,
) *VirtualHostService {
	return &VirtualHostService{
		persistentDbSvc:       persistentDbSvc,
		vhostQueryRepo:        vhostInfra.NewVirtualHostQueryRepo(persistentDbSvc),
		vhostCmdRepo:          vhostInfra.NewVirtualHostCmdRepo(persistentDbSvc),
		mappingQueryRepo:      mappingInfra.NewMappingQueryRepo(persistentDbSvc),
		mappingCmdRepo:        mappingInfra.NewMappingCmdRepo(persistentDbSvc),
		activityRecordCmdRepo: activityRecordInfra.NewActivityRecordCmdRepo(trailDbSvc),
	}
}

func (service *VirtualHostService) Read() ServiceOutput {
	vhostsList, err := useCase.ReadVirtualHosts(service.vhostQueryRepo)
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
		rawVhostType, err = voHelper.InterfaceToString(input["type"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
	}
	vhostType, err := valueObject.NewVirtualHostType(rawVhostType)
	if err != nil {
		return NewServiceOutput(UserError, err.Error())
	}

	var parentHostnamePtr *valueObject.Fqdn
	if input["parentHostname"] != nil {
		parentHostname, err := valueObject.NewFqdn(input["parentHostname"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}
		parentHostnamePtr = &parentHostname
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

	createDto := dto.NewCreateVirtualHost(
		hostname, vhostType, parentHostnamePtr, operatorAccountId, operatorIpAddress,
	)

	err = useCase.CreateVirtualHost(
		service.vhostQueryRepo, service.vhostCmdRepo, service.activityRecordCmdRepo,
		createDto,
	)
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

	deleteDto := dto.NewDeleteVirtualHost(
		hostname, primaryVhost, operatorAccountId, operatorIpAddress,
	)

	err = useCase.DeleteVirtualHost(
		service.vhostQueryRepo, service.vhostCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Success, "VirtualHostDeleted")
}

func (service *VirtualHostService) ReadWithMappings() ServiceOutput {
	vhostsWithMappings, err := useCase.ReadVirtualHostsWithMappings(
		service.mappingQueryRepo,
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
		typedRawMatchPattern, err := voHelper.InterfaceToString(input["matchPattern"])
		if err != nil {
			return NewServiceOutput(UserError, err.Error())
		}

		if len(typedRawMatchPattern) > 0 {
			rawMatchPattern = typedRawMatchPattern
		}
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

	createDto := dto.NewCreateMapping(
		hostname, path, matchPattern, targetType, targetValuePtr,
		targetHttpResponseCodePtr, operatorAccountId, operatorIpAddress,
	)

	servicesQueryRepo := servicesInfra.NewServicesQueryRepo(service.persistentDbSvc)

	err = useCase.CreateMapping(
		service.mappingQueryRepo, service.mappingCmdRepo, service.vhostQueryRepo,
		servicesQueryRepo, service.activityRecordCmdRepo, createDto,
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

	deleteDto := dto.NewDeleteMapping(id, operatorAccountId, operatorIpAddress)

	err = useCase.DeleteMapping(
		service.mappingQueryRepo, service.mappingCmdRepo, service.activityRecordCmdRepo,
		deleteDto,
	)
	if err != nil {
		return NewServiceOutput(InfraError, err.Error())
	}

	return NewServiceOutput(Created, "MappingDeleted")
}
