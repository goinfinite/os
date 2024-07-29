package service

import (
	"github.com/speedianet/os/src/domain/dto"
	"github.com/speedianet/os/src/domain/useCase"
	"github.com/speedianet/os/src/domain/valueObject"
	internalDbInfra "github.com/speedianet/os/src/infra/internalDatabase"
	vhostInfra "github.com/speedianet/os/src/infra/vhost"
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

	hostname, err := valueObject.NewFqdn(input["hostname"].(string))
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

	return NewServiceOutput(Success, "VirtualHostCreated")
}
